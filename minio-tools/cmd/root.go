package cmd

import (
	"log"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "miniotools",
	Short: "上传文件到minio以提供给用户下载",
	Long:  "上传文件到minio以提供给用户下载",
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatalf("rootCmd.Execute:%s", err)
		os.Exit(1)
	}

}

var rootTest string

func init() {
	log.Println("root init")
	rootCmd.PersistentFlags().StringVar(&rootTest, "root", "rootDefault", "aaaa")
}
