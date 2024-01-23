package main

import (
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
	cliArgs.Parse()

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
}
