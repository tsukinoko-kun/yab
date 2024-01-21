package task

import (
	"path/filepath"

	lua "github.com/Frank-Mayer/gopher-lua"
	"github.com/Frank-Mayer/yab/internal/cache"
	"github.com/Frank-Mayer/yab/internal/extensions/pretty"
	"github.com/charmbracelet/log"
)

func Task(l *lua.LState) int {
	in := anyArgString(l, l.CheckAny(1))
	for i, s := range in {
		if abs, err := filepath.Abs(s); err == nil {
			in[i] = abs
		}
	}
	out := anyArgString(l, l.CheckAny(2))
	for i, s := range out {
		if abs, err := filepath.Abs(s); err == nil {
			out[i] = abs
		}
	}
	toolVar := l.CheckAny(3)
	tool := anyArgFunc(l, toolVar)
	var toolStr string
	switch toolVar.Type() {
	case lua.LTTable:
		toolStr = pretty.PrettyTable(toolVar.(*lua.LTable), 0)
	case lua.LTFunction:
		toolStr = toolVar.(*lua.LFunction).Proto.String()
	default:
		toolStr = toolVar.String()
	}

	if upToDate, err, writeCaheFile := cache.LookupToolchain(in, out, toolStr); err != nil {
		l.RaiseError("Error looking up toolchain: %s", err.Error())
	} else {
		if upToDate {
			log.Info("Toolchain up to date")
			l.Push(lua.LFalse)
		} else {
			log.Debug("Toolchain not up to date")
			log.Debug("Executing toolchain", "count", len(tool))
			for i, f := range tool {
				log.Debug("Executing tool", "index", i)
				if err := f(); err != nil {
					l.RaiseError("Error executing tool: %s", err.Error())
					return 0
				}
			}
			// create cache file here using os
			log.Debug("Writing cache file")
			if err := writeCaheFile(); err != nil {
				l.RaiseError("Error writing cache file: %s", err.Error())
				return 0
			}
			// check again if toolchain is up to date
			if upToDate, err, _ := cache.LookupToolchain(in, out, toolStr); err != nil {
				l.RaiseError("Error looking up toolchain: %s", err.Error())
			} else {
				if upToDate {
					log.Info("Toolchain executed successfully")
				} else {
					l.RaiseError("Toolchain not up to date after execution. This should not happen. Check your toolchain.")
				}
			}
			l.Push(lua.LTrue)
		}
		return 1
	}

	return 0
}

func anyArgString(l *lua.LState, val lua.LValue) []string {
	switch val.Type() {
	case lua.LTString:
		return []string{val.String()}
	case lua.LTBool:
		return []string{val.String()}
	case lua.LTNumber:
		return []string{val.String()}
	case lua.LTFunction:
		err := l.CallByParam(lua.P{
			Fn:      val.(*lua.LFunction),
			NRet:    0,
			Protect: true,
		})
		if err != nil {
			l.RaiseError("Error calling function: %s", err.Error())
		}
		return []string{}
	case lua.LTTable:
		t := val.(*lua.LTable)
		var ret []string
		t.ForEach(func(_ lua.LValue, v lua.LValue) {
			ret = append(ret, anyArgString(l, v)...)
		})
		return ret
	default:
		return []string{}
	}
}

type fn func() error

func anyArgFunc(l *lua.LState, val lua.LValue) []fn {
	switch val.Type() {
	case lua.LTString:
		return []fn{func() error {
			return l.DoString(val.String())
		}}
	case lua.LTFunction:
		return []fn{func() error {
			return l.CallByParam(lua.P{
				Fn:      val.(*lua.LFunction),
				NRet:    0,
				Protect: true,
			})
		}}
	case lua.LTTable:
		t := val.(*lua.LTable)
		var ret []fn
		t.ForEach(func(_ lua.LValue, v lua.LValue) {
			ret = append(ret, anyArgFunc(l, v)...)
		})
		return ret
	default:
		return []fn{}
	}
}
