/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"net"
	"os"
	"strings"
	"xc7/bp"
	"xc7/log"
	"xc7/utils"

	"github.com/spf13/cobra"
)

var (
	filename  string
	targets   []string
	bptargets []string
	target    string
	pocfiles  []string
	pocname   string
	report    string
	proxy     string
	gorou     int64
	scan      string
	boom      string
	output    bool
)

// rootCmd represents the base command when called without any subcommands

var rootCmd = &cobra.Command{
	Use:   "CB",
	Short: "CB【漏洞扫描器】",
	Long: `CB【漏洞扫描器】
	V1.0 待加poc版
	`,

	Run: func(cmd *cobra.Command, args []string) {

		// 初始化log
		log.Logger_Init(output)
		log.Xc7Log2()
		// ceye := utils.LoadConfig()

		// 获取目标
		switch {
		case boom != "":
			log.WarLog("爆破模式支持:mssql,oracle,mysql,postgre")
			host, port, err := net.SplitHostPort(boom)

			if err != nil {
				log.ErrorLog("请输入目标:端口")
			} else {
				boom = host + ":" + port
				if pocname != "" {
					bp.CheckWeak(boom, pocname, gorou)

				} else {
					log.ErrorLog("请输入要爆破的数据库或协议")
				}
			}

			// if strings.Contains(boom, "/") {
			// 	bptargets = bp.Parseip(boom)

			// } else {
			// 	bptargets = append(bptargets, boom)

			// }

		case filename != "":
			targets = utils.ReadLine(filename)
		case target != "":
			if strings.HasPrefix(target, "http://") || strings.HasPrefix(target, "https://") {
			} else {
				target = "http://" + target
			}
			targets = append(targets, target)

		case filename == "" && target == "":
			log.ErrorLog("请输入目标 -t || -f")

		}
		if len(targets) > 0 {
			// 加载poc
			switch {
			case pocname != "":
				pocfiles = utils.LoadPoc(pocname)
			case pocname == "":
				pocfiles, _ = utils.LoadAllPoc()
			}

			// 解析ymal
			ymal, err := utils.ParseYml(pocfiles)
			if err != nil {
				log.ErrorLog("解析yaml出现了错误!")
			}

			if ymal != nil {
				//执行poc
				utils.CheckVuls(ymal, targets, gorou, proxy)
			} else {
				log.ErrorLog("为解析到poc,或不支持此格式的poc")
			}
		}

	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	/*
		 *Tips:
			如果在这里定义,对于你的应用程序来说是全局的
	*/

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.xc7.yaml)")
	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	// rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	// rootCmd.Flags().StringVar(&Print, "print", "test", "print test")

	rootCmd.Flags().StringVarP(&filename, "file", "f", "", "从文件读取目标:-f targets.txt")
	rootCmd.Flags().StringVarP(&target, "target", "t", "", "目标:-t http://127.0.0.1 不要多加/")
	rootCmd.Flags().StringVarP(&pocname, "pocname", "p", "", "只测试某应用漏洞:-p tomcat*spring 以*间隔")
	rootCmd.Flags().StringVarP(&boom, "boom", "b", "", "爆破弱口令:-b 192.168.0.1:1443 -p mssql,爆破模式不进行web扫描,不从文件读取目标")
	rootCmd.Flags().StringVarP(&proxy, "proxy", "x", "", "WEB扫描代理:--proxy http://127.0.0.1:10809")
	rootCmd.Flags().Int64VarP(&gorou, "gorou", "g", 1, "设置协程数:-g int")
	rootCmd.Flags().BoolVarP(&output, "output", "o", false, "关闭控制台info日志:-o")

}
