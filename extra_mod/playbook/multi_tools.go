package playbook

import (
	lua "github.com/yuin/gopher-lua"
	"regexp"
	"strings"
)

func initStr(state *lua.LState, table *lua.LTable) {
	table.RawSetString("split", state.NewFunction(split))
	table.RawSetString("hasPrefix", state.NewFunction(hasPrefix))
	table.RawSetString("hasSuffix", state.NewFunction(hasSuffix))
}

func split(state *lua.LState) int {
	arr := state.NewTable()
	defer func() {
		state.Push(arr)
	}()
	var (
		str string
		sep string
	)
	str = state.ToString(1)
	{
		val := state.Get(2)
		switch val.Type() {
		case lua.LTNil:
			str = " "
		case lua.LTString:
			str = val.String()
		}
	}
	_arr := strings.Split(str, sep)
	strSliceToTable(arr, _arr)
	return 1
}

func hasPrefix(state *lua.LState) int {
	var (
		str    string
		prefix string
	)
	str = state.ToString(1)
	prefix = state.ToString(2)
	b := strings.HasPrefix(str, prefix)
	state.Push(lua.LBool(b))
	return 1
}

func hasSuffix(state *lua.LState) int {
	var (
		str    string
		suffix string
	)
	str = state.ToString(1)
	suffix = state.ToString(2)
	b := strings.HasSuffix(str, suffix)
	state.Push(lua.LBool(b))
	return 1
}

func initRe(state *lua.LState, table *lua.LTable) {
	table.RawSetString("match", state.NewFunction(reMatch))
}

func reMatch(state *lua.LState) int {
	var (
		str string
		re  string
	)
	str = state.ToString(1)
	{
		val := state.Get(2)
		switch val.Type() {
		case lua.LTNil:
			state.Push(lua.LFalse)
			return 1
		case lua.LTString:
			re = val.String()
		}
	}
	b, err := regexp.MatchString(re, str)
	if err != nil {
		state.Error(lua.LString(err.Error()), 1)
		return 0
	}
	state.Push(lua.LBool(b))
	return 1
}
