package ostype

import (
	"runtime"

	"github.com/Frank-Mayer/gopher-lua"
)

func OsType(l *lua.LState) int {
	l.Push(lua.LString(runtime.GOOS))
	return 1
}
