package pretty

import (
	"fmt"
	"strings"

	lua "github.com/Frank-Mayer/gopher-lua"
)

func Pretty(l *lua.LState) int {
	val := l.CheckAny(1)
	if val.Type() == lua.LTTable {
		str := PrettyTable(val.(*lua.LTable), 0)
		l.Push(lua.LString(str))
	} else {
		str := val.String()
		l.Push(lua.LString(str))
	}
	return 1
}

func PrettyTable(tbl *lua.LTable, indent int) string {
	sb := strings.Builder{}

	sb.WriteString("{\n")
	tbl.ForEach(func(key lua.LValue, value lua.LValue) {
		sb.WriteString(strings.Repeat("  ", indent+1))
		if key.Type() == lua.LTNumber {
			sb.WriteString("[")
			sb.WriteString(key.String())
			sb.WriteString("] = ")
		} else {
			sb.WriteString(key.String())
			sb.WriteString(" = ")
		}
		switch value.Type() {
		case lua.LTTable:
			sb.WriteString(PrettyTable(value.(*lua.LTable), indent+2))
		case lua.LTString:
			sb.WriteString("\"" + value.String() + "\"")
		case lua.LTNumber:
			sb.WriteString(value.String())
		case lua.LTBool:
			sb.WriteString(value.String())
		case lua.LTFunction:
			f := value.(*lua.LFunction).Proto.SourceName
			n := value.(*lua.LFunction).Proto.LineDefined
			sb.WriteString(fmt.Sprintf("function :: %s:%d", f, n))
		case lua.LTNil:
			sb.WriteString("nil")
		default:
			sb.WriteString("unknown")
		}
		sb.WriteString(",\n")
	})
	sb.WriteString(strings.Repeat("  ", indent))
	sb.WriteString("}")
	return sb.String()
}

func PrintPretty(l *lua.LState) int {
	val := l.CheckAny(1)
	if val.Type() == lua.LTTable {
		str := PrettyTable(val.(*lua.LTable), 0)
		fmt.Println(str)
	} else {
		str := val.String()
		fmt.Println(str)
	}
	return 0
}
