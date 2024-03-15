package cmd

import (
	"os"
	"os/exec"

	argsext "github.com/Frank-Mayer/yab/internal/extensions/args"
	"github.com/Frank-Mayer/yab/internal/mainutil"
	"github.com/Frank-Mayer/yab/internal/util"
	"github.com/charmbracelet/log"
	"github.com/spf13/cobra"
)

var (
	runCmd = &cobra.Command{
		Use:   "run [config...]",
		Short: "Run a configuration",
		Long:  "Run a configuration lua file",
		RunE: func(self *cobra.Command, args []string) error {
			var err error

			mainutil.Prepare()

			if self.Flags().Changed("debug") {
				log.SetLevel(log.DebugLevel)
			} else if self.Flags().Changed("silent") {
				log.SetLevel(log.ErrorLevel)
			} else {
				log.SetLevel(log.WarnLevel)
			}

			defer util.RestoreEnv()

			if util.ConfigPath, err = mainutil.GetConfigPath(); err != nil {
				return err
			}

			var files []string
			argsLenAtDash := self.Flags().ArgsLenAtDash()
			if argsLenAtDash > 0 {
				files = args[:argsLenAtDash]
				argsext.SetArgs(args[argsLenAtDash:])
			} else {
				files = args
			}

			for _, file := range files {
				initFile, err := mainutil.GetInitFile(util.ConfigPath, file)
				if err != nil {
					return err
				}
				err = mainutil.RunLuaFile(initFile)
				if err != nil {
					log.Error("Error running file: " + file)
					return err
				}
			}

			attached := self.Flag("attach").Value.String()
			if attached != "" {
				log.Info("attaching", "command", attached)
				// execute the attached command
				cmd := exec.Command(attached)
				cmd.Stdout = os.Stdout
				cmd.Stderr = os.Stderr
				cmd.Stdin = os.Stdin
				if err := cmd.Run(); err != nil {
					log.Fatal("Error running attached command", "error", err)
				}
			}

			return nil
		},
	}
)

func init() {
	runCmd.Flags().String("attach", "", "Attach a command to run after the configuration files.")
	runCmd.Flags().Bool("debug", false, "Enable debug logging.")
	runCmd.Flags().Bool("silent", false, "Disable logging.")
	rootCmd.AddCommand(runCmd)
}
