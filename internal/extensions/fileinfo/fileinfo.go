package fileinfo

import (
	"os"
	"github.com/Frank-Mayer/gopher-lua"
)

func FileInfo(l *lua.LState) int {
    path := l.CheckString(1)
    file, err := os.Stat(path)
    if err != nil {
        l.Error(lua.LString(err.Error()), 1)
        return 0
    }
    tbl := l.NewTable()
    l.SetField(tbl, "name", lua.LString(file.Name()))
    l.SetField(tbl, "size", lua.LNumber(file.Size()))
    l.SetField(tbl, "mode", lua.LNumber(file.Mode()))
    l.SetField(tbl, "modtime", lua.LNumber(file.ModTime().Unix()))
    l.SetField(tbl, "isdir", lua.LBool(file.IsDir()))
    l.SetField(tbl, "sys", lua.LNil)
    l.Push(tbl)
    return 1
}
