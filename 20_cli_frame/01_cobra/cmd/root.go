/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

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
	"github.com/spf13/cobra"
	"os"

	"github.com/spf13/viper"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands

//Run: 命令执行入口，函数主要写在这个模块中
var rootCmd = &cobra.Command{ // 定义根命令20_cobra
	Use:   "time",
	Short: "Show Current Time",
	Long: `Show Current Time. For example:

With this command, you can view the current time or customize the time format.`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	rootCmd.Execute()
}

// 函数init(): 定义flag和配置处理。
func init() {
	// 配置处理
	cobra.OnInitialize(initConfig)

	// 定义默认flag --config
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.20_cobra.yaml)")

	// 定义默认flag --toggle
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// initConfig(): 用来初始化viper配置文件位置，监听变化。

func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, _ := os.UserHomeDir()

		// Search config in home directory with name ".20_cobra" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".20_cobra")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}
