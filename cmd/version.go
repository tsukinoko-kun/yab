package cmd

import (
	"fmt"

	"github.com/Frank-Mayer/yab/internal/util"
	"github.com/spf13/cobra"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of yab",
	Long:  "Print the version number of yab",
	Run: func(_ *cobra.Command, _ []string) {
		fmt.Println(util.Version)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
