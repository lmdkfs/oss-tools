package pkg

import (
	"log"
	"net/url"
	"oss-tools/minio-tools/config"
	"time"

	"github.com/minio/minio-go/v6"
)

func UploadToMinio() {
	config := config.NewConfig()
	minioClient, err := minio.New(config.EndPoint, config.AccessKeyID, config.SecretAccessKey, config.Secure)
	if err != nil {
		log.Fatalln(err)

	}

	bucketName := "zrq"
	location := "us-east-1"
	err = minioClient.MakeBucket(bucketName, location)
	if err != nil {
		exists, errBucketExists := minioClient.BucketExists(bucketName)
		if errBucketExists == nil && exists {
			log.Printf("we already own %s\n", bucketName)
		} else {
			log.Fatalln(err)
		}
	} else {
		log.Printf("Sucessfully created %s \n", bucketName)
	}
	// Upload the file
	contentType := "application/binary"

	// Upload the rke file with FPutObject
	n, err := minioClient.FPutObject(config.BucketName, objectName, config.FileName, minio.PutObjectOptions{ContentType: contentType})
	if err != nil {
		log.Fatalln(err)
	}
	log.Printf("Successfully uploaded %s of size %d\n", objectName, n)

	reqParams := make(url.Values)
	reqParams.Set("response-content-disposition", "attachment; filename="+objectName)
	presignedURL, err := minioClient.PresignedGetObject(bucketName, objectName, time.Second*24*60*60, reqParams)
	if err != nil {
		log.Println("签名失败", err)
	}
	log.Printf("Successfully generated presigned URL:\n %s", presignedURL)
}
