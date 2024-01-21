package osarch

import (
	"runtime"

	"github.com/Frank-Mayer/gopher-lua"
)

func OsArch(l *lua.LState) int {
	l.Push(lua.LString(runtime.GOARCH))
	return 1
}
