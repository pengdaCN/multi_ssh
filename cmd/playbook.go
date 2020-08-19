package cmd

import (
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	lua "github.com/yuin/gopher-lua"
	"log"
	"multi_ssh/extra_mod/playbook"
	"multi_ssh/m_terminal"
	"os"
	"strings"
)

func init() {
	rootCmd.AddCommand(&playbookCmd)
}

var playbookCmd = cobra.Command{
	Use:     "playbook <file>",
	Short:   "执行一系列的命令",
	Long:    "通过golang内置的lua虚拟机来执行一系列操作",
	Args:    cobra.MinimumNArgs(1),
	Example: "playbook example.play",
	Run: func(cmd *cobra.Command, args []string) {
		ch := make(chan *execResult, 0)
		outFinish := output(ch, outFormat, os.Stdout)
		if err := playbook.VM.DoFile(args[0]); err != nil {
			log.Println(errors.New("错误文件位置"))
			return
		}
		var (
			fn *lua.LFunction
			ok bool
		)
		if fn, ok = playbook.VM.GetGlobal("exec").(*lua.LFunction); !ok {
			log.Println(errors.New("未读取到exec函数，请检查代码"))
			return
		}
		finished := eachTerm(terminals, func(term *m_terminal.Terminal) {
			playbook.Push(term.GetID(), term)
			co, _ := playbook.VM.NewThread()
			_, err, _ := playbook.VM.Resume(co, fn, lua.LNumber(term.GetID()))
			var (
				msg     string
				code    int
				errInfo string
			)
			if err != nil {
				log.Println(err.Error())
			}
			out, ok := term.GetOnceShare(playbook.OutKey)
			if ok {
				sb := out.(*strings.Builder)
				str := sb.String()
				msg = str
			}
			c, ok := term.GetOnceShare(playbook.Code)
			if ok {
				code, _ = c.(int)
			}
			_errInfo, ok := term.GetOnceShare(playbook.ErrInfo)
			if ok {
				errInfo, _ = _errInfo.(string)
			}
			rst := new(execResult)
			{
				rst.errInfo = errInfo
				rst.msg = msg
				rst.code = code
				rst.u = term.GetUser()
			}
			ch <- rst
		})
		<-finished
		close(ch)
		<-outFinish
	},
}
