package util

import (
	"os"
	"runtime"
)

var origEnv map[string]string

func AddToPath(pathToAdd string) {
	PushEnv("PATH", pathToAdd)
}

func PushEnv(key string, pathToAdd string) {
	// check if pathToAdd exists
	_, err := os.Stat(pathToAdd)
	if err != nil {
		return
	}

	if path, ok := os.LookupEnv(key); ok {
        if runtime.GOOS == "windows" {
            SetEnv(key, pathToAdd+";"+path)
        } else {
            SetEnv(key, pathToAdd+":"+path)
        }
	} else {
        SetEnv(key, pathToAdd)
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
