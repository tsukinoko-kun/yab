package util

import (
	"os"
	"runtime"
)

var origEnv map[string]string

func AddToPath(pathToAdd string) {
	// check if pathToAdd exists
	_, err := os.Stat(pathToAdd)
	if err != nil {
		return
	}

	path := os.Getenv("PATH")

	if runtime.GOOS == "windows" {
		SetEnv("PATH", pathToAdd+";"+path)
	} else {
		SetEnv("PATH", pathToAdd+":"+path)
	}
}

func RestoreEnv() {
	for key, value := range origEnv {
		if value == "" {
			os.Unsetenv(key)
		} else {
			os.Setenv(key, value)
		}
	}
}

func SetEnv(key string, value string) {
	if origEnv == nil {
		origEnv = make(map[string]string)
	}
	if _, ok := origEnv[key]; !ok {
		origEnv[key] = os.Getenv(key)
	}
	os.Setenv(key, value)
}
