package use_test

import (
	"fmt"
	"testing"

	"github.com/Frank-Mayer/yab/internal/extensions/use"
)

func TestWinToPosixPath(t *testing.T) {
	tests := []struct {
		in, out string
	}{
		{"C:\\Program Files\\Git\\cmd", "/c/Program Files/Git/cmd"},
		{"C:\\Program Files\\Git\\cmd;C:\\Program Files\\Git\\bin", "/c/Program Files/Git/cmd:/c/Program Files/Git/bin"},
		{"C:\\Program Files\\Git\\cmd;C:\\Program Files\\Git\\bin;C:\\Program Files\\Git\\usr\\bin", "/c/Program Files/Git/cmd:/c/Program Files/Git/bin:/c/Program Files/Git/usr/bin"},
		{"C:\\Program Files\\Git\\cmd;C:\\Program Files\\Git\\bin;C:\\Program Files\\Git\\usr\\bin;C:\\Program Files\\Git\\usr\\bin\\vendor\\bin", "/c/Program Files/Git/cmd:/c/Program Files/Git/bin:/c/Program Files/Git/usr/bin:/c/Program Files/Git/usr/bin/vendor/bin"},
	}
	for i, tt := range tests {
		t.Run(fmt.Sprintf("Test %d", i), func(t *testing.T) {
			if got := use.WinToPosixPath(tt.in); got != tt.out {
				t.Errorf("WinToPosixPath(%q) = %q; want %q", tt.in, got, tt.out)
			}
		})
	}
}
