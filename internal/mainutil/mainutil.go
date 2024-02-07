package mainutil

import (
	"os"
	"path/filepath"
	"strings"
	"time"

	lua "github.com/Frank-Mayer/gopher-lua"
	"github.com/Frank-Mayer/yab/internal/extensions"
	"github.com/Frank-Mayer/yab/internal/util"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/log"
)

func GetInitFile(configPath string, file string) (string, error) {
	initFile := filepath.Join(configPath, file+".lua")
	if _, err := os.Stat(initFile); err != nil {
		return "", err
	}
	return initFile, nil
}

func RunLuaFile(initFile string) error {
	// setup lua
	l := lua.NewState()
	defer l.Close()
	extensions.RegisterExtensions(l)

	packagePath := util.GetPackagePath()
	setupCode := "package.path = '" + strings.ReplaceAll(packagePath, "\\", "\\\\") + ";'"
	err := l.DoString(setupCode)
	if err != nil {
		log.Error("Error setting up lua", "error", err, "code", setupCode)
		return err
	}

	// run lua file
	err = l.DoFile(initFile)
	if err != nil {
		return err
	}

	return nil
}

func GetConfigPath() (string, error) {
	rootPath, err := os.Getwd()
	if err != nil {
		return "", err
	}

	for {
		pathname := filepath.Join(rootPath, ".yab")
		if _, err := os.Stat(pathname); !os.IsNotExist(err) {
			if err := os.Chdir(filepath.Dir(pathname)); err != nil {
				return pathname, err
			}
			return pathname, nil
		}
		// parent directory
		parent := filepath.Dir(rootPath)
		if parent == rootPath {
			break
		}
		rootPath = parent
	}

	return util.GetGlobalConfigPath()
}

func InitDefinitons() {
	configPath, err := util.GetGlobalConfigPath()
	if err != nil {
		log.Fatal(err)
	}

	libPath := filepath.Join(configPath, "lib")
	filename := filepath.Join(libPath, "yab.lua")

	// create directory
	err = os.MkdirAll(libPath, 0775)
	if err != nil {
		log.Fatal(err)
	}

	// does file already exist?
	_, err = os.Stat(filename)
	if err == nil {
		log.Info("Lua API definitions already exist, overwriting", "path", filename)
		// rename old file (for windows)
		old := filename + ".old." + time.Now().Format("20060102150405")
		err = os.Rename(filename, old)
		if err != nil {
			log.Fatal("Error renaming old file", "error", err, "path", old)
		}
		// remove old file
		err = os.Remove(old)
		if err != nil {
			log.Warn("Error removing old file", "error", err, "path", old)
		}
	}

	// create file
	f, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	// write file
	_, err = f.WriteString(extensions.Definitions())
	if err != nil {
		log.Fatal(err)
	}
	log.Info("Lua API definitions created", "location", libPath)
}

func PrintEnv() {
	sb := strings.Builder{}
	sb.WriteString("# Yab environment\n\n")

	configPath, err := util.GetGlobalConfigPath()
	if err != nil {
		log.Fatal(err)
	}

	sb.WriteString("Global config path `")
	sb.WriteString(configPath)
	sb.WriteString("`\n\n")

	sb.WriteString("Local config path `")
	sb.WriteString(util.ConfigPath)
	sb.WriteString("`\n\n")

	sb.WriteString("Lua library path `")
	sb.WriteString(filepath.Join(configPath, "lib"))
	sb.WriteString("`\n\n")

	sb.WriteString("Cache path `")
	sb.WriteString(filepath.Join(configPath, "cache"))
	sb.WriteString("`\n\n")

	width := util.TermWidth()

	r, _ := glamour.NewTermRenderer(
		glamour.WithAutoStyle(),
		glamour.WithWordWrap(width),
	)

	out, err := r.Render(sb.String())
	if err != nil {
		panic(err)
	}
	print(out)
}

func ListConfigs() {
	de, err := os.ReadDir(util.ConfigPath)
	if err != nil {
		log.Fatal(err)
	}

	out := []string{}

	for _, d := range de {
		if d.IsDir() {
			if d.Name() == "cache" || d.Name() == "lib" {
				continue
			}
			err := filepath.WalkDir(filepath.Join(util.ConfigPath, d.Name()), func(path string, d os.DirEntry, err error) error {
				if err != nil {
					return err
				}
				if strings.HasSuffix(d.Name(), ".lua") {
					if rel, err := filepath.Rel(util.ConfigPath, path); err != nil {
						log.Fatal(err)
					} else {
						out = append(out, rel[0:len(rel)-4])
					}
				}
				return nil
			})
			if err != nil {
				log.Fatal(err)
			}
		}
		if strings.HasSuffix(d.Name(), ".lua") {
			out = append(out, d.Name()[0:len(d.Name())-4])
		}
	}

	os.Stdout.WriteString(strings.Join(out, ";"))
}

var attached = make([]string, 0)

func Attach(cmd string) {
	attached = append(attached, cmd)
}
func GetAttached() []string {
	return attached
}
func ClearAttached() {
	attached = make([]string, 0)
}
