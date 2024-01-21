package util_test

import (
	"testing"

	"github.com/Frank-Mayer/yab/internal/util"
)

func TestGetGlobalConfigPath(t *testing.T) {
	util.ConfigPath = ""
	path, err := util.GetGlobalConfigPath()
	if err != nil {
		t.Error(err)
	}
	if path == "" {
		t.Error("Expected path to be non-empty, got empty string")
	}
}

func BenchmarkGetGlobalConfigPath(b *testing.B) {
	for i := 0; i < b.N; i++ {
		path, err := util.GetGlobalConfigPath()
		if err != nil {
			b.Error(err)
		} else if path == "" {
			b.Error("Expected path to be non-empty, got empty string")
		}
	}
}
