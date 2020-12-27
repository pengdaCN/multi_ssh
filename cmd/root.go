package cmd

import (
	"github.com/spf13/cobra"
	"log"
	"multi_ssh/cro"
	"os"
	"time"
)

const version = `0.3.9`

var (
	globalBuilder = cro.New()
)

var (
	hosts        string
	hostLine     string
	outFormat    string
	filterStr    string
	execableNums int
	preInfo      bool
	timeout   time.Duration
)


func init() {
	rootCmd.Flags().StringVarP(&hosts, "hosts", "", "./hosts", "multi_ssh 读取hosts配置文件")
	rootCmd.Flags().StringVarP(&hostLine, "line", "", "", "从cli中读取要连接的信息")
	rootCmd.Flags().StringVarP(&outFormat, "format", "f", cro.DefaultOutputFormat, "以指定格式输出信息")
	rootCmd.Flags().StringVarP(&filterStr, "filter", "F", "", "使用格式选择需要执行的主机")
	rootCmd.Flags().BoolVarP(&preInfo, "uinfo", "", false, "是否在对主机操作之前获取他的信息")
	rootCmd.Flags().DurationVarP(&timeout, "wait", "w", -1, "设置超时，默认不永不超时")
	rootCmd.Flags().IntVarP(&execableNums, "limit-exec", "L", -1, "限制执行连接主机最大个数，默认限制")
}

var rootCmd = cobra.Command{
	Use:              "multi_ssh",
	Short:            "这是一个简单的cli工具",
	Long:             "这是一个简单的cli的并发ssh client工具",
	Version:          version,
	TraverseChildren: true,
	Args:             cobra.MaximumNArgs(0),
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if hostLine != "" {
			globalBuilder.RawHostsInfo(hostLine)
		} else {
			globalBuilder.Hosts(hosts)
		}
		globalBuilder.Format(outFormat).Filter(filterStr).PreInfo(!preInfo).SetMaxExecSeveral(execableNums).Timeout(timeout)
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Println(err.Error())
		os.Exit(1)
	}
}
