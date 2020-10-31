package cmd

import (
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	lua "github.com/yuin/gopher-lua"
	"log"
	"multi_ssh/common"
	"multi_ssh/extra_mod/playbook"
	"multi_ssh/m_terminal"
	"os"
	"regexp"
	"strings"
)

var (
	argsList string
)

func init() {
	rootCmd.AddCommand(&playbookCmd)
	playbookCmd.Flags().StringVarP(&argsList, "set-args", "S", "", "设置lua中全局的变量，key=val形式")
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
			log.Println(errors.WithStack(err))
			return
		}
		if argsList != "" {
			setGlobalVal(argsList)
		}
		var (
			fn *lua.LFunction
			ok bool
		)
		_ = playbook.VM.CallByParam(lua.P{
			Fn:      playbook.VM.GetGlobal("BEGIN"),
			NRet:    0,
			Protect: true,
		})
		if fn, ok = playbook.VM.GetGlobal("exec").(*lua.LFunction); !ok {
			log.Println(errors.New("未读取到exec函数，请检查代码"))
			return
		}
		finished := eachTerm(terminals, func(term *m_terminal.Terminal) {
			co, cancel := playbook.VM.NewThread()
			t := playbook.NewLuaTerm(co, term, cancel)
			_ = playbook.VM.CallByParam(lua.P{
				Fn:      playbook.VM.GetGlobal("EXEC_BEGIN"),
				NRet:    0,
				Protect: true,
			}, t)
			_, err, _ := playbook.VM.Resume(co, fn, t)
			_ = playbook.VM.CallByParam(lua.P{
				Fn:      playbook.VM.GetGlobal("EXEC_OVER"),
				NRet:    0,
				Protect: true,
			}, t)
			var (
				msg     string
				code    int
				errInfo string
			)
			if err != nil {
				log.Println("VM: ", err.Error())
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
		_ = playbook.VM.CallByParam(lua.P{
			Fn:      playbook.VM.GetGlobal("OVER"),
			NRet:    0,
			Protect: true,
		})
		<-finished
		close(ch)
		<-outFinish
	},
}

var (
	spaceLine    = regexp.MustCompile(`^\s*$`)
	assignment   = regexp.MustCompile(`^\s*=`)
	segmentation = regexp.MustCompile(`^\s*,`)
)

func setGlobalVal(str string) {
	var (
		word string
		val  string
	)
	for {
		if spaceLine.MatchString(str) {
			break
		}
		word, str = common.ReadWord(str)
		if !assignment.MatchString(str) {
			panic("ERROR format")
		}
		str = str[strings.IndexRune(str, '='):]
		val, str = common.ReadStr(str)
		playbook.SetGlobalVal(word, val)
		if spaceLine.MatchString(str) {
			break
		}
		if !segmentation.MatchString(str) {
			panic("ERROR format")
		}
	}
}
