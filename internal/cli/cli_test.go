package cli_test

import (
	"os"
	"testing"

	"github.com/Frank-Mayer/yab/internal/cli"
)

func TestCli_Parse(t *testing.T) {
	os.Args = []string{"yab", "test", "test2", "test3"}

	c := cli.Cli{}
	c.Parse()

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
	c.Parse()

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

func BenchmarkCli_Parse(b *testing.B) {
	for i := 0; i < b.N; i++ {
		os.Args = []string{"yab", "test", "test2", "test3", "--", "test4", "test5", "test6"}
		c := cli.Cli{}
		c.Parse()
	}
}
