package playbook

import (
	"context"
	lua "github.com/yuin/gopher-lua"
	"multi_ssh/m_terminal"
	"strconv"
	"strings"
	"time"
)

func LuaValueToGoVal(val lua.LValue) interface{} {
	switch val.Type() {
	case lua.LTString:
		return val.String()
	case lua.LTNumber:
		v := val.(lua.LNumber)
		return float64(v)
	case lua.LTBool:
		v := val.(lua.LBool)
		return bool(v)
	case lua.LTTable:
		v := val.(*lua.LTable)
		m := make(map[interface{}]interface{})
		intoLuaMap(m, v)
		return m
	default:
		return nil
	}
}

func intoLuaMap(m map[interface{}]interface{}, val *lua.LTable) {
	val.ForEach(func(value lua.LValue, value2 lua.LValue) {
		key := LuaValueToGoVal(value)
		_val := LuaValueToGoVal(value2)
		m[key] = _val
	})
}

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
