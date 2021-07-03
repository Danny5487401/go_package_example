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
	"time"

	"github.com/spf13/cobra"
)
var format string

// parseCmd represents the parse command
var parseCmd = &cobra.Command{
	Use:   "parse",
	Short: "Format current time",
	Long: `Format current time. For example:

time show parse --format "2006:01:02 15:04:05".`,
	Run: func(cmd *cobra.Command, args []string) {
		// 输出格式化之后的时间
		fmt.Println(time.Now().Format(format))
	},
}

func init() {
	showCmd.AddCommand(parseCmd)
	parseCmd.Flags().StringVarP(&format, "format", "f", "", "Help message for toggle")
	// 这里指定format flag为必须
	_ = parseCmd.MarkFlagRequired("format")

}
