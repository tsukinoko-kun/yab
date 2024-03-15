package cmd

import (
	"github.com/Frank-Mayer/yab/internal/docs"
	"github.com/spf13/cobra"
)

// docsCmd represents the docs command
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
