package pkg

import (
	"bufio"
	"github.com/gogo/protobuf/proto"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"runtime"
	"strings"
	"sync"
	"time"

	"oss-tools/sync-tools/config"
	"oss-tools/utils"

	"github.com/qiniu/api.v7/auth/qbox"
	"github.com/qiniu/api.v7/storage"
	"github.com/sirupsen/logrus"
	ufsdk "github.com/ufilesdk-dev/ufile-gosdk"
)

type FileInfo struct {
	fileName         string
	qiniuDownloadUrl string
}

var (
	wgp    sync.WaitGroup
	wgc    sync.WaitGroup
	cfg    *config.Config
	Ch     chan FileInfo
	logger *logrus.Logger
)
func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}


func multiPutToUFile(workerID int, ch <-chan FileInfo, wg *sync.WaitGroup, rwMutex *sync.RWMutex, ucache *UfileCache) {
	defer wg.Done()
	for fileInfo := range ch {
		resp, err := http.Get(fileInfo.qiniuDownloadUrl)
		if err != nil {
			logger.Printf("workerID: %d, download fail from qiniu: %s", workerID, err.Error())
		} else {
			rwMutex.Lock()
			if _, ok := ucache.Ufilecache[fileInfo.fileName]; ok {
				logger.Printf("workerID: %d, fileName:%s is already in cache", workerID, fileInfo.fileName)
			} else {
				uploadToUFile(workerID, resp.Body, fileInfo.fileName)

				ucache.Ufilecache[fileInfo.fileName] = true

			}
			rwMutex.Unlock()

		}
		if resp != nil {
			err = resp.Body.Close()
			if err != nil {
				logger.Println("http Body close fail:", err)
			}

		}
	}

}

func uploadToUFile(workerID int, uploadFileReader io.Reader, remoteFileKey string) {
	config := &ufsdk.Config{}
	config.PublicKey = cfg.Ufile.PublicKey
	config.PrivateKey = cfg.Ufile.PrivateKey
	config.FileHost = cfg.Ufile.FileHost
	config.BucketName = cfg.Ufile.Bucket

	if req, err := ufsdk.NewFileRequest(config, nil); err == nil {
		log.Println("正在上传数据...")
		if cfg.Ufile.Prefix != "" {
			remoteFileKey = cfg.Ufile.Prefix + remoteFileKey
		}
		if err := req.IOPut(uploadFileReader, remoteFileKey, ""); err == nil {
			//if err := req.IOPut(uploadFileReader, remoteFileKey, ""); err == nil {
			logger.Printf("workerID:%d, %s 文件上传成功", workerID, remoteFileKey)
		} else {
			logger.Errorf("workerID:%d, %s 文件上传失败， err:%s", workerID, remoteFileKey, err)
			logger.Println(string(req.DumpResponse(true)))
		}
	} else {
		logger.Error("read ufile config fail: ", err)
	}

}

// generate downloadUrl from qiniu
func GenDownloadUrl(key string) (privateUrl string) {
	if strings.HasPrefix(key, "/") {
		key = "@" + key
	}
	if cfg.Qiniu.PrivateBucket {
		mac := qbox.NewMac(cfg.Qiniu.Ak, cfg.Qiniu.Sk)
		deadLine := time.Now().Add(time.Second * time.Duration(1000000)).Unix()
		privateUrl = storage.MakePrivateURL(mac, cfg.Qiniu.QiniuDomain, key, deadLine)
	} else {
		privateUrl = cfg.Qiniu.QiniuDomain + key
	}
	return
}

// getFileFromQiniuSourceFile
func getFileFromQiniuSourceFile(fileName string, wg *sync.WaitGroup) <- chan FileInfo{
	out := make(chan FileInfo, 100)
	currentDir, err := os.Getwd()
	if err != nil {
		logger.Panic("获取当前目录失败:", err)
	}
	var filePath string
	for _, f := range []string{currentDir + "/" + fileName, fileName} {
		if fileExists(f) {
			filePath = f
			break
		}
	}
	if filePath == "" {
		logger.Panicf("文件%s or %s 不存在", currentDir + "/" + fileName, fileName)
	}
	go func() {
		defer wg.Done()
		defer close(out)
		f, err := os.Open(filePath)
		if err != nil {
			logger.Panic("Error:", err)
		}
		defer f.Close()
		br := bufio.NewReader(f)
		for {
			a, _, c := br.ReadLine()
			if c == io.EOF{
				break
			}
			downloadFileUrl := GenDownloadUrl(string(a))
			fileInfo := FileInfo{
				fileName: string(a),
				qiniuDownloadUrl: downloadFileUrl,
			}
			out <- fileInfo
		}
	}()

	return out
}
func getFileFromQiniuMetadata( wg *sync.WaitGroup) <- chan FileInfo {
	out := make(chan FileInfo, 100)
	mac := qbox.NewMac(cfg.Qiniu.Ak, cfg.Qiniu.Sk)
	zone := storage.Zone{
		SrcUpHosts: []string{cfg.Qiniu.UpHost},
		RsfHost:    cfg.Qiniu.RsfHost,
		RsHost:     cfg.Qiniu.RsHost,
		IovipHost:  cfg.Qiniu.IoVipHost,
		ApiHost:    cfg.Qiniu.ApiHost,
	}
	qiniuConfig := storage.Config{
		UseHTTPS: false,
	}
	qiniuConfig.Zone = &zone
	bucketManager := storage.NewBucketManager(mac, &qiniuConfig)

	// 列举所有文件
	go func() {
		defer wg.Done()
		defer close(out)
		prefix, delimiter, marker, limit := "", "", "", 1000
		for {
			entries, _, nextMarker, hashNext, err := bucketManager.ListFiles(cfg.Qiniu.Bucket, prefix, delimiter, marker, limit)
			if err != nil {
				logger.Println("list error, ", err)
				break
			}
			for _, entry := range entries {

				logger.Println(entry.Key)
				downloadFileUrl := GenDownloadUrl(entry.Key)
				fileInfo := FileInfo{
					fileName:         entry.Key,
					qiniuDownloadUrl: downloadFileUrl,
				}
				out <- fileInfo
			}

			if hashNext {
				marker = nextMarker
			} else {
				break
			}

		}

	}()

	return out
}

func Migrate() {
	cfg = config.NewGlobalConfig()
	logger = utils.NewLogger()
	numcpu := runtime.NumCPU()
	runtime.GOMAXPROCS(numcpu)
	//Ch = make(chan FileInfo, 100)
	var ucache UfileCache
	var rwMutex sync.RWMutex
	var Ch <- chan FileInfo
	ucache = UfileCache{Ufilecache: make(map[string]bool)}
	currentDir, err := os.Getwd()
	if err != nil {
		logger.Panic("获取当前执行路径失败,", err)
	}


	if protoBytes, err := ioutil.ReadFile(currentDir +"/" + cfg.Ufile.Bucket); err == nil {
		if err := proto.Unmarshal(protoBytes, &ucache); err != nil {
			log.Panic("parse protobuf  file failed:", err)
		}
	}
	if cfg.Qiniu.File != "" {
		Ch = getFileFromQiniuSourceFile(cfg.Qiniu.File, &wgp)
	} else {
		Ch = getFileFromQiniuMetadata(&wgp)
	}

	wgp.Add(1)
	for i := 0; i < cfg.Worker; i++ {
		go multiPutToUFile(i, Ch, &wgc, &rwMutex, &ucache)
		wgc.Add(1)
	}
	wgp.Wait()
	wgc.Wait()
	out, err := proto.Marshal(&ucache)
	if err != nil {
		log.Fatalln("Failed to encode ucache:", err)
	}
	if err := ioutil.WriteFile(currentDir +"/" + cfg.Ufile.Bucket, out, 0644); err != nil {
		log.Fatalln("Failed to write ucache:", err)
	}
	logger.Println("migrate 完成")
}
