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

// showCmd represents the show command
var showCmd = &cobra.Command{
	Use:   "show",
	Short: "show简短介绍展示当前时间",
	Long:  `show详细介绍`,
	Run:   ShowTime,
}

// ShowTime 显示当前时间
func ShowTime(cmd *cobra.Command, args []string) {
	fmt.Printf("cmd是%+v,args是%+v\n", *cmd, args)
	fmt.Println("现在时间是", time.Now())
}

func init() {
	rootCmd.AddCommand(showCmd)

}

/*
命令：  ./time show hello
结果：
	cmd是{Use:show Aliases:[] SuggestFor:[] Short:show简短介绍展示当前时间 Long:show详细介绍 Example: ValidArgs:[]
	Args:<nil> ArgAliases:[] BashCompletionFunction: Deprecated: Hidden:false Annotations:map[]
	Version: PersistentPreRun:<nil> PersistentPreRunE:<nil> PreRun:<nil> PreRunE:<nil> Run:0x14238f0
	RunE:<nil> PstRunE:<nil> PersistentPostRun:<nil> PersistentPostRunE:<nil> SilenceErrors:false
	SilenceUsage:false DisableFlagParsing:false DisableAutoGenTag:false DisableFlagsInUseLine:false DisableSuggestions:false
	SuggestionsMinimumDistance:0 TraverseChildren:false FParseErrWhitelist:{UnknownFlags:false} commands:[0x18bb360]
	parent:0x18bb5c0 commandsMaxUseLen:5 commandsMaxCommandPathLen:10 commandsMaxNameLen:5 commandsAreSorted:false
	commandCalledAs:{name:show called:true} args:[] flagErrorBuf:0xc0000219e0 flags:0xc0000bea00
	pflags:0xc0000beb00 lflags:<nil> iflags:<nil> parentsPflags:0xc0000be900 globNormFunc:<nil> output:<nil>
	usageFunc:<nil> usageTemplate: flagErrorFunc:<nil> helpTemplate: helpFunc:<nil> helpCommand:<nil> versionTemplate:},

	args是[hello]

	现在时间是 2021-12-01 12:45:17.122819 +0800 CST m=+0.017287141

*/
