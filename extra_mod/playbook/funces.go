package playbook

import (
	"context"
	"fmt"
	lua "github.com/yuin/gopher-lua"
	"io"
	"multi_ssh/m_terminal"
	"multi_ssh/tools"
	"os"
	"strings"
	"time"
)

type (
	FUNCWithTerm func(*m_terminal.Terminal) lua.LGFunction
)

func luaSleep(state *lua.LState) int {
	second := state.ToInt(1)
	time.Sleep(time.Duration(second) * time.Second)
	return 0
}

func newShell(term *m_terminal.Terminal) lua.LGFunction {
	return func(state *lua.LState) int {
		var (
			sudo = false
			cmd  string
			rst  *m_terminal.Result
		)
		args := state.ToTable(1)
		cmd = lvalueToStr(args.RawGetInt(1))
		sudo = lvalueToBool(args.RawGetString("sudo"))
		ctx := useTimeoutFromLvalue(args.RawGetString("timeout"))
		tools.WithCancel(ctx, func() {
			rst = term.Run(sudo, cmd)
		})
		if rst == nil {
			rst = new(m_terminal.Result)
		}
		state.Push(rstToLTable(state, rst))
		return 1
	}
}

func newScript(term *m_terminal.Terminal) lua.LGFunction {
	return func(state *lua.LState) int {
		var (
			sudo = false
			args string
			path string
			text string
			read io.Reader
			rst  *m_terminal.Result
		)
		part := state.ToTable(1)
		sudo = lvalueToBool(part.RawGetString("sudo"))
		args = lvalueToStr(part.RawGetString("args"))
		// 若有text选项，则优先使用text选项内容作为脚本执行
		text = lvalueToStr(part.RawGetString("text"))
		if text == "" {
			path = lvalueToStr(part.RawGetInt(1))
			if path == "" {
				panic("required one path")
			}
			var err error
			read, err = os.Open(path)
			if err != nil {
				panic("script not exists")
			}
		} else {
			read = strings.NewReader(text)
		}
		ctx := useTimeoutFromLvalue(part.RawGetString("timeout"))
		tools.WithCancel(ctx, func() {
			rst = term.Script(sudo, read, args)
		})
		if rst == nil {
			rst = new(m_terminal.Result)
		}
		state.Push(rstToLTable(state, rst))
		return 1
	}
}

func newContext(term *m_terminal.Terminal) lua.LGFunction {
	return func(state *lua.LState) int {
		var (
			filename string
			text     string
			read     io.Reader
			rst      *m_terminal.Result
		)
		part := state.ToTable(1)
		filename = lvalueToStr(part.RawGetString("filename"))
		// 若有text选项，则优先使用text选项内容作为脚本执行
		text = lvalueToStr(part.RawGetString("text"))

		read = strings.NewReader(text)

		ctx := useTimeoutFromLvalue(part.RawGetString("timeout"))
		tools.WithCancel(ctx, func() {
			rst = m_terminal.BuildRstByErr(term.SftpUpdateByReaderWithFunc(filename, read, `./`))
		})
		if rst == nil {
			rst = new(m_terminal.Result)
		}
		state.Push(rstToLTable(state, rst))
		return 1
	}
}

func newCopy(term *m_terminal.Terminal) lua.LGFunction {
	return func(state *lua.LState) int {
		var (
			sudo   = false
			exists = false
			src    []string
			dst    string
			rst    *m_terminal.Result
		)
		args := state.ToTable(1)
		sudo = lvalueToBool(args.RawGetString("sudo"))
		exists = lvalueToBool(args.RawGetString("exists"))
		{
			s := args.RawGetInt(1)
			switch s.Type() {
			case lua.LTTable:
				t := s.(*lua.LTable)
				src = lTableToStrSlice(t)
			case lua.LTString:
				src = []string{
					s.String(),
				}
			default:
				panic("required src path")
			}
		}
		{
			dst = args.RawGetInt(2).String()
			if dst == "" {
				panic("required dst path")
			}
		}
		//{
		//	t := args.RawGetString("attr")
		//	switch t.Type() {
		//	case lua.LTNil:
		//		attr = nil
		//	case lua.LTTable:
		//		tab := t.(*lua.LTable)
		//		attr = buildHandleByFileWithLTable(tab)
		//	default:
		//		attr = nil
		//	}
		//}
		ctx := useTimeoutFromLvalue(args.RawGetString("timeout"))
		tools.WithCancel(ctx, func() {
			rst = term.Copy(exists, sudo, src, dst)
		})
		if rst == nil {
			rst = new(m_terminal.Result)
		}
		state.Push(rstToLTable(state, rst))
		return 1
	}
}

func newExtra(term *m_terminal.Terminal) lua.LGFunction {
	return func(state *lua.LState) int {
		info := term.GetUser().Extra()
		if info == nil {
			state.Push(state.NewTable())
			return 1
		}
		state.Push(mapToLTable(state, info))
		return 1
	}
}

func newHostInfo(term *m_terminal.Terminal) lua.LGFunction {
	return func(state *lua.LState) int {
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
}

const (
	OutKey = "output"
)

func newOut(term *m_terminal.Terminal) lua.LGFunction {
	return func(state *lua.LState) int {
		msg := state.ToString(1)
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
}

func newOutLn(term *m_terminal.Terminal) lua.LGFunction {
	return func(state *lua.LState) int {
		msg := state.ToString(1)
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
}

const (
	Code = "code"
)

// 用于设置脚本结束状态码，模式是0，即为正常
func newSetCode(term *m_terminal.Terminal) lua.LGFunction {
	return func(state *lua.LState) int {
		code := state.ToInt(1)
		term.SetShare(Code, code)
		return 0
	}
}

const (
	ErrInfo = "errinfo"
)

// 设置错误状态码与错误信息
func newSetErrInfo(term *m_terminal.Terminal) lua.LGFunction {
	return func(state *lua.LState) int {
		code := state.ToInt(1)
		errInfo := state.ToString(2)
		term.SetShare(Code, code)
		term.SetShare(ErrInfo, errInfo)
		return 0
	}
}

func newExit(cancel context.CancelFunc) lua.LGFunction {
	return func(state *lua.LState) int {
		cancel()
		return 0
	}
}
