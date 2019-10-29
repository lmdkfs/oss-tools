package cmd

import (
	"fmt"
	miniocfg "oss-tools/minio-tools/config"
	"oss-tools/minio-tools/pkg"
	"path/filepath"
	"time"

	"github.com/spf13/cobra"
)

const (
	DefaultEndPoint        string = "minio.puhuitech.cn:9001"
	DefaultAccessKeyID     string = "minio"
	DefaultSecretAccessKey string = "DvVeTXU9"
	DefaultSecure          bool   = false
	DefaultBucketName      string = "paas"
	DefaultLocation        string = "us-east-1"
)

var uploadCmd = &cobra.Command{
	Use:   "up",
	Short: "up file to minio storage",
	Long:  "up file to minio storage",
	Run: func(cmd *cobra.Command, args []string) {
		pkg.UploadToMinio()
		// fmt.Println("up file to minio")
	},
}

func init() {
	rootCmd.AddCommand(uploadCmd)
	cobra.OnInitialize(initUploadConfig)

	uploadCmd.Flags().String("ep", DefaultEndPoint, "minio endpoint url")
	uploadCmd.Flags().String("ak", "", "minio accessKeyID")
	uploadCmd.Flags().String("sk", "", "minio secretAccessKey")
	uploadCmd.Flags().StringP("file", "f", "", "filename for upload file")
	uploadCmd.MarkFlagRequired("file")

	uploadCmd.Flags().StringP("ob", "o", "", "objectName of minio server filename ")
	// uploadCmd.PersistentFlags().MarkHidden("ak")
}
func initUploadConfig() {
	var config = miniocfg.NewConfig()
	config.EndPoint = DefaultEndPoint
	config.AccessKeyID = DefaultAccessKeyID
	config.SecretAccessKey = DefaultSecretAccessKey
	config.Secure = DefaultSecure
	config.BucketName = DefaultBucketName
	config.Location = DefaultLocation
	if ep := uploadCmd.Flags().Lookup("ep").Value.String(); ep != "" {
		config.EndPoint = ep
	}
	if ak := uploadCmd.Flags().Lookup("ak").Value.String(); ak != "" {
		config.AccessKeyID = ak
	}
	if sk := uploadCmd.Flags().Lookup("sk").Value.String(); sk != "" {
		config.SecretAccessKey = sk
	}
	config.FileName = uploadCmd.Flags().Lookup("file").Value.String()
	if obj := uploadCmd.Flags().Lookup("sk").Value.String(); obj != "" {
		config.ObjectName = obj
	} else {
		obj = time.Now().Format("2006-01-02-15-04-05")
		config.ObjectName = fmt.Sprint(obj) + "-" + filepath.Base(config.FileName)
	}

}
