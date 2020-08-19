package playbook

import (
	"fmt"
	"github.com/pkg/errors"
	lua "github.com/yuin/gopher-lua"
	"log"
	"multi_ssh/m_terminal"
	"os"
	"strings"
)

// shell(id int, sudo bool, cmd string) out -> playbookResult
func shell(state *lua.LState) int {
	genericFuncSendRst(state, func(state *lua.LState) *m_terminal.Result {
		id := state.ToInt(1)
		sudo := state.ToBool(2)
		cmd := state.ToString(3)
		term, ok := Get(id)
		if !ok {
			log.Println("错误，执行lua shell 执行时，错误的id")
			return m_terminal.BuildRstWithCode(1)
		}
		rst := term.Run(sudo, cmd)
		return rst
	})
	return 1
}

// script(id int, sudo bool, script_path string, args string) out -> playbookResult
func script(state *lua.LState) int {
	genericFuncSendRst(state, func(state *lua.LState) *m_terminal.Result {
		id := state.ToInt(1)
		sudo := state.ToBool(2)
		scriptPath := state.ToString(3)
		args := state.ToString(4)
		term, ok := Get(id)
		if !ok {
			log.Println("错误，执行lua script 执行时，错误的id")
			return m_terminal.BuildRstWithCode(1)
		}
		fil, err := os.Open(scriptPath)
		if err != nil {
			log.Println(errors.WithStack(err))
			return m_terminal.BuildRstWithCode(1)

		}
		rst := term.Script(sudo, fil, args)
		return rst
	})
	return 1
}

// copy(id int, sudo, exists bool, src []string, dst string, attr map<lua table>) out -> playbookResult
func cp(state *lua.LState) int {
	genericFuncSendRst(state, func(state *lua.LState) *m_terminal.Result {
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
			return m_terminal.BuildRstWithCode(1)
		}
		rst := term.Copy(exists, sudo, src, dst, attr)

		return rst
	})
	return 1
}

// extraInfo(id int) map
// extraInfo 函数返回主机在hosts文件中的扩展信息，以map形式返回
func extraInfo(state *lua.LState) int {
	id := state.ToInt(1)
	term, ok := Get(id)
	if !ok {
		log.Printf("错误，执行lua extraInfo 执行时，错误的id: %d\n", id)
		state.Push(state.NewTable())
		return 1
	}
	info := term.GetUser().Extra()
	if info == nil {
		state.Push(state.NewTable())
		return 1
	}
	state.Push(mapToLTable(state, info))
	return 1
}

func hostInfo(state *lua.LState) int {
	id := state.ToInt(1)
	term, ok := Get(id)
	if !ok {
		log.Printf("错误，执行lua hostInfo 执行时，错误的id: %d\n", id)
		state.Push(state.NewTable())
		return 1
	}
	info := term.GetUser()
	if info == nil {
		state.Push(state.NewTable())
		return 1
	}
	host := info.Host()
	username := info.User()
	m := map[string]string{
		"host": host,
		"user": username,
	}
	state.Push(mapToLTable(state, m))
	return 1
}

const (
	OutKey = "output"
)

// out(id int, msg string)
// out 函数可以设置格式，在tab中的'format'键定义格式，与#{key}
func out(state *lua.LState) int {
	id := state.ToInt(1)
	msg := state.ToString(2)
	term, ok := Get(id)
	if !ok {
		log.Printf("错误，执行lua out 执行时，错误的id: %d\n", id)
		return 0
	}
	var b *strings.Builder
	v, ok := term.GetOnceShare(OutKey)
	if !ok {
		b = new(strings.Builder)
		term.SetShare(OutKey, b)
	} else {
		b = v.(*strings.Builder)
	}
	b.WriteString(msg)
	return 0
}

// outln(id int, msg string)
// outln 函数可以设置格式，在tab中的'format'键定义格式，与#{key}
func outln(state *lua.LState) int {
	id := state.ToInt(1)
	msg := state.ToString(2)
	term, ok := Get(id)
	if !ok {
		log.Printf("错误，执行lua out 执行时，错误的id: %d\n", id)
		return 0
	}
	var b *strings.Builder
	v, ok := term.GetOnceShare(OutKey)
	if !ok {
		b = new(strings.Builder)
		term.SetShare(OutKey, b)
	} else {
		b = v.(*strings.Builder)
	}
	b.WriteString(fmt.Sprintln(msg))
	return 0
}

const (
	Code = "code"
)

// 用于设置脚本结束状态码，模式是0，即为正常
// setCode(id int, code int)
func setCode(state *lua.LState) int {
	id := state.ToInt(1)
	code := state.ToInt(2)
	term, ok := Get(id)
	if !ok {
		log.Printf("错误，执行lua setCode 执行时，错误的id: %d\n", id)
		return 0
	}
	term.SetShare(Code, code)
	return 0
}

const (
	ErrInfo = "errinfo"
)

// 设置错误状态码与错误信息
// setErrInfo(id, code int, err string)
func setErrInfo(state *lua.LState) int {
	id := state.ToInt(1)
	code := state.ToInt(2)
	errInfo := state.ToString(3)
	term, ok := Get(id)
	if !ok {
		log.Printf("错误，执行lua setErrInfo 执行时，错误的id: %d\n", id)
		return 0
	}
	term.SetShare(Code, code)
	term.SetShare(ErrInfo, errInfo)
	return 0
}

// 通用处理方法，回自动发送返回值
func genericFuncSendRst(state *lua.LState, fn func(*lua.LState) *m_terminal.Result) {
	rst := fn(state)
	state.Push(rstToLTable(state, rst))
}
