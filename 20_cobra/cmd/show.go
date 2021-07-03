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

/*
需求：
（1）show 查看当前时间
（2）parse 指定时间格式 --format，parse为show的子命令。
 */
// showCmd represents the show command
var showCmd = &cobra.Command{
	Use:   "show",
	Short: "show简短介绍展示当前时间",
	Long: `show详细介绍`,
	Run: ShowTime,
}
// ShowTime 显示当前时间
func ShowTime(cmd *cobra.Command, args []string) {
	fmt.Println(time.Now())
}


func init() {
	rootCmd.AddCommand(showCmd)

}
