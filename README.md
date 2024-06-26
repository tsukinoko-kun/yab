# Yab

Yet another build tool :construction_worker: :construction:

[![Homebrew](https://img.shields.io/badge/Homebrew-FBB040?logo=homebrew&logoColor=000)](#homebrew)

[![Latest Release)](https://img.shields.io/github/v/release/tsukinoko-kun/yab?logo=github)](https://github.com/tsukinoko-kun/yab/releases/latest)

[![CodeQL](https://github.com/tsukinoko-kun/yab/actions/workflows/github-code-scanning/codeql/badge.svg)](https://github.com/tsukinoko-kun/yab/actions/workflows/github-code-scanning/codeql)

Wouldn't it be great if you could use the same build tool for every project?
Regardless of operating system, programming language...

Yab is just that.

Use Lua scripts to define specific actions and execute them from the command line.

> [!CAUTION]
> This project is still in very early stages of development.
> Expect breaking changes.
> Stable API is guaranteed at v1.0.0

## Does that not already exist?

No!

<table>
    <thead>
        <tr>
            <td></td>
            <td><sup>Builtin support for many technologies</sup></td>
            <td><sup>Smart incremental build</sup></td>
            <td><sup>Easy to setup and extend</sup></td>
            <td><sup>Imperative syntax (loops, functions, ...)</sup></td>
            <td><sup>Parameters</sup></td>
            <td><sup><b>No</b> domain specific language</sup></td>
            <td><sup>Cross-platform by default</sup></td>
        </tr>
    </thead>
    <tbody>
        <tr>
            <td>Yab</td>
            <td>:x: / :white_check_mark:</td>
            <td>:white_check_mark:</td>
            <td>:white_check_mark:</td>
            <td>:white_check_mark:</td>
            <td>:white_check_mark:</td>
            <td>:white_check_mark:</td>
            <td>:white_check_mark:</td>
        </tr>
        <tr>
            <td><a href="https://bazel.build/" target="blank">Bazel</a></td>
            <td>:white_check_mark:</td>
            <td>:white_check_mark:</td>
            <td>:x:</td>
            <td>:x:</td>
            <td>:x:</td>
            <td>:x:</td>
            <td>:white_check_mark:</td>
        </tr>
        <tr>
            <td><a href="https://gradle.org/" target="blank">Gradle</a></td>
            <td>:white_check_mark:</td>
            <td>:grey_question:</td>
            <td>:x:</td>
            <td>:x:</td>
            <td>:x:</td>
            <td>:x:</td>
            <td>:white_check_mark:</td>
        </tr>
        <tr>
            <td><a href="https://www.gnu.org/software/make/" target="blank">Make</a></td>
            <td>:x:</td>
            <td>:x:</td>
            <td>:x:</td>
            <td>:x:</td>
            <td>:x:</td>
            <td>:x:</td>
            <td>:x:</td>
        </tr>
    </tbody>
</table>

### Explanation

#### Built-in support for many technologies

How many technologies are supported by the tool?

#### Smart incremental build

Skip tasks that are already done.
This should be done using file content hashes.

#### Easy to set up and extend

How difficult is it to get started with the tool?

#### Imperative syntax

Use the tool to describe how the build process should be done.

#### Parameters

Optionally parametrize build configurations.

#### No domain specific language

Use a language to describe the build that is popular outside the tool.

#### Cross-platform by default

It is easy to get the build working on different platforms.

---

Lua is a common and performant programming language.
Yab offers some useful functions in addition to the Lua standard library that might be useful when building configurations.

Looking for an example configuration?
Take a look at [this projects `.yab` folder](https://github.com/tsukinoko-kun/yab/tree/main/.yab).

## Docs

Documentation is in the [DOCS.md](https://github.com/tsukinoko-kun/yab/blob/main/DOCS.md) file.

## Usage

Run one or more configs:

```bash
yab run [configs ...]
```

Pass arguments to the scripts:

```bash
yab run [configs ...] -- [args ...]
```

A config is a lua file inside the config directory.

The following directories are used as configs (first found wins)

1. `./.yab/`
1. `$XDG_CONFIG_HOME/yab/`
1. `$APPDATA/yab/`
1. `$HOME/.config/yab/`

If you run `yab run` without any arguments, the default config (`init.lua`) is used.

## Lua definitions

You can use a Lua language server to get autocompletion and type checking.

Run `yab def` to create a definitions file in your **global** config directory.
Add it to your Lua language server configuration.

Global config is one of those directories:

1. `$XDG_CONFIG_HOME/yab/`
1. `$APPDATA/yab/`
1. `$HOME/.config/yab/`

## Example Code

```lua
local yab = require("yab")

-- use a specific Go version
yab.use("golang", "1.21.6")

-- os specific name
local bin_name = yab.os_type() == "windows" and "yab.exe" or "yab"

-- incremental build
yab.task(yab.find("**.go"), bin_name, function()
	os.execute('go build -ldflags="-s -w" -o ' .. bin_name .. " ./cmd/yab/")
end)
```

## Attach external applications

I like to create an `env.lua` config where all shared tools are defined.

```lua
local yab = require("yab")
yab.use("golang", "1.21.6")
yab.use("nodejs", "20.11.1")
```

I now can start a shell using the specified environment using `yab run env --attach sh` or start neovim using `yab run env --attach nvim`.

## Installation

### Homebrew

```sh
brew tap tsukinoko-kun/tap
brew install tsukinoko-kun/tap/yab
```

### Go install

```bash
go install github.com/tsukinoko-kun/yab@latest
```

### Binary from release

[Latest release](https://github.com/tsukinoko-kun/yab/releases/latest)

## GitHub Actions

```yaml
- name: Use Yab
  uses: tsukinoko-kun/use-yab@v1
```

You can specify a version using the `version` input. The default is `latest`.

```yaml
- name: Use Yab
  uses: tsukinoko-kun/use-yab@v1
  with:
      version: 0.3.0
```

## Badge

Use the badge to show that your project uses Yab.
This helps to spread the word about Yab and makes it easier for others to work on your project.

This is just a suggestion.
You don't have to use it.

[![Yab Project](https://img.shields.io/badge/Yab_Project-2C2D72?logo=lua)](https://github.com/tsukinoko-kun/yab)

```markdown
[![Yab Project](https://img.shields.io/badge/Yab_Project-2C2D72?logo=lua)](https://github.com/tsukinoko-kun/yab)
```
