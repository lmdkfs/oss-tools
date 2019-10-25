package cmd

import (
	"fmt"
	"log"
	"os"
	miniocfg "oss-tools/minio-tools/config"
	"oss-tools/minio-tools/pkg"
	"path/filepath"

	"github.com/spf13/cobra"
)

const (
	DefaultEndPoint        string = "10.10.186.58:9001"
	DefaultAccessKeyID     string = "minio"
	DefaultSecretAccessKey string = "minio123"
	DefaultSecure          bool   = false
	DefaultBucketName      string = "public"
	DefaultLocation        string = "us-east-1"
)

var uploadCmd = &cobra.Command{
	Use:   "up",
	Short: "up file to minio storage",
	Long:  "up file to minio storage",
	Run: func(cmd *cobra.Command, args []string) {
		pkg.UploadToMinio()
		fmt.Println("up file to minio")
	},
}

func init() {
	rootCmd.AddCommand(uploadCmd)
	cobra.OnInitialize(initUploadConfig)

	uploadCmd.PersistentFlags().String("ep", DefaultEndPoint, "minio endpoint url")
	uploadCmd.PersistentFlags().String("ak", "", "minio accessKeyID")
	uploadCmd.PersistentFlags().String("sk", "", "minio secretAccessKey")
	uploadCmd.PersistentFlags().StringP("file", "f", "", "filename for upload file")
	uploadCmd.MarkFlagRequired("file")
	uploadCmd.PersistentFlags().String("obj", "", "objectName")

	// uploadCmd.PersistentFlags().MarkHidden("ak")
}
func initUploadConfig() {
	var config = miniocfg.NewConfig()
	config.EndPoint = DefaultEndPoint
	config.AccessKeyID = DefaultAccessKeyID
	config.SecretAccessKey = DefaultSecretAccessKey
	config.Secure = DefaultSecure
	if ep := uploadCmd.PersistentFlags().Lookup("ep").Value.String(); ep != "" {
		config.EndPoint = ep
	}
	if ak := uploadCmd.PersistentFlags().Lookup("ak").Value.String(); ak != "" {
		config.AccessKeyID = ak
	}
	if sk := uploadCmd.PersistentFlags().Lookup("sk").Value.String(); sk != "" {
		config.SecretAccessKey = sk
	}
	fileName := uploadCmd.PersistentFlags().Lookup("file").Value.String()

	if _, err := os.Stat(fileName); os.IsNotExist(err) {
		log.Printf("Errof: %s doesn't exist, please check fileName ", fileName)
		os.Exit(1)
	} else {
		config.FileName = fileName
	}
	if obj := uploadCmd.PersistentFlags().Lookup("sk").Value.String(); obj != "" {
		config.ObjectName = obj
	} else {
		config.ObjectName = filepath.Base(fileName)
	}

}
