package use

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/Frank-Mayer/yab/internal/cache"
	"github.com/Frank-Mayer/yab/internal/util"
	"github.com/charmbracelet/log"
)

func useGo(version string) error {
	p, err := cache.InstallPath("go", version)
	if err != nil {
		return errors.Join(
			fmt.Errorf("Error getting install path for go version '%s'", version),
			err,
		)
	}

	defer func() {
		util.AddToPath(filepath.Join(p, "go", "bin"))
		util.SetEnv("GOROOT", filepath.Join(p, "go"))
		if projectCachePath, err := cache.ProjCachePath(); err == nil {
			util.SetEnv("GOPATH", filepath.Join(projectCachePath, "goworkspace"))
		} else {
			log.Warn("Error setting GOPATH", "error", err)
		}
	}()

	if ok, err := cache.LookupInstall("go", version); err == nil {
		if ok {
			log.Debug("Go version already installed", "version", version)
			return nil
		}
	} else {
		return errors.Join(
			fmt.Errorf("Error checking cache for go version '%s'", version),
			err,
		)
	}

	log.Info("Installing dependency", "package", "go", "version", version)

	var filename string
	if runtime.GOOS == "windows" {
		filename = "go" + version + "." + runtime.GOOS + "-" + runtime.GOARCH + ".zip"
	} else {
		filename = "go" + version + "." + runtime.GOOS + "-" + runtime.GOARCH + ".tar.gz"
	}
	url := "https://dl.google.com/go/" + filename

	filepath := filepath.Join(p, filename)

	if err := util.Download(url, filepath); err != nil {
		return err
	}
	defer func() {
		log.Debug("Removing file", "filepath", filepath)
		if err := os.Remove(filepath); err != nil {
			log.Warn("Error removing file", "filepath", filepath, "error", err)
		}
	}()

	if err := util.Unzip(filepath); err != nil {
		return errors.Join(
			fmt.Errorf("Error unzipping file '%s'", filepath),
			err,
		)
	}

	return nil
}
