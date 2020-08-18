package playbook

import (
	"github.com/pkg/sftp"
	lua "github.com/yuin/gopher-lua"
	"golang.org/x/crypto/ssh"
	"io"
	"multi_ssh/m_terminal"
	"multi_ssh/model"
	"multi_ssh/tools"
	"strconv"
)

// 不知道的错误，在执行bash脚本是方式不知道的错误
const UNKNOWNErr = -1

type playbookResult struct {
	u       model.SHHUser
	msg     []byte
	errInfo string
	code    int
}

func buildPlaybookResult(rst []byte, err error) *playbookResult {
	p := new(playbookResult)
	p.msg = rst
	var (
		code    int
		errInfo string
	)
	{
		if _err, ok := err.(*ssh.ExitError); ok {
			code = _err.ExitStatus()
		} else if err == nil {
			code = 0
		} else {
			errInfo = err.Error()
			code = UNKNOWNErr
		}
	}
	p.code = code
	p.errInfo = errInfo
	return p
}

func buildPlaybookResultWithErr(err error) *playbookResult {
	p := buildPlaybookResult(nil, err)
	if err == nil {
		p.msg = []byte("OK")
	}
	return p
}

func buildPlaybookWithCode(code int) *playbookResult {
	p := new(playbookResult)
	p.code = code
	return p
}

func rstToLTable(state *lua.LState, p *playbookResult) lua.LValue {
	tab := state.NewTable()
	state.SetTable(tab, lua.LString("u"), sshUserToLTable(state, p.u))
	state.SetTable(tab, lua.LString("msg"), lua.LString(p.msg))
	state.SetTable(tab, lua.LString("errInfo"), lua.LString(p.errInfo))
	state.SetTable(tab, lua.LString("code"), lua.LNumber(p.code))
	return tab
}

func sshUserToLTable(state *lua.LState, u model.SHHUser) lua.LValue {
	tab := state.NewTable()
	state.SetTable(tab, lua.LString("user"), lua.LString(u.User()))
	state.SetTable(tab, lua.LString("host"), lua.LString(u.Host()))
	return tab
}

type _output struct {
	receive []io.WriteCloser
}

func newOutput(closer ...io.WriteCloser) _output {
	return _output{
		receive: closer,
	}
}

func (o *_output) append(out ...io.WriteCloser) {
	o.receive = append(o.receive, out...)
}

func lTableToStrSlice(src *lua.LTable) []string {
	rst := make([]string, 0, src.Len())
	src.ForEach(func(key lua.LValue, val lua.LValue) {
		rst = append(rst, val.String())
	})
	return rst
}

func lTableToMapStrLValue(src *lua.LTable) map[string]lua.LValue {
	rst := make(map[string]lua.LValue)
	src.ForEach(func(key lua.LValue, val lua.LValue) {
		rst[key.String()] = val
	})
	return rst
}

func buildHandleByFileWithLTable(src *lua.LTable) m_terminal.HandleByFile {
	m := lTableToMapStrLValue(src)
	var (
		uid  int
		gid  int
		mode string
	)
	if u, ok := m["uid"]; ok {
		uid, _ = strconv.Atoi(u.String())
	} else {
		uid = -1
	}

	if g, ok := m["gid"]; ok {
		gid, _ = strconv.Atoi(g.String())
	} else {
		gid = -1
	}
	if m, ok := m["mode"]; ok {
		mode = m.String()
	}

	return func(file *sftp.File) error {
		if mode != "" {
			m, err := tools.String2FileMode(mode)
			if err != nil {
				return err
			}
			if err := file.Chmod(m); err != nil {
				return err
			}
		}
		if uid != -1 && gid != -1 {
			if err := file.Chown(uid, gid); err != nil {
				return err
			}
		}
		return nil
	}
}
