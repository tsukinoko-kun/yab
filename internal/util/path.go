package util

import (
	"os"
	"runtime"
)

var origPath string

func AddToPath(pathToAdd string) {
	// check if pathToAdd exists
	_, err := os.Stat(pathToAdd)
	if err != nil {
		return
	}

	path := os.Getenv("PATH")
	if origPath == "" {
		origPath = path
	}
	if runtime.GOOS == "windows" {
		os.Setenv("PATH", pathToAdd+";"+path)
	} else {
		os.Setenv("PATH", pathToAdd+":"+path)
	}
}

func RestorePath() {
	if origPath != "" {
		os.Setenv("PATH", origPath)
	}
}
