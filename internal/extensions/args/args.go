package args

import "github.com/tsukinoko-kun/gopher-lua"

var luaArgs []string

func SetArgs(args []string) {
	luaArgs = args
}

func Args(l *lua.LState) int {
	table := l.NewTable()
	for i := 0; i < len(luaArgs); i++ {
		l.SetTable(table, lua.LNumber(i+1), lua.LString(luaArgs[i]))
	}
	l.Push(table)
	return 1
}
