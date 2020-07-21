package cmd

import (
	"github.com/spf13/cobra"
	"io"
	"multi_ssh/m_terminal"
	"os"
)

func init() {
	shellCmd.Flags().BoolVar(&shellSudo, "sudo", false, "设置是否以sudo方式执行命令")
	shellCmd.Flags().StringVar(&shellSaveFile, "save", "", "将shell命令执行的输出保存到本地的文件中")
	rootCmd.AddCommand(&shellCmd)
}

var (
	shellSudo     bool
	shellSaveFile string
)

var shellCmd = cobra.Command{
	Use:     "shell command",
	Short:   "执行一行shell命令",
	Example: "shell --sudo 'command'",
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		ch := make(chan *execResult, 0)
		out := []io.Writer{
			os.Stdout,
		}
		if shellSaveFile != "" {
			fil, err := os.Create(scriptSaveFile)
			if err != nil {
				panic(err)
			}
			out = append(out, fil)
		}
		outFinish := output(ch, outFormat, out...)
		execFinish := eachTerm(terminals, func(term *m_terminal.Terminal) {
			rst, err := term.Run(shellSudo, args[0])
			ch <- buildExecResult(term, rst, err)
		})
		<-execFinish
		close(ch)
		<-outFinish
	},
}
