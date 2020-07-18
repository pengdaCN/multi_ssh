package cmd

import (
	"github.com/spf13/cobra"
	"io"
	"multi_ssh/m_terminal"
	"os"
	"sync"
)

func init() {
	shellCmd.Flags().BoolVar(&enableSudo, "sudo", false, "设置是否以sudo方式执行命令")
	shellCmd.Flags().StringVar(&shellSaveFile, "save", "", "将shell命令执行的输出保存到本地的文件中")
	rootCmd.AddCommand(&shellCmd)
}

var (
	enableSudo    bool
	shellSaveFile string
)

var shellCmd = cobra.Command{
	Use:     "shell command",
	Short:   "执行一行shell命令",
	Example: "shell --sudo true 'command'",
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		ch := make(chan *commandResult, 0)
		out := []io.Writer{
			os.Stdout,
		}
		if shellSaveFile != "" {
			fil, err := os.Open(scriptSaveFile)
			if err != nil {
				panic(err)
			}
			out = append(out, fil)
		}
		finish := output(ch, outFormat, out...)
		var w sync.WaitGroup
		for _, t := range terminals {
			w.Add(1)
			go func(term *m_terminal.Terminal) {
				defer w.Done()
				bs, err := term.Run(enableSudo, args[0])
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
		close(ch)
		<-finish
	},
}
