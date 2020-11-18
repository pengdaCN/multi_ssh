package cmd

import (
	"github.com/spf13/cobra"
	"log"
	"multi_ssh/m_terminal"
	"multi_ssh/model"
	"os"
	"sort"
	"sync"
	"time"
)

const version = `0.3.5`

type (
	userSlice []model.SHHUser
)

var (
	hosts        string
	hostLine     string
	outFormat    string
	filterStr    string
	execableNums int
	preInfo      bool
)

var (
	terminals []*m_terminal.Terminal
	users     userSlice
	timeout   time.Duration
)

func init() {
	rootCmd.Flags().StringVarP(&hosts, "hosts", "", "./hosts", "multi_ssh 读取hosts配置文件")
	rootCmd.Flags().StringVarP(&hostLine, "line", "", "", "从cli中读取要连接的信息")
	rootCmd.Flags().StringVarP(&outFormat, "format", "f", defaultOutputFormat, "以指定格式输出信息")
	rootCmd.Flags().StringVarP(&filterStr, "filter", "F", "", "使用格式选择需要执行的主机")
	rootCmd.Flags().BoolVarP(&preInfo, "uinfo", "", true, "是否在对主机操作之前获取他的信息")
	rootCmd.Flags().DurationVarP(&timeout, "wait", "w", -1, "设置超时，默认不永不超时")
	rootCmd.Flags().IntVarP(&execableNums, "limit-exec", "L", -1, "限制执行连接主机最大个数，默认限制")
}

func (u userSlice) Less(v1, v2 int) bool {
	return u[v1].Line() < u[v2].Line()
}

func (u userSlice) Swap(v1, v2 int) {
	u[v1], u[v2] = u[v2], u[v1]
}

func (u userSlice) Len() int {
	return len(u)
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
			u := model.ReadLine(hostLine)
			if u == nil {
				log.Fatalln("ERROR 使用命令行参数错误")
			}
			users = append(users, u)
		} else {
			us, err := model.ReadHosts(hosts)
			if err != nil {
				log.Println(err.Error())
				os.Exit(1)
			}
			for _, u := range us {
				users = append(users, u)
			}
		}
		if filterStr != "" {
			users = filters(users, filterStr)
		}
		// 使用行号进行排序
		sort.Sort(users)
		ch := make(chan *m_terminal.Terminal, 0)
		var w sync.WaitGroup
		for i, u := range users {
			w.Add(1)
			go func(user model.SHHUser, bi int) {
				defer w.Done()
				c, err := m_terminal.DefaultWithPassphrase(user)
				if err != nil {
					log.Printf("打开%s失败 %s", user.Host(), err)
					return
				} else {
					log.Printf("打开%s成功", user.Host())
				}
				if preInfo {
					m_terminal.GetRemoteHostInfo(c)
				}
				c.SetBirthID(bi + 1)
				ch <- c
			}(u, i)
		}
		var w2 sync.WaitGroup
		w2.Add(1)
		go func() {
			defer w2.Done()
			for i := range ch {
				terminals = append(terminals, i)
			}
		}()
		w.Wait()
		close(ch)
		w2.Wait()
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Println(err.Error())
		os.Exit(1)
	}
}
