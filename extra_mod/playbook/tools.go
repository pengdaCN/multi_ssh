package playbook

import (
	"context"
	lua "github.com/yuin/gopher-lua"
	"multi_ssh/m_terminal"
	"strconv"
	"strings"
	"time"
)

func SetReadOnly(l *lua.LState, table *lua.LTable) *lua.LUserData {
	ud := l.NewUserData()
	mt := l.NewTable()
	// 设置表中域的指向为 table
	l.SetField(mt, "__index", table)
	// 限制对表的更新操作
	l.SetField(mt, "__newindex", l.NewFunction(func(state *lua.LState) int {
		state.RaiseError("not allow to modify table")
		return 0
	}))
	ud.Metatable = mt
	return ud
}

func useTimeoutFromLvalue(value lua.LValue) context.Context {
	switch value.Type() {
	case lua.LTNil:
		return context.Background()
	case lua.LTNumber:
		t := value.(lua.LNumber).String()
		if v, err := strconv.ParseInt(t, 0, 32); err != nil {
			return context.Background()
		} else {
			ctx, _ := context.WithTimeout(context.Background(), time.Second*time.Duration(v))
			return ctx
		}
	default:
		return context.Background()
	}
}

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
	state.SetTable(tab, lua.LString("totalMsg"), lua.LString(p.TotalMsg()))
	state.SetTable(tab, lua.LString("stdout"), lua.LString(p.Stdout()))
	state.SetTable(tab, lua.LString("stderr"), lua.LString(p.Stderr()))
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

func strSliceToTable(table *lua.LTable, arr []string) {
	for i, v := range arr {
		table.Insert(i+1, lua.LString(v))
	}
}

//func buildHandleByFileWithLTable(src *lua.LTable) m_terminal.HandleByFile {
//	m := lTableToMapStrLValue(src)
//	var (
//		uid  int
//		gid  int
//		mode string
//	)
//	if u, ok := m["uid"]; ok {
//		uid, _ = strconv.Atoi(u.String())
//	} else {
//		uid = -1
//	}
//
//	if g, ok := m["gid"]; ok {
//		gid, _ = strconv.Atoi(g.String())
//	} else {
//		gid = -1
//	}
//	if m, ok := m["mode"]; ok {
//		mode = m.String()
//	}
//
//	return func(file *sftp.File) error {
//		if mode != "" {
//			m, err := tools.String2FileMode(mode)
//			if err != nil {
//				return err
//			}
//			if err := file.Chmod(m); err != nil {
//				return err
//			}
//		}
//		if uid != -1 && gid != -1 {
//			if err := file.Chown(uid, gid); err != nil {
//				return err
//			}
//		}
//		return nil
//	}
//}
