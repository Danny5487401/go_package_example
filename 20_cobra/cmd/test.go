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
)
// 全局
var (
	Foo *string  //指针类型
	Print string
	show bool
)
// 局部
var (
	FooL *string
	showL bool
	PrintL string

 )
// testCmd represents the test command
var testCmd = &cobra.Command{
	Use:   "test",
	// 定义arguments数量最少为1个,不传Error: requires at least 1 arg(s), only received 0
	Args:  cobra.MinimumNArgs(1),
	Short: "test简短介绍",
	Long: `test详细介绍`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("test called")
		if show {
			fmt.Println("Show")
		}
		fmt.Println("Print:", Print)
		fmt.Println("Foo:", *Foo)
	},
}

func init() {
	rootCmd.AddCommand(testCmd)

	// 1。全局
	// 下面定义了一个Flag foo, foo后面接的值会被赋值给Foo
	Foo = testCmd.PersistentFlags().String("foo", "", "A help for foo")
	// 下面定义了一个Flag print ,print后面的值会被赋值给Print变量
	testCmd.PersistentFlags().StringVar(&Print, "print", "", "print")
	// 下面定义了一个Flag show,show默认为false, 有两种调用方式--show\-s，命令后面接了show则上面定义的show变量就会变成true
	testCmd.PersistentFlags().BoolVarP(&show, "show", "s", false, "show")

	// 2.局部
	// 下面定义了一个Flag foo, foo后面接的值会被赋值给Foo
	FooL = testCmd.Flags().String("fooL", "", "A help for foo")
	// 下面定义了一个Flag print ,print后面的值会被赋值给Print变量
	testCmd.Flags().StringVar(&PrintL, "printL", "", "print")
	// 下面定义了一个Flag show,show默认为false, 有两种调用方式--show\-s，命令后面接了show则上面定义的show变量就会变成true
	showL = *testCmd.Flags().BoolP("showL", "S", false, "showL")

	// 必须设置某些选项
	// 设置使用test的时候后面必须接show
	_ = testCmd.MarkFlagRequired("showL") // 不写会提示Error: required flag(s) "showLd" not set

}
