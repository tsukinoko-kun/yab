package cli_test

import (
	"os"
	"strings"
	"testing"

	"github.com/Frank-Mayer/yab/internal/cli"
	"github.com/Frank-Mayer/yab/internal/mainutil"
)

func TestCli_Parse(t *testing.T) {
	os.Args = []string{"yab", "test", "test2", "test3"}

	c := cli.Cli{}
	if err := c.Parse(); err != nil {
		t.Errorf("Unexpected error during c.Parse(): %s", err)
	}

	if len(c.Configs) != 3 {
		t.Errorf("Expected len(c.Configs) to be 3, got %d", len(c.Configs))
	}
	if c.Configs[0] != "test" {
		t.Errorf("Expected c.Configs[0] to be 'test', got %s", c.Configs[0])
	}
	if c.Configs[1] != "test2" {
		t.Errorf("Expected c.Configs[1] to be 'test2', got %s", c.Configs[1])
	}
	if c.Configs[2] != "test3" {
		t.Errorf("Expected c.Configs[2] to be 'test3', got %s", c.Configs[2])
	}
}

func TestCli_Parse2(t *testing.T) {
	os.Args = []string{"yab", "test", "test2", "test3", "--", "test4", "test5", "test6"}

	c := cli.Cli{}
	if err := c.Parse(); err != nil {
		t.Errorf("Unexpected error during c.Parse(): %s", err)
	}

	if len(c.Configs) != 3 {
		t.Errorf("Expected len(c.Configs) to be 3, got %d", len(c.Configs))
	}
	if c.Configs[0] != "test" {
		t.Errorf("Expected c.Configs[0] to be 'test', got %s", c.Configs[0])
	}
	if c.Configs[1] != "test2" {
		t.Errorf("Expected c.Configs[1] to be 'test2', got %s", c.Configs[1])
	}
	if c.Configs[2] != "test3" {
		t.Errorf("Expected c.Configs[2] to be 'test3', got %s", c.Configs[2])
	}
}

func TestAttach(t *testing.T) {
	t.Run("NoAttach", func(t *testing.T) {
		defer mainutil.ClearAttached()
		os.Args = []string{"yab"}
		c := cli.Cli{}
		if err := c.Parse(); err != nil {
			t.Errorf("Unexpected error during c.Parse(): %s", err)
		}
		attached := mainutil.GetAttached()
		if len(attached) != 0 {
			t.Errorf("Expected len(attached) to be 0, got %s", vecStr(attached))
		}
	})
	t.Run("Attach one", func(t *testing.T) {
		defer mainutil.ClearAttached()
		os.Args = []string{"yab", "--attach", "vim"}
		c := cli.Cli{}
		if err := c.Parse(); err != nil {
			t.Errorf("Unexpected error during c.Parse(): %s", err)
		}
		attached := mainutil.GetAttached()
		if len(attached) != 1 {
			t.Errorf("Expected [vim], got %s", vecStr(attached))
		}
		if attached[0] != "vim" {
			t.Errorf("Expected attached[0] to be 'vim', got %s", attached[0])
		}
	})
	t.Run("Attach two", func(t *testing.T) {
		defer mainutil.ClearAttached()
		os.Args = []string{"yab", "--attach", "vim", "--attach", "emacs"}
		c := cli.Cli{}
		if err := c.Parse(); err != nil {
			t.Errorf("Unexpected error during c.Parse(): %s", err)
		}
		attached := mainutil.GetAttached()
		if len(attached) != 2 {
			t.Errorf("Expected [vim,emacs], got %s", vecStr(attached))
		}
		if attached[0] != "vim" {
			t.Errorf("Expected attached[0] to be 'vim', got %s", attached[0])
		}
		if attached[1] != "emacs" {
			t.Errorf("Expected attached[1] to be 'emacs', got %s", attached[1])
		}
	})
	t.Run("Attach false", func(t *testing.T) {
		defer mainutil.ClearAttached()
		os.Args = []string{"yab", "--attach"}
		c := cli.Cli{}
		if err := c.Parse(); err == nil {
			t.Errorf("Expected error, got nil")
		}
		attached := mainutil.GetAttached()
		if len(attached) != 0 {
			t.Errorf("Expected len(attached) to be 0, got %s", vecStr(attached))
		}
	})
	t.Run("Attach with config", func(t *testing.T) {
		defer mainutil.ClearAttached()
		os.Args = []string{"yab", "test", "test2", "test3", "--attach", "vim", "--", "test4", "test5", "test6"}
		c := cli.Cli{}
		if err := c.Parse(); err != nil {
			t.Errorf("Unexpected error during c.Parse(): %s", err)
		}
		attached := mainutil.GetAttached()
		if len(attached) != 1 {
			t.Errorf("Expected [vim], got %s", vecStr(attached))
		}
		if attached[0] != "vim" {
			t.Errorf("Expected attached[0] to be 'vim', got %s", attached[0])
		}

		if len(c.Configs) != 3 {
			t.Errorf("Expected len(c.Configs) to be 3, got %s", vecStr(c.Configs))
		}
		if c.Configs[0] != "test" {
			t.Errorf("Expected c.Configs[0] to be 'test', got %s", c.Configs[0])
		}
		if c.Configs[1] != "test2" {
			t.Errorf("Expected c.Configs[1] to be 'test2', got %s", c.Configs[1])
		}
		if c.Configs[2] != "test3" {
			t.Errorf("Expected c.Configs[2] to be 'test3', got %s", c.Configs[2])
		}
	})
}

func BenchmarkCli_Parse(b *testing.B) {
	for i := 0; i < b.N; i++ {
		os.Args = []string{"yab", "test", "test2", "test3", "--", "test4", "test5", "test6"}
		c := cli.Cli{}
		if err := c.Parse(); err != nil {
			b.Errorf("Unexpected error during c.Parse(): %s", err)
		}
	}
}

func vecStr(v []string) string {
	return "[" + strings.Join(v, ",") + "]"
}
