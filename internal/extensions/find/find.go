package find

import (
	"github.com/Frank-Mayer/gopher-lua"

	"os"
	"path/filepath"
)

func Find(l *lua.LState) int {
	// can be called (<pattern>) or (<root>, <pattern>)
	var root string
	var pattern string

	switch l.GetTop() {
	case 1:
		root = "."
		pattern = l.CheckString(1)
	case 2:
		root = l.CheckString(1)
		pattern = l.CheckString(2)
	default:
		l.ArgError(1, "expected 1 or 2 arguments")
	}

	matches, err := WalkMatch(root, pattern)
	if err != nil {
		l.RaiseError(err.Error())
		return 0
	}

	table := l.NewTable()
	for i := 0; i < len(matches); i++ {
		l.SetTable(table, lua.LNumber(i+1), lua.LString(matches[i]))
	}
	l.Push(table)
	return 1
}

func WalkMatch(root string, pattern string) ([]string, error) {
	var matches []string
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		if matched, err := filepath.Match(pattern, filepath.Base(path)); err != nil {
			return err
		} else if matched {
			matches = append(matches, path)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return matches, nil
}
