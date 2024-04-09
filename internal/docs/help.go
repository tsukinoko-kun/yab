package docs

import (
	"github.com/Frank-Mayer/yab/internal/extensions"
	"github.com/Frank-Mayer/yab/internal/util"

	"strings"

	"github.com/charmbracelet/glamour"
)

func Help() {
	width := util.TermWidth()

	r, _ := glamour.NewTermRenderer(
		glamour.WithAutoStyle(),
		glamour.WithWordWrap(width),
	)

	out, err := r.Render(Markdown())
	if err != nil {
		panic(err)
	}
	print(out)
}

func Markdown() string {
	var sb strings.Builder
	sb.WriteString("# Documentation v")
	sb.WriteString(util.Version)
	sb.WriteString("\n\n")

	binName := util.BinName()

	sb.WriteString("## Usage\n\n")
	sb.WriteString("**`" + binName + " run [configs ...]`**\n\n")
	sb.WriteString("**`" + binName + " run [configs ...] -- [args ...]`**\n\n")
	sb.WriteString("Configs are Lua files in your local `.yab` folder or in the global config folder.\n\n")
	sb.WriteString("**`" + binName + " [configs ...] --attach [command]`**\n\n")
	sb.WriteString("Attaches a command to the yab environment after running all given configs.\n\n")
	sb.WriteString("**`" + binName + " version`**\n\n")
	sb.WriteString("Prints the version of the program.\n\n")
	sb.WriteString("**`" + binName + " docs`**\n\n")
	sb.WriteString("Prints this documentation.\n\n")
	sb.WriteString("**`" + binName + " def`**\n\n")
	sb.WriteString("Creates definitions file in global config.\n\n")
	sb.WriteString("**`" + binName + " env`**\n\n")
	sb.WriteString("Prints the yab environment.\n\n")

	sb.WriteString("### Flags\n\n")
	sb.WriteString("**`--debug`**\n\n")
	sb.WriteString("Enables debug logging.\n\n")
	sb.WriteString("**`--silent`**\n\n")
	sb.WriteString("Disables logging.\n\n")

	sb.WriteString("## Command Line Arguments\n\n")

	sb.WriteString("## Lua API Functions (in the `yab` module)\n\n")
	for _, f := range extensions.Functions {
		sb.WriteString(addFunction(&f))
		sb.WriteString("\n")
	}
	return sb.String()
}

func addFunction(f *extensions.Function) string {
	var sb strings.Builder
	sb.WriteString("### 𝑓 " + f.Name + "\n\n")
	sb.WriteString("*" + f.Description + "*\n\n")
	sb.WriteString("**Parameters:** ")
	if len(f.Parameters) > 0 {
		sb.WriteString("\n")
		for _, p := range f.Parameters {
			sb.WriteString("* ")
			paramWords := strings.Split(p, " ")
			sb.WriteString(paramWords[0])
			sb.WriteString(" `")
			sb.WriteString(strings.Join(paramWords[1:], " "))
			sb.WriteString("`\n")
		}
	} else {
		sb.WriteString("None\n")
	}
	sb.WriteString("\n")
	sb.WriteString("**Returns:** ")
	switch len(f.Returns) {
	case 0:
		sb.WriteString("None\n")
	case 1:
		sb.WriteString(f.Returns[0] + "\n")
	default:
		sb.WriteString("\n")
		for _, r := range f.Returns {
			sb.WriteString("* " + r + "\n")
		}
	}

	if f.Example != "" {
		sb.WriteString("\n**Example:**\n\n```lua\n")
		sb.WriteString(f.Example)
		sb.WriteString("\n```\n")
	}

	return sb.String()
}
