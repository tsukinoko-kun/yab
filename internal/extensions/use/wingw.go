package use

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	posixpath "path"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/Frank-Mayer/yab/internal/cache"
	"github.com/Frank-Mayer/yab/internal/util"
	"github.com/charmbracelet/log"
)

// wingwLatest returns the latest wingw version from github api
func wingwLatest() (string, error) {
	resp, err := http.Get("https://api.github.com/repos/brechtsanders/winlibs_mingw/releases/latest")
	if err != nil {
		return "", errors.Join(fmt.Errorf("Failed to request latest wingw version"), err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", errors.Join(fmt.Errorf("Failed to read response body for latest wingw version"), err)
	}
	var result map[string]any
	err = json.Unmarshal([]byte(body), &result)
	if err != nil {
		return "", errors.Join(fmt.Errorf("Failed to unmarshal response body for latest wingw version"), err)
	}
	if tagName, ok := result["tag_name"]; ok {
		return tagName.(string), nil
	}
	if message, ok := result["message"]; ok {
		return "", errors.New(message.(string))
	}
	return "", errors.New("Failed to read latest wingw version from github api response")
}

// wingwTagAssets returns the assets url for the given tag from github api
func wingwTagAssets(tag string) (string, error) {
	resp, err := http.Get("https://api.github.com/repos/brechtsanders/winlibs_mingw/releases/tags/" + tag)
	if err != nil {
		return "", errors.Join(fmt.Errorf("Failed to request wingw version '%s'", tag), err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", errors.Join(fmt.Errorf("Failed to read response body for wingw version '%s'", tag), err)
	}
	var result map[string]any
	err = json.Unmarshal([]byte(body), &result)
	if err != nil {
		return "", errors.Join(fmt.Errorf("Failed to unmarshal response body for wingw version '%s'", tag), err)
	}
	if assetsUrl, ok := result["assets_url"]; ok {
		return assetsUrl.(string), nil
	}
	if message, ok := result["message"]; ok {
		return "", errors.New(message.(string))
	}
	return "", errors.New("Failed to read assets url from github api response")
}

// wingwFindAsset returns the url of the asset with the given tag for the current architecture
func wingwFindAsset(tag string) (string, error) {
	assetsUrl, err := wingwTagAssets(tag)
	if err != nil {
		return "", errors.Join(fmt.Errorf("Failed to get assets url for wingw version '%s'", tag), err)
	}
	resp, err := http.Get(assetsUrl)
	if err != nil {
		return "", errors.Join(fmt.Errorf("Failed to request assets url for wingw version '%s'", tag), err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", errors.Join(fmt.Errorf("Failed to read response body for assets url for wingw version '%s'", tag), err)
	}
	var result []map[string]any
	err = json.Unmarshal([]byte(body), &result)
	if err != nil {
		return "", errors.Join(fmt.Errorf("Failed to unmarshal response body for assets url for wingw version '%s'", tag), err)
	}
	var arch string
	switch runtime.GOARCH {
	case "386":
		arch = "i686"
	case "amd64":
		arch = "x86_64"
	default:
		// return "", fmt.Errorf("Unsupported architecture '%s'", runtime.GOARCH)
		arch = "x86_64"
	}
	for _, asset := range result {
		if name, ok := asset["name"]; ok {
			name := name.(string)
			ext := posixpath.Ext(name)
			if ext != ".zip" {
				continue
			}
			if !strings.Contains(name, "posix") {
				continue
			}
			if !strings.Contains(name, "llvm") {
				continue
			}
			if !strings.Contains(name, arch) {
				continue
			}
			if browserDownloadUrl, ok := asset["browser_download_url"]; ok {
				return browserDownloadUrl.(string), nil
			}
		}
	}
	return "", fmt.Errorf("Failed to find asset for wingw version '%s'", tag)
}

func useWingw(version string) error {
	if runtime.GOOS != "windows" {
		log.Warn("Not on windows, not installing wingw")
		return nil
	}

	if version == "latest" {
		log.Debug("Getting latest wingw version")
		var err error
		version, err = wingwLatest()
		if err != nil {
			return errors.Join(
				fmt.Errorf("Error getting latest wingw version"),
				err,
			)
		}
		log.Warnf(`Latest wingw version is '%s' you should use this version with 'yab.use("wingw", "%s")'`, version, version)
	}

	p, err := cache.InstallPath("mingw", version)
	if err != nil {
		return errors.Join(
			fmt.Errorf("Error getting install path for mingw version '%s'", version),
			err,
		)
	}

	defer func() {
		var path string
		var mingw32path string
		switch runtime.GOARCH {
		case "386":
			path = filepath.Join(p, "mingw32")
			mingw32path = filepath.Join(path, "i686-w64-mingw32")
		case "amd64":
			path = filepath.Join(p, "mingw64")
			mingw32path = filepath.Join(path, "x86_64-w64-mingw32")
		default:
			log.Error("Unsupported architecture", "architecture", runtime.GOARCH)
		}
		util.AddToPath(filepath.Join(mingw32path, "bin"))
		util.AddToPath(filepath.Join(path, "bin"))
		util.SetEnv("GCC_EXEC_PREFIX", filepath.Join(mingw32path, "lib", "gcc"))
		util.SetEnv("CPATH", filepath.Join(mingw32path, "include"))
		util.PushEnv("CPATH", filepath.Join(path, "include"))
		util.SetEnv("C_INCLUDE_PATH", filepath.Join(mingw32path, "include"))
		util.PushEnv("C_INCLUDE_PATH", filepath.Join(path, "include"))
		util.SetEnv("CPLUS_INCLUDE_PATH", filepath.Join(mingw32path, "include"))
		util.PushEnv("CPLUS_INCLUDE_PATH", filepath.Join(path, "include"))
		util.SetEnv("LIBRARY_PATH", filepath.Join(path, "lib")+":"+filepath.Join(mingw32path, "lib"))
	}()

	if ok, err := cache.LookupInstall("mingw", version); err == nil {
		if ok {
			log.Debug("Mingw version already installed", "version", version)
			return nil
		}
	} else {
		return errors.Join(
			fmt.Errorf("Error checking cache for mingw version '%s'", version),
			err,
		)
	}

	log.Info("Installing dependency", "package", "mingw", "version", version)

	assetUrl, err := wingwFindAsset(version)
	if err != nil {
		return errors.Join(
			fmt.Errorf("Error getting assets url for wingw version '%s'", version),
			err,
		)
	}

	filename := posixpath.Base(assetUrl)
	filepath := filepath.Join(p, filename)

	resp, err := http.Get(assetUrl)
	if err != nil {
		return errors.Join(
			fmt.Errorf("Error downloading wingw version '%s'", version),
			err,
		)
	}
	defer resp.Body.Close()

	f, err := os.Create(filepath)
	if err != nil {
		return errors.Join(
			fmt.Errorf("Error creating wingw version '%s'", version),
			err,
		)
	}
	defer f.Close()

	_, err = io.Copy(f, resp.Body)
	if err != nil {
		return errors.Join(
			fmt.Errorf("Error writing wingw version '%s'", version),
			err,
		)
	}

	if err := util.Unzip(filepath); err != nil {
		return errors.Join(
			fmt.Errorf("Error unzipping file '%s'", filepath),
			err,
		)
	}

	defer func() {
		if err := os.Remove(filepath); err != nil {
			log.Warn("Error removing file", "filepath", filepath)
		}
	}()

	return nil
}
