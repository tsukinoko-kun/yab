package cache

import (
	"os"
	"path/filepath"
	"runtime"

	"github.com/tsukinoko-kun/yab/internal/util"
)

func InstallPath(pack string, version string) (string, error) {
	p, err := util.GetGlobalConfigPath()
	if err != nil {
		return "", err
	} else {
		p = filepath.Join(p, "cache_"+runtime.GOARCH, "install", pack, version)
		// make sure the path exists
		if err := os.MkdirAll(p, 0777); err != nil {
			return "", err
		}
		return p, nil
	}
}

func LookupInstall(pack string, version string) (bool, error) {
	p, err := InstallPath(pack, version)
	if err != nil {
		return false, err
	}
	_, err = os.Stat(p)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}
	// check if there is a file in the directory
	// if there is, we assume it is installed

	// loop through the directory and check for files
	entries, err := os.ReadDir(p)
	if err != nil {
		return false, err
	}
	for _, entry := range entries {
		// if a directory
		if entry.IsDir() {
			return true, nil
		}
		// if not a archive
		ext := filepath.Ext(entry.Name())
		if ext == ".zip" || ext == ".gz" || ext == ".tar" {
			return true, nil
		}
	}

	return false, nil
}
