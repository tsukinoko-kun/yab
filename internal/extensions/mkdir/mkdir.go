package mkdir

import (
	"os"

	lua "github.com/Frank-Mayer/yab/internal/lua"
)

func Mkdir(l *lua.LState) int {
	dir := l.CheckString(1)
	err := os.Mkdir(dir, 0755)
	if err != nil {
		l.RaiseError("Error creating directory: %s", err)
	}
	return 0
}
