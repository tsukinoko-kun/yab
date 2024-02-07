package cli

import (
	"fmt"
	"os"

	"github.com/Frank-Mayer/yab/internal/docs"
	"github.com/Frank-Mayer/yab/internal/extensions/args"
	"github.com/Frank-Mayer/yab/internal/mainutil"
	"github.com/Frank-Mayer/yab/internal/util"
	"github.com/charmbracelet/log"
	hash "github.com/segmentio/fasthash/fnv1a"
)

type Cli struct {
	Configs []string
}

// hashes of the command line arguments
// this is 10% faster than using a switch statement on the string
const (
	// -v
	_v = 0x59cd76ca
	// --version
	__version = 0x2d151c11
	// -h
	_h = 0x5bcd79f0
	// --help
	__help = 0x7fb28c5c
	// --def
	__def = 0x2ee45cee
	// --debug
	__debug = 0x149e2a4e
	// --silent
	__silent = 0xd979d3e6
	// --
	__ = 0x20cd1d0f
	// --env
	__env = 0xa2734492
	// --ls
	__ls = 0xf1e7f19e
)

func (c *Cli) Parse() error {
	// attach loop
	for i := 0; i < len(os.Args); {
		if os.Args[i] == "--attach" {
			if len(os.Args) <= i+1 {
				return fmt.Errorf("No argument provided for --attach")
			}
			attach := os.Args[i+1]
			// remove the --attach and the argument
			os.Args = append(os.Args[:i], os.Args[i+2:]...)
			mainutil.Attach(attach)
		} else {
			i++
		}
	}
argLoop:
	for i, arg := range os.Args[1:] {
		switch hash.HashString32(arg) {
		case _v, __version:
			fmt.Println(util.Version)
		case _h, __help:
			docs.Help()
		case __def:
			mainutil.InitDefinitons()
		case __debug:
			log.SetLevel(log.DebugLevel)
			log.Debug("Debug mode enabled")
		case __silent:
			log.SetLevel(10)
		case __env:
			mainutil.PrintEnv()
		case __ls:
			mainutil.ListConfigs()
		case __:
			args.SetArgs(os.Args[i+2:])
			break argLoop
		default:
			c.Configs = append(c.Configs, arg)
		}
	}

	return nil
}
