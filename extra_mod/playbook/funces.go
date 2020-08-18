package playbook

import (
	"github.com/pkg/errors"
	lua "github.com/yuin/gopher-lua"
	"log"
	"os"
)

// shell(id int, sudo bool, cmd string) out -> playbookResult
func shell(state *lua.LState) int {
	genericFuncSendRst(state, func(state *lua.LState) *playbookResult {
		id := state.ToInt(1)
		sudo := state.ToBool(2)
		cmd := state.ToString(3)
		term, ok := Get(id)
		if !ok {
			log.Println("错误，执行lua shell 执行时，错误的id")
			return buildPlaybookWithCode(1)
		}
		rst, err := term.Run(sudo, cmd)
		p := buildPlaybookResult(rst, err)
		p.u = term.GetUser()
		return p
	})
	return 1
}

// script(id int, sudo bool, script_path string, args string) out -> playbookResult
func script(state *lua.LState) int {
	genericFuncSendRst(state, func(state *lua.LState) *playbookResult {
		id := state.ToInt(1)
		sudo := state.ToBool(2)
		scriptPath := state.ToString(3)
		args := state.ToString(4)
		term, ok := Get(id)
		if !ok {
			log.Println("错误，执行lua script 执行时，错误的id")
			return buildPlaybookWithCode(1)
		}
		fil, err := os.Open(scriptPath)
		if err != nil {
			log.Println(errors.WithStack(err))
			return buildPlaybookWithCode(1)
		}
		rst, err := term.Script(sudo, fil, args)
		p := buildPlaybookResult(rst, err)
		p.u = term.GetUser()
		return p
	})
	return 1
}

// copy(id int, sudo, exists bool, src []string, dst string, attr map<lua table>) out -> playbookResult
func cp(state *lua.LState) int {
	genericFuncSendRst(state, func(state *lua.LState) *playbookResult {
		id := state.ToInt(1)
		sudo := state.ToBool(2)
		exists := state.ToBool(3)
		_src := state.ToTable(4)
		src := lTableToStrSlice(_src)
		dst := state.ToString(5)
		_attr := state.ToTable(6)
		attr := buildHandleByFileWithLTable(_attr)
		term, ok := Get(id)
		if !ok {
			log.Printf("错误，执行lua copy 执行时，错误的id: %d\n", id)
			return buildPlaybookWithCode(1)
		}
		err := term.Copy(exists, sudo, src, dst, attr)
		p := buildPlaybookResultWithErr(err)
		return p
	})
	return 1
}

// output(id int, format, msg string)
// output 函数可以设置格式，在tab中的'format'键定义格式，与#{key}
func out(state *lua.LState) int {

	return 0
}

// 通用处理方法，回自动发送返回值
func genericFuncSendRst(state *lua.LState, fn func(*lua.LState) *playbookResult) {
	rst := fn(state)
	state.Push(rstToLTable(state, rst))
}
