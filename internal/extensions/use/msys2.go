package use

import (
	"fmt"
	"os"
	"os/exec"
	posixpath "path"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/Frank-Mayer/yab/internal/cache"
	"github.com/Frank-Mayer/yab/internal/shell"
	"github.com/Frank-Mayer/yab/internal/util"
	"github.com/charmbracelet/log"
	"github.com/pkg/errors"
)

func useMsys2(version string) error {
	if runtime.GOOS != "windows" {
		log.Warn("Not on windows, msys2 installation skipped")
		return nil
	}

	p, err := cache.InstallPath("msys2", version)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("Error getting install path for msys2 version '%s'", version))
	}

	defer func() {
		util.AddToPath(filepath.Join(p, "msys64"))

		msys2Loc := filepath.Join(p, "msys64", "msys2_shell.cmd")
		log.Debug("Now using msys2 shell", "location", msys2Loc)
		shell.SetShell(func(c string) *exec.Cmd {
			log.Debug("Running msys2 shell", "command", c)
			if wd, err := os.Getwd(); err == nil {
				util.SetEnv("__CD__", wd)
			} else {
				log.Warn("Error getting current working directory", "error", err)
			}
			cmd := exec.Command(msys2Loc,
				"-defterm", "-no-start", "-here",
				"-c", fmt.Sprintf("export PATH=\"$PATH:%s\";%s", WinToPosixPath(util.UsedPath), c))
			return cmd
		})
	}()

	if ok, err := cache.LookupInstall("msys2", version); err == nil {
		if ok {
			log.Debug("Msys2 version already installed", "version", version)
			return nil
		}
	} else {
		return errors.Wrap(err, fmt.Sprintf("Error checking cache for msys2 version '%s'", version))
	}

	log.Info("Installing dependency", "package", "msys2", "version", version)

	filename := "msys2-base-"
	switch runtime.GOARCH {
	case "386":
		return errors.New("msys2 not supported on 32-bit systems")
	case "amd64":
		filename += "x86_64-"
	default:
		return errors.New("Unsupported architecture " + runtime.GOARCH)
	}
	filename += strings.Replace(version, "-", "", -1)
	filename += ".sfx.exe"

	url := "https://github.com/msys2/msys2-installer/releases/download/" + version + "/" + filename

	fp := filepath.Join(p, filename)

	if err := util.Download(url, fp); err != nil {
		return errors.Wrap(err, fmt.Sprintf("Error downloading msys2 version '%s'", version))
	}

	defer func() {
		log.Debug("Removing file", "filepath", fp)
		if err := os.Remove(fp); err != nil {
			log.Error("Error removing file", "filepath", fp, "error", err)
		}
	}()

	log.Debug("Running msys2 installer", "filepath", fp)
	if wd, err := os.Getwd(); err == nil {
		log.Debugf("Got current working directory '%s'", wd)
		if err := os.Chdir(p); err == nil {
			log.Debug("Changed working directory", "path", p)
			err := func() error {
				cmd := exec.Command(fp)
				cmd.Stdout = os.Stdout
				cmd.Stderr = os.Stderr
				cmd.Stdin = os.Stdin
				if err := cmd.Run(); err != nil {
					return errors.Wrap(err, fmt.Sprintf("Error running msys2 installer '%s'", fp))
				}
				log.Debug("Run msys2 installer", "command", cmd)
				return nil
			}()
			if err := os.Chdir(wd); err != nil {
				return errors.Wrap(err, fmt.Sprintf("Error changing directory to '%s'", wd))
			}
			if err != nil {
				return errors.Wrap(err, fmt.Sprintf("Error changing directory to '%s'", p))
			}
		} else {
			return errors.Wrap(err, fmt.Sprintf("Error changing directory to '%s'", p))
		}
	} else {
		return errors.Wrap(err, "Error getting current working directory")
	}

	return nil
}

// WinToPosixPath converts a windows path environment variable to a posix path environment variable.
// It replaces all backslashes with forward slashes and replaces the semicolon with a colon.
// Drive letters are converted to lowercase. (e.g. "C:\Windows" becomes "/c/Windows")
func WinToPosixPath(path string) string {
	parts := strings.Split(path, ";")
	for i := 0; i < len(parts); i++ {
		if len(parts[i]) <= 2 {
			continue
		}
		parts[i] = strings.Replace(parts[i], "\\", "/", -1)
		letter := strings.ToLower(parts[i][:1])
		parts[i] = "/" + letter + parts[i][2:]
		parts[i] = posixpath.Clean(parts[i])
	}
	return strings.Join(parts, ":")
}
