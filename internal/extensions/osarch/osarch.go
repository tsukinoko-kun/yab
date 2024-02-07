package osarch

import (
	"runtime"

	"github.com/Frank-Mayer/yab/internal/lua"
)

func OsArch(l *lua.LState) int {
	l.Push(lua.LString(runtime.GOARCH))
	return 1
}
