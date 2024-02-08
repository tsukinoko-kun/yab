package shell

import "os/exec"

var shell func(string) error = nil

func SetShell(f func(string) error) {
	shell = f
}

func UseSell(command string) (bool, error, int) {
	if shell == nil {
		return false, nil, 0
	}
	err := shell(command)
	if exitError, ok := err.(*exec.ExitError); ok {
		exitCode := exitError.ExitCode()
		return true, err, exitCode
	}
	return true, err, 0
}
