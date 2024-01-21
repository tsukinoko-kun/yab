package test_test

import (
	"testing"

	"github.com/Frank-Mayer/yab/internal/mainutil"
	"github.com/Frank-Mayer/yab/internal/util"
	"github.com/charmbracelet/log"
)

// func Test_Main(t *testing.T) {
// 	log.SetLevel(log.ErrorLevel)
// 	defer util.RestorePath()
//
// 	var err error
// 	if util.ConfigPath, err = mainutil.GetConfigPath(); err != nil {
// 		log.Fatal(err)
// 	}
//
// 	file := "build"
// 	initFile, err := mainutil.GetInitFile(util.ConfigPath, file)
// 	if err != nil {
// 		t.Fail()
// 		return
// 	}
// 	err = mainutil.RunLuaFile(initFile)
// 	if err != nil {
// 		t.Fail()
// 		return
// 	}
// }

func Benchmark_Main(b *testing.B) {
    log.SetLevel(log.ErrorLevel)
	for i := 0; i < b.N; i++ {
        var err error
		if util.ConfigPath, err = mainutil.GetConfigPath(); err != nil {
			log.Fatal(err)
		}

		file := "build"
		initFile, err := mainutil.GetInitFile(util.ConfigPath, file)
		if err != nil {
			b.Fail()
			return
		}
		err = mainutil.RunLuaFile(initFile)
		if err != nil {
			b.Fail()
			return
		}
        util.RestorePath()
	}
}
