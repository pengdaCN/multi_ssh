package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"multi_ssh/m_terminal"
	"multi_ssh/model"
	"sync"
	"time"
)

func init() {
	shellCmd.Flags().BoolVar(&enableSudo, "sudo", false, "设置是否以sudo方式执行命令")
	shellCmd.Flags().StringVar(&saveFile, "save", "", "将shell命令执行的输出保存到本地的文件中")
	rootCmd.AddCommand(&shellCmd)
}

var (
	enableSudo bool
	saveFile   string
)

type commandResult struct {
	u   model.SHHUser
	msg []byte
}

var shellCmd = cobra.Command{
	Use:     "shell command",
	Short:   "执行一行shell命令",
	Example: "shell --sudo true 'command'",
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		ch := make(chan *commandResult, 0)
		go func() {
			if saveFile != "" {
				//	TODO 后续添加保存功能
			}
			for r := range ch {
				fmt.Printf("%s\n\t%s\n", r.u.Host(), string(r.msg))
			}
		}()
		var w sync.WaitGroup
		for _, t := range terminals {
			w.Add(1)
			go func(term *m_terminal.Terminal) {
				defer w.Done()
				bs, err := term.Run2(args[0], enableSudo)
				if err == nil {
					ch <- &commandResult{
						u:   term.GetUser(),
						msg: bs,
					}
				} else {
					ch <- &commandResult{
						u:   term.GetUser(),
						msg: []byte(err.Error()),
					}
				}
			}(t)
		}
		w.Wait()
		time.Sleep(time.Second)
		close(ch)
	},
}
