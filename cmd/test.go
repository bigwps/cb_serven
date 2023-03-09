package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var getCmd = &cobra.Command{
	Use:   "get",
	Short: "待续",
	Long:  `待开发子命令`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("get called")
	},
}

func init() {
	rootCmd.AddCommand(getCmd) //注意这里rootCmd 表示当前这个命令是rootCmd的子命令
	//它会继承rootcmd的所有操作方法，在主体命令下层
}
