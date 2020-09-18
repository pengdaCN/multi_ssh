package playbook

import (
	"github.com/pkg/sftp"
	lua "github.com/yuin/gopher-lua"
	"multi_ssh/m_terminal"
	"multi_ssh/tools"
	"strconv"
	"strings"
)

func lvalueToBool(value lua.LValue) bool {
	switch value.Type() {
	case lua.LTNil:
		return false
	case lua.LTString:
		s := value.(lua.LString).String()
		s = strings.ToLower(s)
		if s == "" || s == "false" {
			return false
		} else if s == "true" {
			return true
		}
		return false
	case lua.LTBool:
		b := bool(value.(lua.LBool))
		return b
	default:
		return false
	}
}

func lvalueToStr(value lua.LValue) string {
	switch value.Type() {
	case lua.LTNil:
		return ""
	case lua.LTString:
		return value.(lua.LString).String()
	default:
		return value.String()
	}
}

func rstToLTable(state *lua.LState, p *m_terminal.Result) lua.LValue {
	tab := state.NewTable()
	state.SetTable(tab, lua.LString("msg"), lua.LString(p.Msg()))
	state.SetTable(tab, lua.LString("errInfo"), lua.LString(p.ErrInfo()))
	state.SetTable(tab, lua.LString("code"), lua.LNumber(p.Code()))
	return tab
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

func mapToLTable(src *lua.LState, m map[string]string) *lua.LTable {
	tab := src.NewTable()
	for k, v := range m {
		tab.RawSetString(k, lua.LString(v))
	}
	return tab
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
