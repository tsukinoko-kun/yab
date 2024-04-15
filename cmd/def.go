package cmd

import (
	"github.com/tsukinoko-kun/yab/internal/mainutil"
	"github.com/spf13/cobra"
)

var defCmd = &cobra.Command{
	Use:   "def",
	Short: "Create a definitions file in your global config directory",
	Long:  "Create a definitions file in your global config directory. Use this to configure your Lua language server.",
	Run: func(_ *cobra.Command, _ []string) {
		mainutil.InitDefinitons()
	},
}

func init() {
	rootCmd.AddCommand(defCmd)
}
