package cmd

import (
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	lua "github.com/yuin/gopher-lua"
	"log"
	"multi_ssh/extra_mod/playbook"
	"multi_ssh/m_terminal"
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
			_, _, _ = playbook.VM.Resume(co, fn, lua.LNumber(term.GetID()))

		})
		<-finished
	},
}
