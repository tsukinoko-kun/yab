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

func useNode(version string) error {
	p, err := cache.InstallPath("node", version)
	if err != nil {
		return errors.Join(
			fmt.Errorf("Error getting install path for node version '%s'", version),
			err,
		)
	}

	defer func() {
		dir := "node-v" + version + "-"
		switch runtime.GOOS {
		case "darwin":
			dir += "darwin-"
		case "linux":
			dir += "linux-"
		case "windows":
			dir += "win-"
		default:
			return
		}
		switch runtime.GOARCH {
		case "386":
			dir += "x86"
		case "amd64":
			dir += "x64"
		case "arm64":
			dir += "arm64"
		case "arm":
			dir += "armv7l"
		default:
			return
		}
		util.AddToPath(filepath.Join(p, dir))
		util.AddToPath(filepath.Join(p, dir, "bin"))
		util.AddToPath(filepath.Join(p, dir, "node_modules", "corepack", "shims"))
		if runtime.GOOS == "windows" {
			util.AddToPath(filepath.Join(p, dir, "node_modules", "corepack", "shims", "nodewin"))
		}
		if nodeModulesPath, err := util.Find("node_modules", p); err == nil {
			util.AddToPath(filepath.Join(nodeModulesPath, ".bin"))
		}
		nodeModulesPath := filepath.Join(p, "node_modules")
		util.SetEnv("NODE_PATH", nodeModulesPath)
		util.AddToPath(filepath.Join(nodeModulesPath, ".bin"))
	}()

	if ok, err := cache.LookupInstall("node", version); err == nil {
		if ok {
			log.Debug("Node version already installed", "version", version)
			return nil
		}
	} else {
		return errors.Join(
			fmt.Errorf("Error checking cache for node version '%s'", version),
			err,
		)
	}

	log.Info("Installing dependency", "package", "node", "version", version)

	filename := "node-v" + version + "-"
	switch runtime.GOOS {
	case "darwin":
		filename += "darwin-"
	case "linux":
		filename += "linux-"
	case "windows":
		filename += "win-"
	default:
		return fmt.Errorf("Unsupported OS '%s'", runtime.GOOS)
	}

	switch runtime.GOARCH {
	case "386":
		filename += "x86"
	case "amd64":
		filename += "x64"
	case "arm64":
		filename += "arm64"
	case "arm":
		filename += "armv7l"
	default:
		return fmt.Errorf("Unsupported architecture '%s'", runtime.GOARCH)
	}

	switch runtime.GOOS {
	case "windows":
		filename += ".zip"
	case "darwin":
		filename += ".tar.gz"
	case "linux":
		filename += ".tar.xz"
	}

	url := "https://nodejs.org/dist/v" + version + "/" + filename

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
