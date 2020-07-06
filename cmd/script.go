package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"log"
	"multi_ssh/m_terminal"
	"path"
	"sync"
	"time"
)

var (
	scriptSudo bool
)

func init() {
	rootCmd.AddCommand(&scriptCmd)
	scriptCmd.Flags().BoolVar(&scriptSudo, "sudo", false, "是否以sudo方式执行脚本")
}

var scriptCmd = cobra.Command{
	Use:   "script file",
	Short: "将本地脚本上传到远端并执行",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		chTerm := make(chan *m_terminal.Terminal, 0)
		var w sync.WaitGroup
		go func() {
			for _, t := range terminals {
				w.Add(1)
				go func(term *m_terminal.Terminal) {
					defer w.Done()
					err := term.SftpUpdate(args[0], "/tmp", nil)
					if err != nil {
						log.Println("脚本上传失败", err)
						return
					}
					chTerm <- term
				}(t)
			}
		}()
		var w2 sync.WaitGroup
		chRst := make(chan *commandResult, 0)
		go func() {
			for r := range chRst {
				fmt.Printf("%s\n\t%s\n", r.u.Host(), string(r.msg))
			}
		}()
		go func() {
			for t := range chTerm {
				w2.Add(1)
				go func(term *m_terminal.Terminal) {
					defer w2.Done()
					bs, err := term.Run(scriptSudo, fmt.Sprintf("bash /tmp/%s", path.Base(args[0])))
					if err == nil {
						chRst <- &commandResult{
							u:   term.GetUser(),
							msg: bs,
						}
					} else {
						chRst <- &commandResult{
							u:   term.GetUser(),
							msg: []byte(err.Error()),
						}
					}
					_ = term.Remove(path.Join("/", "tmp", path.Base(args[0])))
				}(t)
			}
		}()
		w.Wait()
		time.Sleep(time.Second)
		close(chTerm)
		w2.Wait()
		time.Sleep(time.Second)
		close(chRst)
	},
}
