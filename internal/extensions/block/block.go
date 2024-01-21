package block

import (
	"os"
	"os/signal"

	"github.com/Frank-Mayer/gopher-lua"
	"github.com/Frank-Mayer/yab/internal/util"
	"github.com/charmbracelet/log"
)

func Block(_ *lua.LState) int {
	stopSpinner := util.Spin()
	defer stopSpinner()

	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, os.Interrupt)
	<-sigchan

	log.Debug("Received interrupt signal, stog blocking")

	return 0
}
