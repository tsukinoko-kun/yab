# Documentation v0.4.0

## Usage

**yab run [configs ...]**

**yab run [configs ...] -- [args ...]**

Configs are Lua files in your local `.yab` folder or in the global config folder.

**yab [configs ...] --attach [command]**

Attaches a command to the yab environment after running all given configs.

**yab version**

Prints the version of the program.

**yab docs**

Prints this documentation.

**yab def**

Creates definitions file in global config.

**yab env**

Prints the yab environment.

### Flags

**--debug**

Enables debug logging.

**--silent**

Disables logging.

## Command Line Arguments

## Lua API Functions (in the `yab` module)

### 𝑓 use

_Makes the specified package available for use in the script. Currently supported packages are: golang, nodejs, mingw, msys2._

**Parameters:**

- package `"golang"|"nodejs"|"mingw"|"msys2"`
- version `string`

**Returns:** None

**Example:**

```lua
yab.use("golang", "1.22.0")
yab.use("nodejs", "14.17.6")
yab.use("msys2", "2024-01-13")
```

### 𝑓 task

_Checks if the given task is up to date and if not, executes the given task. This is useful for incremental builds._

**Parameters:**

- src `any`
- out `any`
- tool `function|table`

**Returns:** true if the toolchain was executed, false otherwise.

**Example:**

```lua
yab.task({ "foo.c" }, { "foo.o" }, function()
	os.execute("gcc -c foo.c -o foo.o")
end)
```

### 𝑓 os_type

_Returns the operating system type._

**Parameters:** None

**Returns:** "windows", "linux" or "darwin" on the respective system.

### 𝑓 os_arch

_Returns the operating system architecture._

**Parameters:** None

**Returns:** "amd64" or "arm64" on the respective system.

### 𝑓 setenv

_Sets an environment variable._

**Parameters:**

- key `string`
- value `string`

**Returns:** None

**Example:**

```lua
yab.setenv("FOO", "bar")
```

### 𝑓 args

_Returns the command line arguments passed to the program._

**Parameters:** None

**Returns:** A table containing the command line arguments.

### 𝑓 cd

_Changes the current working directory to the given path for one function call._

**Parameters:**

- path `string`
- fn `function`

**Returns:** None

### 𝑓 mkdir

_Creates a new directory._

**Parameters:**

- path `string`

**Returns:** None

**Example:**

```lua
yab.mkdir('foo')
```

### 𝑓 rm

_Removes a file or directory._

**Parameters:**

- path `string`

**Returns:** None

**Example:**

```lua
yab.rm("./foo/bar")
```

### 𝑓 rm

_Removes a file or directory._

**Parameters:**

- path `string`
- recursive `boolean`

**Returns:** None

**Example:**

```lua
yab.rm("./foo/bar", true)
```

### 𝑓 check_exec

_Checks if an executable is available in the system's PATH._

**Parameters:**

- executable `string`

**Returns:** true if the executable is available, false otherwise.

### 𝑓 stdall

_Call a shell command and return the full output (stdout + stderr) in one string._

**Parameters:**

- command `string`

**Returns:** The output of the command.

### 𝑓 stdout

_Call a shell command and return the output (stdout) in one string._

**Parameters:**

- command `string`

**Returns:** The output of the command.

### 𝑓 stderr

_Call a shell command and return the error output (stderr) in one string._

**Parameters:**

- command `string`

**Returns:** The output of the command.

### 𝑓 git_clone_or_pull

_Clones a git repository to a specified destination. If the repository already exists, it will pull the latest changes instead._

**Parameters:**

- url `string`
- destination `string`

**Returns:** None

### 𝑓 zip

_Create a zip file containing the given files._

**Parameters:**

- files `table`
- output `string`

**Returns:** None

**Example:**

```lua
yab.zip({ "foo.txt", "bar.txt", "baz/" }, "archive.zip")
```

### 𝑓 download

_Download a file from the internet._

**Parameters:**

- url `string`

**Returns:** The name of the downloaded file.

**Example:**

```lua
yab.download("https://example.com/foo.txt")
```

### 𝑓 download

_Download a file from the internet to a specified destination._

**Parameters:**

- url `string`
- destination `string`

**Returns:** The name of the downloaded file.

**Example:**

```lua
yab.download("https://example.com/foo.txt", "foo.txt")
```

### 𝑓 watch

_Watch file or directory paths for changes and call a function when a change occurs. The callback function will be called with the file path and the event type as arguments. The event type can be one of 'create', 'write', 'remove', 'rename' or 'chmod'._

**Parameters:**

- paths `table`
- callback `function(string, string)`

**Returns:** None

**Example:**

```lua
yab.watch("foo.txt", function(file, event)
	print("foo.txt changed!")
end)
```

### 𝑓 block

_Block the current thread and wait for an interrupt signal._

**Parameters:** None

**Returns:** None

**Example:**

```lua
yab.block()
```

### 𝑓 find

_Find files matching a pattern in a directory._

**Parameters:**

- pattern `string`

**Returns:** A table containing the matching file paths.

**Example:**

```lua
yab.find("*.txt")
```

### 𝑓 find

_Find files matching a pattern in a directory._

**Parameters:**

- root `string`
- pattern `string`

**Returns:** A table containing the matching file paths.

**Example:**

```lua
yab.find("foo", "*.txt")
```

### 𝑓 fileinfo

_Get information about a file._

**Parameters:**

- path `string`

**Returns:** A table containing the file information (name, size, mode, modtime, isdir, sys). See https://pkg.go.dev/io/fs#FileInfo for details.

**Example:**

```lua
local foo_info = yab.fileinfo("foo.txt")
print(foo_info.size)
```

### 𝑓 pretty

_Pretty print a table._

**Parameters:**

- value `any`

**Returns:** A string representation of the table.

**Example:**

```lua
yab.pretty({foo = "bar", baz = "qux"})
```

### 𝑓 print

_Pretty print a table._

**Parameters:**

- value `any`

**Returns:** None

**Example:**

```lua
yab.print({foo = "bar", baz = "qux"})
```
