/*
Copyright © 2019 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"fmt"
	"io"
	"log"
	"os"
	"oss-tools/sync-tools/config"
	"oss-tools/utils"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "sync-tools",
	Short: "A brief description of your application",
	Long:  `从七牛迁移到ucloud的ufile`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	//	Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.cfg.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	rootCmd.PersistentFlags().StringP("log", "l", "oss-tools.log", "logname")
	rootCmd.PersistentFlags().StringP("worker", "c", "15", "worker count")
	rootCmd.PersistentFlags().StringP("file", "f", "", "qiniu metadata from file")
	viper.BindPFlag("log.logname", rootCmd.PersistentFlags().Lookup("log"))
	viper.BindPFlag("worker", rootCmd.PersistentFlags().Lookup("worker"))
	viper.BindPFlag("qiniu.file", rootCmd.PersistentFlags().Lookup("file"))
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	globalConfig := config.NewGlobalConfig()
	currentDirName, err := os.Getwd()
	if err != nil {
		log.Println("get currentdir fail")
	}
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".sync-tools" (without extension).
		viper.AddConfigPath(currentDirName)
		viper.AddConfigPath(home)
		viper.SetConfigName(".cfg")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		log.Println("Using config file:", viper.ConfigFileUsed())
	}
	viper.SetDefault("log.logpath", currentDirName)
	globalConfig.Qiniu.Ak = viper.GetString("qiniu.ak")
	globalConfig.Qiniu.Sk = viper.GetString("qiniu.sk")
	globalConfig.Qiniu.UpHost = viper.GetString("qiniu.uphost")
	globalConfig.Qiniu.RsHost = viper.GetString("qiniu.rshost")
	globalConfig.Qiniu.RsfHost = viper.GetString("qiniu.rsfhost")
	globalConfig.Qiniu.ApiHost = viper.GetString("qiniu.apihost")
	globalConfig.Qiniu.IoVipHost = viper.GetString("qiniu.iovipHost")
	globalConfig.Qiniu.QiniuDomain = viper.GetString("qiniu.qiniudomain")
	globalConfig.Qiniu.Bucket = viper.GetString("qiniu.bucket")
	globalConfig.Qiniu.PrivateBucket = viper.GetBool("qiniu.privateBucket")
	globalConfig.Qiniu.File = viper.GetString("qiniu.file")

	globalConfig.Ufile.PrivateKey = viper.GetString("ufile.private_key")
	globalConfig.Ufile.PublicKey = viper.GetString("ufile.public_key")
	globalConfig.Ufile.Bucket = viper.GetString("ufile.bucket")
	globalConfig.Ufile.FileHost = viper.GetString("ufile.file_host")
	globalConfig.Ufile.Prefix = viper.GetString("ufile.prefix")
	globalConfig.Log.LogPath = viper.GetString("log.logpath")
	globalConfig.Log.LogName = viper.GetString("log.logname")

	globalConfig.Worker = viper.GetInt("worker")
	logger := utils.NewLogger()
	if logfile, err := os.Create(globalConfig.Log.LogPath + "/" + globalConfig.Log.LogName); err != nil {
		logger.Panicf("Create LogFile Fail%s", err)

	} else {
		mw := io.MultiWriter(os.Stdout, logfile)
		logger.SetOutput(mw)
	}
	fmt.Println("logname", globalConfig.Log.LogName)
	fmt.Println("initconfig**")
}
