package cmd

import (
	"errors"
	"github.com/spf13/cobra"
	"log"
	"multi_ssh/m_terminal"
	"multi_ssh/model"
	"os"
	"sync"
)

var (
	hosts     string
	hostLine  string
	outFormat string
	preInfo   bool
)

var (
	terminals []*m_terminal.Terminal
	users     []model.SHHUser
)

func init() {
	rootCmd.Flags().StringVarP(&hosts, "hosts", "", "./hosts", "multi_ssh 读取hosts配置文件")
	rootCmd.Flags().StringVarP(&hostLine, "line", "", "", "从cli中读取要连接的信息")
	rootCmd.Flags().StringVarP(&outFormat, "format", "f", defaultOutputFormat, "以指定格式输出信息")
	rootCmd.Flags().BoolVarP(&preInfo, "uinfo", "", true, "是否在对主机操作之前获取他的信息")
}

var rootCmd = cobra.Command{
	Use:              "multi_ssh",
	Short:            "这是一个简单的cli工具",
	Long:             "这是一个简单的cli的并发ssh client工具",
	TraverseChildren: true,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 0 {
			return errors.New("错误的位置参数")
		}
		return nil
	},
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if hostLine != "" {
			u, err := model.NewSSHUserByPassphraseWithStringLine(hostLine)
			if err != nil {
				log.Fatalln("ERROR 使用命令行参数错误", err)
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
		ch := make(chan *m_terminal.Terminal, 0)
		var w sync.WaitGroup
		for _, u := range users {
			w.Add(1)
			go func(user model.SHHUser) {
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
				ch <- c
			}(u)
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
