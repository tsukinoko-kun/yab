package main

import (
	"os"
	"os/exec"

	"github.com/Frank-Mayer/yab/internal/cli"
	"github.com/Frank-Mayer/yab/internal/mainutil"
	"github.com/Frank-Mayer/yab/internal/util"
	"github.com/charmbracelet/log"
)

func main() {
	var err error

	mainutil.Prepare()

	log.SetLevel(log.WarnLevel)

	defer util.RestoreEnv()

	if util.ConfigPath, err = mainutil.GetConfigPath(); err != nil {
		log.Fatal(err)
	}

	cliArgs := cli.Cli{}
	if err := cliArgs.Parse(); err != nil {
		log.Fatal("Failed to parse command line arguments", "error", err)
	}

	for _, file := range cliArgs.Configs {
		initFile, err := mainutil.GetInitFile(util.ConfigPath, file)
		if err != nil {
			log.Fatal(err)
		}
		err = mainutil.RunLuaFile(initFile)
		if err != nil {
			log.Fatal("Error running file: "+file, "error", err)
		}
	}

	for _, attached := range mainutil.GetAttached() {
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
}
