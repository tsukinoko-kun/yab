package std

import (
	"bufio"
	"strings"
	"sync"

	"github.com/Frank-Mayer/yab/internal/lua"
	"github.com/Frank-Mayer/yab/internal/util"
)

// call a shell command and return the full output (stdout + stderr) in one string
func All(l *lua.LState) int {
	command := l.CheckString(1)

	cmd := util.System(command)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		l.Error(lua.LString("Error creating stdout pipe. "+err.Error()), 0)
		return 0
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		l.Error(lua.LString("Error creating stderr pipe. "+err.Error()), 0)
		return 0
	}

	if err := cmd.Start(); err != nil {
		l.Error(lua.LString("Error starting command. "+err.Error()), 0)
		return 0
	}

	sb := strings.Builder{}
	var wg sync.WaitGroup

	wg.Add(2)

	go func() {
		defer wg.Done()
		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			sb.WriteString(scanner.Text() + "\n")
		}
	}()

	go func() {
		defer wg.Done()
		scanner := bufio.NewScanner(stderr)
		for scanner.Scan() {
			sb.WriteString(scanner.Text() + "\n")
		}
	}()
	wg.Wait()

	if err := cmd.Wait(); err != nil {
		l.Error(lua.LString("Error executing command. "+err.Error()), 0)
		return 0
	}

	l.Push(lua.LString(strings.TrimSpace(sb.String())))
	return 1
}

// call a shell command and return the output (stdout) in one string
func Out(l *lua.LState) int {
	command := l.CheckString(1)

	cmd := util.System(command)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		l.Error(lua.LString("Error creating stdout pipe. "+err.Error()), 0)
		return 0
	}

	if err := cmd.Start(); err != nil {
		l.Error(lua.LString("Error starting command. "+err.Error()), 0)
		return 0
	}

	sb := strings.Builder{}
	var wg sync.WaitGroup

	wg.Add(1)

	go func() {
		defer wg.Done()
		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			sb.WriteString(scanner.Text() + "\n")
		}
	}()

	wg.Wait()

	if err := cmd.Wait(); err != nil {
		l.Error(lua.LString("Error executing command. "+err.Error()), 0)
		return 0
	}

	l.Push(lua.LString(strings.TrimSpace(sb.String())))
	return 1
}

// call a shell command and return the error (stderr) in one string
func Err(l *lua.LState) int {
	command := l.CheckString(1)

	cmd := util.System(command)

	stderr, err := cmd.StderrPipe()
	if err != nil {
		l.Error(lua.LString("Error creating stderr pipe. "+err.Error()), 0)
		return 0
	}

	if err := cmd.Start(); err != nil {
		l.Error(lua.LString("Error starting command. "+err.Error()), 0)
		return 0
	}

	sb := strings.Builder{}
	var wg sync.WaitGroup

	wg.Add(1)

	go func() {
		defer wg.Done()
		scanner := bufio.NewScanner(stderr)
		for scanner.Scan() {
			sb.WriteString(scanner.Text() + "\n")
		}
	}()
	wg.Wait()

	if err := cmd.Wait(); err != nil {
		l.Error(lua.LString("Error executing command. "+err.Error()), 0)
		return 0
	}

	l.Push(lua.LString(strings.TrimSpace(sb.String())))
	return 1
}
