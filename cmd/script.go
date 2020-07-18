package cmd

import (
	"bytes"
	"github.com/spf13/cobra"
	"io"
	"io/ioutil"
	"multi_ssh/m_terminal"
	"os"
	"sync"
)

var (
	scriptSudo     bool
	scriptArgs     string
	scriptSaveFile string
)

func init() {
	rootCmd.AddCommand(&scriptCmd)
	scriptCmd.Flags().BoolVar(&scriptSudo, "sudo", false, "是否以sudo方式执行脚本")
	scriptCmd.Flags().StringVar(&scriptArgs, "args", "", "添加脚本执行的参数")
	scriptCmd.Flags().StringVar(&scriptSaveFile, "save", "", "将脚本输出保存到文件中")
}

var scriptCmd = cobra.Command{
	Use:   "script file",
	Short: "将本地脚本上传到远端并执行",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		ch := make(chan *execResult, 0)
		out := []io.Writer{
			os.Stdout,
		}
		if scriptSaveFile != "" {
			fil, err := os.Create(scriptSaveFile)
			if err != nil {
				panic(err)
			}
			out = append(out, fil)
		}
		finish := output(ch, outFormat, out...)
		scriptContext, err := ioutil.ReadFile(args[0])
		if err != nil {
			panic(err)
		}
		var w sync.WaitGroup
		for _, v := range terminals {
			w.Add(1)
			go func(term *m_terminal.Terminal) {
				defer w.Done()
				rst, err := term.Script(copySudo, bytes.NewReader(scriptContext), scriptArgs)
				ch <- buildExecResult(term, rst, err)
			}(v)
		}
		w.Wait()
		close(ch)
		<-finish
	},
}
