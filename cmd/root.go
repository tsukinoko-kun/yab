package cmd

import (
	"os"

	"github.com/tsukinoko-kun/yab/internal/mainutil"
	"github.com/tsukinoko-kun/yab/internal/util"
	"github.com/charmbracelet/log"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "yab",
	Short: "Yet another build tool",
	Long: `Wouldn't it be great if you could use the same build tool for every project?
Regardless of operating system, programming language...
Yab is just that.
Use Lua scripts to define specific actions and execute them from the command line.`,
	PersistentPreRunE: func(cmd *cobra.Command, _ []string) error {
		mainutil.Prepare()

		if cmd.Flags().Changed("debug") {
			log.SetLevel(log.DebugLevel)
		} else if cmd.Flags().Changed("silent") {
			log.SetLevel(log.ErrorLevel)
		} else {
			log.SetLevel(log.WarnLevel)
		}

		var err error
		util.ConfigPath, err = mainutil.GetConfigPath()
		return err
	},
	PersistentPostRun: func(_ *cobra.Command, _ []string) {
		util.RestoreEnv()
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().Bool("debug", false, "Enable debug output")
	rootCmd.PersistentFlags().Bool("silent", false, "Silent output")
}
