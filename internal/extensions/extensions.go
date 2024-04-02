package extensions

import (
	"strings"

	"github.com/Frank-Mayer/gopher-lua"
	"github.com/Frank-Mayer/yab/internal/extensions/args"
	"github.com/Frank-Mayer/yab/internal/extensions/block"
	"github.com/Frank-Mayer/yab/internal/extensions/cd"
	"github.com/Frank-Mayer/yab/internal/extensions/checkexec"
	"github.com/Frank-Mayer/yab/internal/extensions/download"
	"github.com/Frank-Mayer/yab/internal/extensions/env"
	"github.com/Frank-Mayer/yab/internal/extensions/fileinfo"
	"github.com/Frank-Mayer/yab/internal/extensions/find"
	"github.com/Frank-Mayer/yab/internal/extensions/git"
	"github.com/Frank-Mayer/yab/internal/extensions/mkdir"
	"github.com/Frank-Mayer/yab/internal/extensions/osarch"
	"github.com/Frank-Mayer/yab/internal/extensions/ostype"
	"github.com/Frank-Mayer/yab/internal/extensions/pretty"
	"github.com/Frank-Mayer/yab/internal/extensions/rm"
	"github.com/Frank-Mayer/yab/internal/extensions/std"
	"github.com/Frank-Mayer/yab/internal/extensions/task"
	"github.com/Frank-Mayer/yab/internal/extensions/use"
	"github.com/Frank-Mayer/yab/internal/extensions/watch"
	"github.com/Frank-Mayer/yab/internal/extensions/zip"
	"github.com/Frank-Mayer/yab/internal/util"
)

type Function struct {
	Name        string
	Description string
	Parameters  []string
	Returns     []string
	Function    func(l *lua.LState) int
	Ret         string
	Example     string
}

var Functions = []Function{
	{
		"use",
		"Makes the specified package available for use in the script. " +
			"Currently supported packages are: " +
			strings.Join(use.Usable, ", ") + ".",
		[]string{
			"package \"" + strings.Join(use.Usable, "\"|\"") + "\"",
			"version string"},
		[]string{},
		use.Use,
		"",
		"yab.use(\"golang\", \"1.22.0\")\n" +
			"yab.use(\"nodejs\", \"14.17.6\")\n" +
			"yab.use(\"msys2\", \"2024-01-13\")",
	},

	{
		"task",
		"Checks if the given task is up to date and if not, executes the given task. This is useful for incremental builds.",
		[]string{"src any", "out any", "tool function|table"},
		[]string{"true if the toolchain was executed, false otherwise."},
		task.Task,
		"boolean",
		"yab.task({ \"foo.c\" }, { \"foo.o\" }, function()\n" +
			"\tos.execute(\"gcc -c foo.c -o foo.o\")\n" +
			"end)",
	},

	{
		"os_type",
		"Returns the operating system type.",
		[]string{},
		[]string{"\"windows\", \"linux\" or \"darwin\" on the respective system."},
		ostype.OsType,
		"'windows'|'linux'|'darwin'",
		"",
	},

	{
		"os_arch",
		"Returns the operating system architecture.",
		[]string{},
		[]string{"\"amd64\" or \"arm64\" on the respective system."},
		osarch.OsArch,
		"'amd64'|'arm64'",
		"",
	},

	{
		"setenv",
		"Sets an environment variable.",
		[]string{"key string", "value string"},
		[]string{},
		env.SetEnv,
		"",
		`yab.setenv("FOO", "bar")`,
	},

	{
		"args",
		"Returns the command line arguments passed to the program.",
		[]string{},
		[]string{"A table containing the command line arguments."},
		args.Args,
		"table",
		"",
	},

	{
		"cd",
		"Changes the current working directory to the given path for one function call.",
		[]string{"path string", "fn function"},
		[]string{},
		cd.Cd,
		"",
		"",
	},

	{
		"mkdir",
		"Creates a new directory.",
		[]string{"path string"},
		[]string{},
		mkdir.Mkdir,
		"",
		"yab.mkdir('foo')",
	},

	{
		"rm",
		"Removes a file or directory.",
		[]string{"path string"},
		[]string{},
		rm.Rm,
		"",
		"yab.rm(\"./foo/bar\")",
	},
	{
		"rm",
		"Removes a file or directory.",
		[]string{"path string", "recursive boolean"},
		[]string{},
		rm.Rm,
		"",
		"yab.rm(\"./foo/bar\", true)",
	},

	{
		"check_exec",
		"Checks if an executable is available in the system's PATH.",
		[]string{"executable string"},
		[]string{"true if the executable is available, false otherwise."},
		checkexec.CheckExec,
		"boolean",
		"",
	},

	{
		"stdall",
		"Call a shell command and return the full output (stdout + stderr) in one string.",
		[]string{"command string"},
		[]string{"The output of the command."},
		std.All,
		"string",
		"",
	},

	{
		"stdout",
		"Call a shell command and return the output (stdout) in one string.",
		[]string{"command string"},
		[]string{"The output of the command."},
		std.Out,
		"string",
		"",
	},

	{
		"stderr",
		"Call a shell command and return the error output (stderr) in one string.",
		[]string{"command string"},
		[]string{"The output of the command."},
		std.Err,
		"string",
		"",
	},

	{
		"git_clone_or_pull",
		"Clones a git repository to a specified destination. If the repository already exists, it will pull the latest changes instead.",
		[]string{"url string", "destination string"},
		[]string{},
		git.GitCloneOrPull,
		"",
		"",
	},

	{
		"zip",
		"Create a zip file containing the given files.",
		[]string{"files table", "output string"},
		[]string{},
		zip.MakeZip,
		"",
		"yab.zip({ \"foo.txt\", \"bar.txt\", \"baz/\" }, \"archive.zip\")",
	},

	{
		"download",
		"Download a file from the internet.",
		[]string{"url string"},
		[]string{"The name of the downloaded file."},
		download.Download,
		"string",
		"yab.download(\"https://example.com/foo.txt\")",
	},
	{
		"download",
		"Download a file from the internet to a specified destination.",
		[]string{"url string", "destination string"},
		[]string{"The name of the downloaded file."},
		download.Download,
		"string",
		"yab.download(\"https://example.com/foo.txt\", \"foo.txt\")",
	},

	{
		"watch",
		"Watch file or directory paths for changes and call a function when a change occurs. " +
			"The callback function will be called with the file path and the event type as arguments. " +
			"The event type can be one of 'create', 'write', 'remove', 'rename' or 'chmod'.",
		[]string{"paths table", "callback function(string, string)"},
		[]string{},
		watch.Watch,
		"",
		"yab.watch(\"foo.txt\", function(file, event)\n\tprint(\"foo.txt changed!\")\nend)",
	},

	{
		"block",
		"Block the current thread and wait for an interrupt signal.",
		[]string{},
		[]string{},
		block.Block,
		"",
		"yab.block()",
	},

	{
		"find",
		"Find files matching a pattern in a directory.",
		[]string{"pattern string"},
		[]string{"A table containing the matching file paths."},
		find.Find,
		"table",
		"yab.find(\"*.txt\")",
	},
	{
		"find",
		"Find files matching a pattern in a directory.",
		[]string{"root string", "pattern string"},
		[]string{"A table containing the matching file paths."},
		find.Find,
		"table",
		"yab.find(\"foo\", \"*.txt\")",
	},

	{
		"fileinfo",
		"Get information about a file.",
		[]string{"path string"},
		[]string{"A table containing the file information (name, size, mode, modtime, isdir, sys). See https://pkg.go.dev/io/fs#FileInfo for details."},
		fileinfo.FileInfo,
		"table",
		"local foo_info = yab.fileinfo(\"foo.txt\")\nprint(foo_info.size)",
	},

	{
		"pretty",
		"Pretty print a table.",
		[]string{"value any"},
		[]string{"A string representation of the table."},
		pretty.Pretty,
		"string",
		"yab.pretty({foo = \"bar\", baz = \"qux\"})",
	},

	{
		"print",
		"Pretty print a table.",
		[]string{"value any"},
		[]string{},
		pretty.PrintPretty,
		"",
		"yab.print({foo = \"bar\", baz = \"qux\"})",
	},
}

func Definitions() string {
	sb := strings.Builder{}
	sb.WriteString("---@meta\n")
	sb.WriteString("---@class Yab\n")
	sb.WriteString("---@version " + util.Version + "\n")
	sb.WriteString("local yab = {}\n")
	for _, f := range Functions {
		sb.WriteString("\n")
		for _, p := range f.Parameters {
			sb.WriteString("---@param ")
			sb.WriteString(p)
			sb.WriteString("\n")
		}
		sb.WriteString("---@return ")
		sb.WriteString(f.Ret)
		sb.WriteString("\n")
		sb.WriteString("---")
		sb.WriteString(f.Description)
		sb.WriteString("\n")
		sb.WriteString("yab.")
		sb.WriteString(f.Name)
		sb.WriteString(" = function(")
		for i, p := range f.Parameters {
			if i > 0 {
				sb.WriteString(", ")
			}
			sb.WriteString(strings.Split(p, " ")[0])
		}
		sb.WriteString(")\nend\n")
	}
	sb.WriteString("\n")
	sb.WriteString("return yab\n")
	return sb.String()
}

func RegisterExtensions(l *lua.LState) {
	l.PreloadModule("yab", func(l *lua.LState) int {
		table := l.NewTable()
		for _, f := range Functions {
			l.SetTable(table, lua.LString(f.Name), l.NewFunction(f.Function))
		}
		l.Push(table)
		return 1
	})
}
