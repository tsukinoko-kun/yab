//go:build !windows
// +build !windows

package mainutil

import "syscall"

func Prepare() {
	syscall.Umask(0)
}
