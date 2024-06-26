package cmd

import (
	"github.com/tsukinoko-kun/yab/internal/docs"
	"github.com/spf13/cobra"
)

var docsCmd = &cobra.Command{
	Use:     "docs",
	Aliases: []string{"help"},
	Short:   "Prints the documentation",
	Long:    "Prints the yab documentation for your installed version",
	Run: func(cmd *cobra.Command, args []string) {
		docs.Help()
	},
}

func init() {
	rootCmd.AddCommand(docsCmd)
}
