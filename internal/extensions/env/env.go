package env

import (
	"os"

	"github.com/Frank-Mayer/gopher-lua"
)

func SetEnv(l *lua.LState) int {
	key := l.CheckString(1)
	val := l.CheckString(2)
	l.Env.RawSet(lua.LString(key), lua.LString(val))
	os.Setenv(key, val)
	return 0
}
