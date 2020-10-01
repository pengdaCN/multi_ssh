package playbook

import lua "github.com/yuin/gopher-lua"

func setGlobalVal(vm *lua.LState, key, val string) {
	vm.SetGlobal(key, lua.LString(val))
}

func SetGlobalVal(key, val string) {
	setGlobalVal(VM, key, val)
}
