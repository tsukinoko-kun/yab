package download

import (
	"strings"

	lua "github.com/Frank-Mayer/gopher-lua"
	"github.com/Frank-Mayer/yab/internal/util"
)

func Download(l *lua.LState) int {
	src := l.CheckString(1)
	dest := l.OptString(2, "")

	stopSpinner := util.Spin()
	defer stopSpinner()

	if dest == "" {
		i := strings.LastIndexByte(src, '/') + 1
		dest = src[i:]
	}

	if err := util.Download(src, dest); err != nil {
		l.RaiseError("Error downloading file. " + err.Error())
		return 0
	}

	l.Push(lua.LString(dest))
	return 1
}
