package pkg

import (
	"io"
	"log"
	"net/http"
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

func multiPutToUFile(ch <-chan FileInfo, wg *sync.WaitGroup) {

	for fileInfo := range ch {
		resp, err := http.Get(fileInfo.qiniuDownloadUrl)
		if err != nil {
			logger.Println("download fail from qiniu:", err.Error())
		} else {
			uploadToUFile(resp.Body, fileInfo.fileName)
		}
		if resp != nil {
			err = resp.Body.Close()
			if err != nil {
				logger.Println("http Body close fail:", err)
			}

		}
	}

	wg.Done()
}

func uploadToUFile(uploadFileReader io.Reader, remoteFileKey string) {
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
			logger.Printf("%s 文件上传成功", remoteFileKey)
		} else {
			logger.Errorf("%s 文件上传失败， err:%s", remoteFileKey, err)
		}
	} else {
		logger.Error("read ufile config fail: ", err)
	}

}

// generate downloadUrl from qiniu
func GenPrivateUrl(key string) (privateUrl string) {
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

func getFileFromQiniu(ch chan FileInfo, wg *sync.WaitGroup) {
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
	prefix, delimiter, marker, limit := "", "", "", 1000
	for {
		entries, _, nextMarker, hashNext, err := bucketManager.ListFiles(cfg.Qiniu.Bucket, prefix, delimiter, marker, limit)
		if err != nil {
			logger.Println("list error, ", err)
			break
		}
		for _, entry := range entries {

			logger.Println(entry.Key)
			downloadFileUrl := GenPrivateUrl(entry.Key)
			fileInfo := FileInfo{
				fileName:         entry.Key,
				qiniuDownloadUrl: downloadFileUrl,
			}
			ch <- fileInfo
		}

		if hashNext {
			marker = nextMarker
		} else {
			break
		}

	}
	wg.Done()
}

func Migrate() {
	cfg = config.NewGlobalConfig()
	logger = utils.NewLogger()
	numcpu := runtime.NumCPU()
	runtime.GOMAXPROCS(numcpu)
	Ch = make(chan FileInfo, 100)
	go getFileFromQiniu(Ch, &wgp)
	wgp.Add(1)
	for i := 0; i < cfg.Worker; i++ {
		go multiPutToUFile(Ch, &wgc)
		wgc.Add(1)
	}
	wgp.Wait()
	close(Ch)
	wgc.Wait()
	logger.Println("migrate 完成")
}
