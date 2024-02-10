package shell

import (
	"os"
	"os/exec"
)

var shell func(string) *exec.Cmd = nil

func SetShell(f func(string) *exec.Cmd) {
	shell = f
}

func UseSell(command string) (bool, error, int) {
	if shell == nil {
		return false, nil, 0
	}
	cmd := shell(command)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	if err := cmd.Run(); err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			return true, nil, exitError.ExitCode()
		} else {
			return true, err, 1
		}
	} else {
		return true, nil, 0
	}
}

func GetShell() (func(string) *exec.Cmd, bool) {
	if shell == nil {
		return nil, false
	}
	return shell, true
}
