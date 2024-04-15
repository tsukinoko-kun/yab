package cd

import (
	"os"
	"path/filepath"
	"sync"

	"github.com/tsukinoko-kun/gopher-lua"
)

// mutex is used to prevent concurrent calls to change the current working directory.
var mutex = &sync.Mutex{}

// Cd changes the current working directory to the given path for one function call.
func Cd(l *lua.LState) int {
	mutex.Lock()
	defer mutex.Unlock()

	path := l.CheckString(1)
	fn := l.CheckFunction(2)

	absPath, err := filepath.Abs(path)
	if err != nil {
		l.RaiseError("Error getting absolute path: %s", err.Error())
		return 0
	}

	cwd, err := os.Getwd()
	if err != nil {
		l.RaiseError("Error getting current working directory: %s", err.Error())
		return 0
	}

	err = os.Chdir(absPath)
	if err != nil {
		l.RaiseError("Error changing directory: %s", err.Error())
		return 0
	}

	defer func() {
		err = os.Chdir(cwd)
		if err != nil {
			l.RaiseError("Error changing directory back: %s", err.Error())
		}
	}()

	err = l.CallByParam(lua.P{
		Fn:      fn,
		NRet:    0,
		Protect: true,
	})

	if err != nil {
		l.RaiseError("Error calling function: %s", err.Error())
	}

	return 0
}
