package use

import (
	"fmt"

	lua "github.com/Frank-Mayer/yab/internal/lua"
	"github.com/Frank-Mayer/yab/internal/util"
	"github.com/charmbracelet/log"
)

// usable is an array of usable packages
var Usable = []string{"golang", "nodejs", "mingw", "msys2"}

func Use(l *lua.LState) int {
	pack := l.CheckString(1)
	version := l.CheckString(2)

	cancelSpinner := util.Spin()
	defer cancelSpinner()

	if err := use(pack, version); err != nil {
		l.RaiseError("Error ensuring package '%s': %s", pack, err.Error())
	}

	return 0
}

func use(pack string, version string) error {
	if version == "" {
		return fmt.Errorf("No version specified for package '%s'", pack)
	}
	if version[0] == 'v' {
		version = version[1:]
	}
	log.Debug("Ensuring", "package", pack, "version", version)
	switch pack {
	case "golang":
		return useGo(version)
	case "nodejs":
		return useNode(version)
	case "mingw":
		return useMingw(version)
	case "msys2":
		return useMsys2(version)
	}
	return fmt.Errorf("Package '%s' not supported", pack)
}
