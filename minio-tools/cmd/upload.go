package cmd

import (
	"fmt"
	miniocfg "oss-tools/minio-tools/config"
	"oss-tools/minio-tools/pkg"

	"github.com/spf13/cobra"
)

const (
	DefaultEndPoint        string = "endponit"
	DefaultAccessKeyID     string = "accessKeyID"
	DefaultSecretAccessKey string = "secretAccessKey"
	DefaultSecure          bool   = false
)

var uploadCmd = &cobra.Command{
	Use:   "up",
	Short: "up file to minio storage",
	Long:  "up file to minio storage",
	Run: func(cmd *cobra.Command, args []string) {
		pkg.Test()
		fmt.Println("up file to minio")
	},
}

func init() {
	rootCmd.AddCommand(uploadCmd)
	cobra.OnInitialize(initUploadConfig)

	uploadCmd.PersistentFlags().String("ep", "", "minio endpoint url")
	uploadCmd.PersistentFlags().String("ak", "", "minio accessKeyID")
	uploadCmd.PersistentFlags().String("sk", "", "minio secretAccessKey")
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
}
