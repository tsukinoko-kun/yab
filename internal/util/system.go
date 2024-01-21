package util

import (
	"os/exec"
	"strings"
)

func System(cmd string) *exec.Cmd {
	parts := strings.Fields(cmd)
	command := parts[0]
	args := parts[1:]

	return exec.Command(command, args...)
}
