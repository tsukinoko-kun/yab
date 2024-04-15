package parallel

import (
	"sync"

	"github.com/tsukinoko-kun/gopher-lua"
)

// Parallel executes a given set of functions in parallel.
func Parallel(l *lua.LState) int {
	n := l.GetTop()
	if n == 0 {
		l.RaiseError("No functions given to parallel()")
		return 0
	}

	var wg sync.WaitGroup
	for i := 1; i <= n; i++ {
		wg.Add(1)
		fn := l.CheckFunction(i)
		thread, _ := l.NewThread()
		go func() {
			defer wg.Done()
			_ = thread.CallByParam(lua.P{
				Fn:      fn,
				NRet:    0,
				Protect: true,
			})
			thread.Close()
		}()
	}

	wg.Wait()
	return 0
}
