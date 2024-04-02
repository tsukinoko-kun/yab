package util

import (
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/Frank-Mayer/gopher-lua/shell"
)

func System(cmdStr string) *exec.Cmd {
	if sh, ok := shell.GetShell(); ok {
		return sh(cmdStr)
	}
	cmd, args := PopenArgs(cmdStr)
	return exec.Command(cmd, args...)
}

func PopenArgs(arg string) (string, []string) {
	cmd := "/bin/sh"
	args := []string{"-c"}
	if runtime.GOOS == "windows" {
		cmd = "C:\\Windows\\system32\\cmd.exe"
		args = []string{"/c"}

		a := strings.Builder{}
		quoted := false
		for i := 0; i < len(arg); i++ {
			c := arg[i]
		sw:
			switch c {
			case '%':
				// might be a environment variable like %PATH%
				v := strings.Builder{}
				for j := i + 1; j < len(arg); j++ {
					switch arg[j] {
					case '%':
						i = j
						a.WriteString(os.Getenv(v.String()))
						break sw
					case ' ':
						a.WriteByte(c)
						break sw
					default:
						v.WriteByte(arg[j])
					}
				}
			case '"':
				// differenciate "" as escape sequence inside quoted string and begin/end of quoted string
				if quoted {
					j := i + 1
					if j < len(arg) && arg[j] == '"' {
						a.WriteRune('"')
						i = j
					} else {
						quoted = false
					}
				} else {
					quoted = true
				}
			case ' ', '\t', '\n', '\r', '\v':
				if quoted {
					a.WriteByte(c)
				} else if a.Len() > 0 {
					args = append(args, a.String())
					a.Reset()
				}
			default:
				a.WriteByte(c)
			}
		}
		if a.Len() > 0 {
			args = append(args, a.String())
		}
	} else {
		args = append(args, arg)
	}
	return cmd, args
}
