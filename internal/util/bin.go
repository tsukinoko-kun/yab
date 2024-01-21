package util

import (
	"errors"
	"os"
	"path/filepath"
)

var (
	ConfigPath       string
	globalConfigPath string
)

func BinName() string {
	binName := os.Args[0]
	if len(binName) > 24 {
		return "yab"
	}
	return binName
}

func GetPackagePath() string {
	return filepath.Join(ConfigPath, "?.lua") + ";" +
		filepath.Join(ConfigPath, "?", "init.lua") + ";" +
		filepath.Join(".", "?.lua") + ";" +
		filepath.Join("?", "init.lua")
}

func GetGlobalConfigPath() (string, error) {
	if globalConfigPath != "" {
		return globalConfigPath, nil
	}
	var err error
	if globalConfigPath, err = getGlobalConfigPath(); err != nil {
		return "", err
	}
	err = os.MkdirAll(globalConfigPath, 0777)
	return globalConfigPath, err
}

func getGlobalConfigPath() (string, error) {
	if xdgConfigHome, exists := os.LookupEnv("XDG_CONFIG_HOME"); exists {
		pathname := filepath.Join(xdgConfigHome, "yab")
		return pathname, nil
	}

	if home, exists := os.LookupEnv("APPDATA"); exists {
		pathname := filepath.Join(home, "yab")
		return pathname, nil
	}

	if home, exists := os.LookupEnv("HOME"); exists {
		pathname := filepath.Join(home, ".config", "yab")
		return pathname, nil
	}

	return "", errors.New("could not find or create global config path")
}
