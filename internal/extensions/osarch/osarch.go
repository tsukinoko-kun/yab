package osarch

import (
	"runtime"

	"github.com/tsukinoko-kun/gopher-lua"
)

func OsArch(l *lua.LState) int {
	l.Push(lua.LString(runtime.GOARCH))
	return 1
}
