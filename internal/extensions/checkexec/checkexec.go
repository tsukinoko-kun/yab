package checkexec

import (
	"os/exec"

	"github.com/Frank-Mayer/gopher-lua"
)

func CheckExec(l *lua.LState) int {
	pkg := l.CheckString(1)
	_, err := exec.LookPath(pkg)
	if err != nil {
		l.Push(lua.LFalse)
	} else {
		l.Push(lua.LTrue)
	}
	return 1
}
