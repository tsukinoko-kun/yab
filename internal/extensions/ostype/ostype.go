package ostype

import (
	"runtime"

	"github.com/tsukinoko-kun/gopher-lua"
)

func OsType(l *lua.LState) int {
	l.Push(lua.LString(runtime.GOOS))
	return 1
}
