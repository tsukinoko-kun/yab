package util

import (
	"path/filepath"
	"strings"
)

func IsInDir(p string, dir string) bool {
	absDir, err := filepath.Abs(dir)
	if err != nil {
		return false
	}
	absPath, err := filepath.Abs(p)
	if err != nil {
		return false
	}
	return strings.HasPrefix(absPath, absDir)
}
