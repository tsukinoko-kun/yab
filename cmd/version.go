package cmd

import (
	"fmt"

	"github.com/tsukinoko-kun/yab/internal/util"
	"github.com/spf13/cobra"
)

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
