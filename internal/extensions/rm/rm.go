package rm

import (
	"os"

	"github.com/tsukinoko-kun/gopher-lua"
)

func Rm(l *lua.LState) int {
	p := l.CheckString(1)
	rec := l.OptBool(2, false)
	if rec {
		err := os.RemoveAll(p)
		if err != nil {
			l.RaiseError("Error removing directory: %s", err)
		}
	} else {
		err := os.Remove(p)
		if err != nil {
			l.RaiseError("Error removing file: %s", err)
		}
	}
	return 0
}
