package watch

import (
	"github.com/tsukinoko-kun/gopher-lua"
	"github.com/charmbracelet/log"
	"github.com/fsnotify/fsnotify"
)

func Watch(l *lua.LState) int {
	paths := l.CheckTable(1)
	callback := l.CheckFunction(2)

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		l.Error(lua.LString("Error creating watcher. "+err.Error()), 0)
		return 0
	}

	go func() {
		defer watcher.Close()
		for {
			select {
			case event := <-watcher.Events:
				if event.Op != 0 {
					var op string
					switch event.Op {
					case fsnotify.Create:
						op = "create"
					case fsnotify.Write:
						op = "write"
					case fsnotify.Remove:
						op = "remove"
					case fsnotify.Rename:
						op = "rename"
					case fsnotify.Chmod:
						op = "chmod"
					default:
						op = "unknown"
					}
					err := l.CallByParam(lua.P{
						Fn:      callback,
						NRet:    0,
						Protect: true,
					}, lua.LString(event.Name), lua.LString(op))
					if err != nil {
						l.Error(lua.LString("Error calling callback. "+err.Error()), 0)
						return
					}
				}
			case err := <-watcher.Errors:
				if err != nil {
					l.Error(lua.LString("Error watching path. "+err.Error()), 0)
					return
				}
			}
		}
	}()

	for i := 1; i <= paths.Len(); i++ {
		path := paths.RawGetInt(i).String()
		log.Debug("Adding path to watcher", "path", path)
		err = watcher.Add(path)
		if err != nil {
			l.Error(lua.LString("Error adding path to watcher. "+err.Error()), 0)
			return 0
		}
	}

	return 0
}
