package pkg

import (
	"log"
	"net/url"
	"os"
	"oss-tools/minio-tools/config"
	"time"

	"github.com/minio/minio-go/v6"
)

func UploadToMinio() {
	config := config.NewConfig()
	if _, err := os.Stat(config.FileName); os.IsNotExist(err) {
		log.Printf("Errof: %s doesn't exist, please check fileName ", config.FileName)
		os.Exit(1)
	}
	minioClient, err := minio.New(config.EndPoint, config.AccessKeyID, config.SecretAccessKey, config.Secure)
	if err != nil {
		log.Fatalln("New minioClient error:", err)

	}

	err = minioClient.MakeBucket(config.BucketName, config.Location)
	if err != nil {
		exists, errBucketExists := minioClient.BucketExists(config.BucketName)
		if errBucketExists == nil && exists {
			log.Printf("we already own bucketName: %s \n", config.BucketName)
		} else {
			log.Fatalln(err)
		}
	} else {
		log.Printf("Sucessfully created %s \n", config.BucketName)
	}
	// Upload the file
	contentType := "application/binary"

	// Upload the rke file with FPutObject
	n, err := minioClient.FPutObject(config.BucketName, config.ObjectName, config.FileName, minio.PutObjectOptions{ContentType: contentType})
	if err != nil {
		log.Fatalln(err)
	}
	log.Printf("Successfully uploaded %s of size %d\n", config.ObjectName, n)

	reqParams := make(url.Values)

	reqParams.Set("response-content-disposition", "attachment; filename="+config.ObjectName)
	presignedURL, err := minioClient.PresignedGetObject(config.BucketName, config.ObjectName, time.Second*24*60*60, reqParams)
	if err != nil {
		log.Println("签名失败", err)
	}
	log.Printf("Successfully generated presigned URL: 复制连接到浏览器即可下载\n %s", presignedURL)
}
