package ostype

import (
	"runtime"

	"github.com/Frank-Mayer/yab/internal/lua"
)

func OsType(l *lua.LState) int {
	l.Push(lua.LString(runtime.GOOS))
	return 1
}
