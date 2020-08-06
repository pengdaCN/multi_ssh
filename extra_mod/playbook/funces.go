package playbook

import (
	"github.com/pkg/errors"
	lua "github.com/yuin/gopher-lua"
	"log"
	"os"
)

// shell(id int, sudo bool, cmd string) out -> playbookResult
func shell(state *lua.LState) int {
	id := state.ToInt(1)
	sudo := state.ToBool(2)
	cmd := state.ToString(3)
	term := GetTerm(id)
	rst, err := term.Run(sudo, cmd)
	p := buildPlaybookResult(rst, err)
	p.u = term.GetUser()
	tab := rstToLTable(state, p)
	state.Push(tab)
	return 1
}

// script(id int, sudo bool, script_path string, args string) out -> playbookResult
func script(state *lua.LState) int {
	id := state.ToInt(1)
	sudo := state.ToBool(2)
	scriptPath := state.ToString(3)
	args := state.ToString(4)
	term := GetTerm(id)
	fil, err := os.Open(scriptPath)
	if err != nil {
		log.Println(errors.WithStack(err))
		state.Push(nil)
		return 1
	}
	rst, err := term.Script(sudo, fil, args)
	p := buildPlaybookResult(rst, err)
	p.u = term.GetUser()
	tab := rstToLTable(state, p)
	state.Push(tab)
	return 1
}

// copy(id int, sudo, exists bool, src []string, dst string, attr map<lua table>) out -> playbookResult
func cp(state *lua.LState) int {
	id := state.ToInt(1)
	sudo := state.ToBool(2)
	exists := state.ToBool(3)
	_src := state.ToTable(4)
	src := lTableToStrSlice(_src)
	dst := state.ToString(5)
	_attr := state.ToTable(6)
	attr := buildHandleByFileWithLTable(_attr)
	term := GetTerm(id)
	err := term.Copy(exists, sudo, src, dst, attr)
	p := buildPlaybookResultWithErr(err)
	tab := rstToLTable(state, p)
	state.Push(tab)
	return 1
}

// output(id int, tab map)
// output 函数可以设置格式，在tab中的'format'键定义格式，与#{key}
func out(state *lua.LState) int {

	return 0
}
