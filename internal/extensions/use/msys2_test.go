package use_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/tsukinoko-kun/yab/internal/extensions/use"
)

func TestWinToPosixPath(t *testing.T) {
	tests := []struct {
		in, out string
	}{
		{"C:\\Program Files\\Git\\cmd", "/c/Program Files/Git/cmd"},
		{"C:\\Program Files\\Git\\cmd;C:\\Program Files\\Git\\bin", "/c/Program Files/Git/cmd:/c/Program Files/Git/bin"},
		{"C:\\Program Files\\Git\\cmd;C:\\Program Files\\Git\\bin;C:\\Program Files\\Git\\usr\\bin", "/c/Program Files/Git/cmd:/c/Program Files/Git/bin:/c/Program Files/Git/usr/bin"},
		{"C:\\Program Files\\Git\\cmd;C:\\Program Files\\Git\\bin;C:\\Program Files\\Git\\usr\\bin;C:\\Program Files\\Git\\usr\\bin\\vendor\\bin", "/c/Program Files/Git/cmd:/c/Program Files/Git/bin:/c/Program Files/Git/usr/bin:/c/Program Files/Git/usr/bin/vendor/bin"},
		{"", ""},
		{"C:\\Program Files\\Alacritty\\;C:\\Windows\\system32;C:\\Windows;C:\\Windows\\System32\\Wbem;C:\\Windows\\System32\\WindowsPowerShell\\v1.0\\;C:\\Windows\\System32\\OpenSSH\\;C:\\Program Files\\NVIDIA Corporation\\NVIDIA NvDLISR;C:\\Program Files\\dotnet\\;C:\\Program Files (x86)\\Windows Kits\\10\\Windows Performance Toolkit\\;C:\\Program Files\\Go\\bin;C:\\Program Files\\Microsoft VS Code\\bin;C:\\Program Files\\Git\\cmd;C:\\Program Files\\nodejs\\;C:\\Program Files\\LLVM\\bin;C:\\Users\\Shadow\\.cargo\\bin;C:\\Users\\Shadow\\AppData\\Local\\Programs\\Python\\Python312\\Scripts\\;C:\\Users\\Shadow\\AppData\\Local\\Programs\\Python\\Python312\\;C:\\Users\\Shadow\\AppData\\Local\\Microsoft\\WindowsApps;C:\\Users\\Shadow\\go\\bin;C:\\Users\\Shadow\\AppData\\Local\\JetBrains\\Toolbox\\scripts;C:\\Users\\Shadow\\.dotnet\\tools;C:\\Users\\Shadow\\go\\bin;C:\\Users\\Shadow\\Applications;C:\\Program Files\\Neovim\\bin;C:\\Users\\Shadow\\AppData\\Roaming\\npm",
			"/c/Program Files/Alacritty:/c/Windows/system32:/c/Windows:/c/Windows/System32/Wbem:/c/Windows/System32/WindowsPowerShell/v1.0:/c/Windows/System32/OpenSSH:/c/Program Files/NVIDIA Corporation/NVIDIA NvDLISR:/c/Program Files/dotnet:/c/Program Files (x86)/Windows Kits/10/Windows Performance Toolkit:/c/Program Files/Go/bin:/c/Program Files/Microsoft VS Code/bin:/c/Program Files/Git/cmd:/c/Program Files/nodejs:/c/Program Files/LLVM/bin:/c/Users/Shadow/.cargo/bin:/c/Users/Shadow/AppData/Local/Programs/Python/Python312/Scripts:/c/Users/Shadow/AppData/Local/Programs/Python/Python312:/c/Users/Shadow/AppData/Local/Microsoft/WindowsApps:/c/Users/Shadow/go/bin:/c/Users/Shadow/AppData/Local/JetBrains/Toolbox/scripts:/c/Users/Shadow/.dotnet/tools:/c/Users/Shadow/go/bin:/c/Users/Shadow/Applications:/c/Program Files/Neovim/bin:/c/Users/Shadow/AppData/Roaming/npm"},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("Test %d", i), func(t *testing.T) {
			got := use.WinToPosixPath(tt.in)
			errCount := 0
			for i, c := range got {
				if c != rune(tt.out[i]) {
					errCount++
					// print result and point out the difference
					// abcdefg
					//   ^ expected 'c' got 'd'

					start := i - 15
					if start < 0 {
						start = 0
					}
					end := i + 5
					if end > len(got) {
						end = len(got)
					}
					sb := strings.Builder{}
					sb.WriteString("\n")
					sb.WriteString(got[start:end])
					sb.WriteString("\n")
					for j := 0; j < i-start; j++ {
						sb.WriteString(" ")
					}
					sb.WriteString("^ expected '")
					sb.WriteString(string(tt.out[i]))
					sb.WriteString("' got '")
					sb.WriteString(string(c))
					sb.WriteString("'")
					t.Errorf(sb.String())

					if errCount > 5 {
						t.Errorf("Too many errors, aborting")
						break
					}
				}
			}
		})
	}
}
