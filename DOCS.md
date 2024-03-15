# Documentation v0.4.1

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

*Makes the specified package available for use in the script. Currently supported packages are: golang, nodejs, mingw, msys2.*

**Parameters:** 
* package `"golang"|"nodejs"|"mingw"|"msys2"`
* version `string`

**Returns:** None

**Example:**

```lua
yab.use("golang", "1.22.0")
yab.use("nodejs", "14.17.6")
yab.use("msys2", "2024-01-13")
```

### 𝑓 task

*Checks if the given task is up to date and if not, executes the given task. This is useful for incremental builds.*

**Parameters:** 
* src `any`
* out `any`
* tool `function|table`

**Returns:** true if the toolchain was executed, false otherwise.

**Example:**

```lua
yab.task({ "foo.c" }, { "foo.o" }, function()
	os.execute("gcc -c foo.c -o foo.o")
end)
```

### 𝑓 os_type

*Returns the operating system type.*

**Parameters:** None

**Returns:** "windows", "linux" or "darwin" on the respective system.

### 𝑓 os_arch

*Returns the operating system architecture.*

**Parameters:** None

**Returns:** "amd64" or "arm64" on the respective system.

### 𝑓 setenv

*Sets an environment variable.*

**Parameters:** 
* key `string`
* value `string`

**Returns:** None

**Example:**

```lua
yab.setenv("FOO", "bar")
```

### 𝑓 args

*Returns the command line arguments passed to the program.*

**Parameters:** None

**Returns:** A table containing the command line arguments.

### 𝑓 cd

*Changes the current working directory to the given path for one function call.*

**Parameters:** 
* path `string`
* fn `function`

**Returns:** None

### 𝑓 mkdir

*Creates a new directory.*

**Parameters:** 
* path `string`

**Returns:** None

**Example:**

```lua
yab.mkdir('foo')
```

### 𝑓 rm

*Removes a file or directory.*

**Parameters:** 
* path `string`

**Returns:** None

**Example:**

```lua
yab.rm("./foo/bar")
```

### 𝑓 rm

*Removes a file or directory.*

**Parameters:** 
* path `string`
* recursive `boolean`

**Returns:** None

**Example:**

```lua
yab.rm("./foo/bar", true)
```

### 𝑓 check_exec

*Checks if an executable is available in the system's PATH.*

**Parameters:** 
* executable `string`

**Returns:** true if the executable is available, false otherwise.

### 𝑓 stdall

*Call a shell command and return the full output (stdout + stderr) in one string.*

**Parameters:** 
* command `string`

**Returns:** The output of the command.

### 𝑓 stdout

*Call a shell command and return the output (stdout) in one string.*

**Parameters:** 
* command `string`

**Returns:** The output of the command.

### 𝑓 stderr

*Call a shell command and return the error output (stderr) in one string.*

**Parameters:** 
* command `string`

**Returns:** The output of the command.

### 𝑓 git_clone_or_pull

*Clones a git repository to a specified destination. If the repository already exists, it will pull the latest changes instead.*

**Parameters:** 
* url `string`
* destination `string`

**Returns:** None

### 𝑓 zip

*Create a zip file containing the given files.*

**Parameters:** 
* files `table`
* output `string`

**Returns:** None

**Example:**

```lua
yab.zip({ "foo.txt", "bar.txt", "baz/" }, "archive.zip")
```

### 𝑓 download

*Download a file from the internet.*

**Parameters:** 
* url `string`

**Returns:** The name of the downloaded file.

**Example:**

```lua
yab.download("https://example.com/foo.txt")
```

### 𝑓 download

*Download a file from the internet to a specified destination.*

**Parameters:** 
* url `string`
* destination `string`

**Returns:** The name of the downloaded file.

**Example:**

```lua
yab.download("https://example.com/foo.txt", "foo.txt")
```

### 𝑓 watch

*Watch file or directory paths for changes and call a function when a change occurs. The callback function will be called with the file path and the event type as arguments. The event type can be one of 'create', 'write', 'remove', 'rename' or 'chmod'.*

**Parameters:** 
* paths `table`
* callback `function(string, string)`

**Returns:** None

**Example:**

```lua
yab.watch("foo.txt", function(file, event)
	print("foo.txt changed!")
end)
```

### 𝑓 block

*Block the current thread and wait for an interrupt signal.*

**Parameters:** None

**Returns:** None

**Example:**

```lua
yab.block()
```

### 𝑓 find

*Find files matching a pattern in a directory.*

**Parameters:** 
* pattern `string`

**Returns:** A table containing the matching file paths.

**Example:**

```lua
yab.find("*.txt")
```

### 𝑓 find

*Find files matching a pattern in a directory.*

**Parameters:** 
* root `string`
* pattern `string`

**Returns:** A table containing the matching file paths.

**Example:**

```lua
yab.find("foo", "*.txt")
```

### 𝑓 fileinfo

*Get information about a file.*

**Parameters:** 
* path `string`

**Returns:** A table containing the file information (name, size, mode, modtime, isdir, sys). See https://pkg.go.dev/io/fs#FileInfo for details.

**Example:**

```lua
local foo_info = yab.fileinfo("foo.txt")
print(foo_info.size)
```

### 𝑓 pretty

*Pretty print a table.*

**Parameters:** 
* value `any`

**Returns:** A string representation of the table.

**Example:**

```lua
yab.pretty({foo = "bar", baz = "qux"})
```

### 𝑓 print

*Pretty print a table.*

**Parameters:** 
* value `any`

**Returns:** None

**Example:**

```lua
yab.print({foo = "bar", baz = "qux"})
```

