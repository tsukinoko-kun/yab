package cmd

import (
	"github.com/tsukinoko-kun/yab/internal/mainutil"
	"github.com/spf13/cobra"
)

var envCmd = &cobra.Command{
	Use:   "env",
	Short: "Prints the yab environment",
	Long:  "Prints the yab environment",
	Run: func(cmd *cobra.Command, args []string) {
		mainutil.PrintEnv()
	},
}

func init() {
	rootCmd.AddCommand(envCmd)
}
