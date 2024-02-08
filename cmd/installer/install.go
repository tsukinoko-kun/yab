package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/Frank-Mayer/yab/internal/mainutil"
	"github.com/Frank-Mayer/yab/internal/util"
	"github.com/charmbracelet/log"
	"github.com/pkg/errors"
)

func Install(w fyne.Window) {
	w.SetIcon(theme.DownloadIcon())
	log.Info("Showing install")

	mainutil.Prepare()

	// set up progress bar
	progress := widget.NewProgressBar()
	progress.Min = 0
	progress.Max = 4

	info := widget.NewLabel("")

	// set up window content
	w.SetContent(container.NewVBox(
		widget.NewLabel("Installing"),
		progress,
		info,
	))

	go func() {
		// find target location
		targetDir := installPath()
		if err := os.MkdirAll(targetDir, 0777); err != nil {
			errorLog(info, "Could not create target directory", "path", targetDir, "error", err)
			return
		}
		targetLocation := filepath.Join(targetDir, "yab")
		if runtime.GOOS == "windows" {
			targetLocation += ".exe"
		}
		func() {
			infoLog(info, "Installing to "+targetLocation)
			if exists(targetLocation) {
				if err := os.Remove(targetLocation); err != nil {
					errorLog(info, "Could not remove old file", "path", targetLocation, "error", err)
					return
				}
			}
			binTarget, err := os.OpenFile(targetLocation, os.O_CREATE|os.O_WRONLY, 0777)
			if err != nil {
				errorLog(info, "Could not create target file", "path", targetLocation, "error", err)
				return
			}
			infoLog(info, "Created target file", "path", targetLocation)
			defer binTarget.Close()
			progress.SetValue(progress.Value + 1)

			// start download
			url := "https://frank-mayer.github.io/yab/yab-" + runtime.GOOS + "-" + runtime.GOARCH
			if runtime.GOOS == "windows" {
				url += ".exe"
			}
			req, err := http.NewRequest("GET", url, nil)
			if err != nil {
				errorLog(info, "Could not create request", "url", url, "error", err)
				return
			}
			infoLog(info, "Downloading binary", "url", url)
			req.Header.Set("User-Agent", "yab-installer")
			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				errorLog(info, "Could not download binary", "url", url, "error", err)
				return
			}
			defer resp.Body.Close()
			if _, err := io.Copy(binTarget, resp.Body); err != nil {
				errorLog(info, "Could not write binary", "error", err)
				return
			}
			infoLog(info, "Downloaded binary", "url", url)
			progress.SetValue(progress.Value + 1)

			// set executable flag
			if runtime.GOOS != "windows" {
				if err := os.Chmod(targetLocation, 0777); err != nil {
					errorLog(info, "Could not set executable flag", "error", err)
					return
				}
				infoLog(info, "Set executable flag", "path", targetLocation)
			}
		}()

		// write lua definitions
		libPath := filepath.Join(targetDir, "lib")
		if err := os.RemoveAll(libPath); err != nil {
			errorLog(info, "Could not remove lib directory", "error", err)
		}
		if err := exe(targetLocation, "--def"); err != nil {
			errorLog(info, "Could not write lua definitions", "error", err)
			return
		}
		infoLog(info, "Wrote lua definitions")
		progress.SetValue(progress.Value + 1)

		// add to path
		path := os.Getenv("PATH")
		if !strings.Contains(path, targetDir) {
			if runtime.GOOS == "windows" {
				if err := addWindowsPath(targetDir); err != nil {
					errorLog(info, "Could not add path to shell", "error", err)
					return
				}
			} else {
				if err := addUnixPath(targetDir); err != nil {
					errorLog(info, "Could not add path to shell", "error", err)
					return
				}
			}
			infoLog(info, "Added to path")
		} else {
			infoLog(info, "Already in path")
		}
		progress.SetValue(progress.Value + 1)

		// show success
		infoLog(info, "Yab was successfully installed!")
		infoLog(info, "You can now use it from the command line. For help, run 'yab --help'")
		infoLog(info, "You have to restart your terminal to make it reload the PATH variable.")

		w.Content().(*fyne.Container).Add(widget.NewButton("Close", func() {
			w.Close()
		}))
	}()
}

func addUnixPath(p string) error {
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	xdgConfigHome := os.Getenv("XDG_CONFIG_HOME")
	shellEnv := os.Getenv("SHELL")
	shell := filepath.Base(shellEnv)
	switch shell {
	case "bash":
		bashrc := filepath.Join(home, ".bashrc")
		if exists(bashrc) {
			log.Info("Adding to bashrc")
			return appent(bashrc, `export PATH="`+p+`:$PATH"`)
		}
		bash_profile := filepath.Join(home, ".bash_profile")
		if exists(bash_profile) {
			log.Info("Adding to bash_profile")
			return appent(bash_profile, `export PATH="`+p+`:$PATH"`)
		}
	case "zsh":
		zshrc := filepath.Join(home, ".zshrc")
		if exists(zshrc) {
			log.Info("Adding to zshrc")
			return appent(zshrc, `export PATH="`+p+`:$PATH"`)
		}
	case "fish":
		fishrc := filepath.Join(xdgConfigHome, "fish", "config.fish")
		if exists(fishrc) {
			log.Info("Adding to XDG_CONFIG_HOME config.fish")
			return appent(fishrc, `export PATH="`+p+`:$PATH"`)
		}
		fishrc = filepath.Join(home, ".config", "fish", "config.fish")
		if exists(fishrc) {
			log.Info("Adding to home config.fish")
			return appent(fishrc, `export PATH="`+p+`:$PATH"`)
		}
	case "csh":
		cshrc := filepath.Join(home, ".cshrc")
		if exists(cshrc) {
			return appent(cshrc, `setenv PATH "`+p+`:$PATH"`)
		}
	case "tcsh":
		tcshrc := filepath.Join(home, ".tcshrc")
		if exists(tcshrc) {
			return appent(tcshrc, `setenv PATH "`+p+`:$PATH"`)
		}
	}
	profile := filepath.Join(home, ".profile")
	if exists(profile) {
		return appent(profile, `export PATH="`+p+`:$PATH"`)
	}
	return fmt.Errorf("could not add path to shell")
}

func addWindowsPath(directory string) error {
	currentPath := os.Getenv("PATH")
	cmd := exec.Command("reg", "add", "HKCU\\Environment", "/v", "Path", "/t", "REG_SZ", "/d", directory+";"+currentPath, "/f")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		return errors.Wrap(err, "could not add path to shell")
	}
	return nil
}

func appent(path string, content string) error {
	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return errors.Wrap(
			err,
			fmt.Sprintf("could not open %s", path),
		)
	}
	defer f.Close()

	if _, err := f.WriteString("\n" + content + "\n"); err != nil {
		return errors.Wrap(err, fmt.Sprintf("could not write to %s", path))
	}
	return nil
}

func installPath() string {
	// is there a yab binary in path?
	if bin, err := exec.LookPath("yab"); err == nil {
		// yes, use that
		p := filepath.Dir(bin)
		log.Info("Found existing yab binary", "path", bin, "dir", p)
		return p
	}
	globalConfigPath, err := util.GetGlobalConfigPath()
	if err != nil {
		log.Fatal(err)
	}
	return filepath.Join(globalConfigPath, "bin")
}

func exists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func exe(bin string, args ...string) error {
	cmd := exec.Command(bin, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	return cmd.Run()
}

func infoLog(el *widget.Label, msg interface{}, keyvals ...interface{}) {
	log.Info(msg, keyvals...)

	sb := strings.Builder{}
	sb.WriteString(fmt.Sprintf("%v ", msg))

	if len(keyvals)%2 != 0 {
		keyvals = append(keyvals, "")
	}

	for i := 0; i < len(keyvals); i += 2 {
		sb.WriteString(fmt.Sprintf("%s=%v ", keyvals[i], keyvals[i+1]))
	}

	el.SetText(el.Text + strings.TrimSpace(sb.String()) + "\n")
}

func errorLog(el *widget.Label, msg interface{}, keyvals ...interface{}) {
	log.Error(msg, keyvals...)

	content := el.Text
	content += fmt.Sprintf("ERROR %v ", msg)

	if len(keyvals)%2 != 0 {
		keyvals = append(keyvals, "")
	}

	for i := 0; i < len(keyvals); i += 2 {
		content += fmt.Sprintf("%s=%v\n", keyvals[i], keyvals[i+1])
	}

	content += "\n"
	el.SetText(content)
}
