package app

import (
	"testing"
)

func TestWatch(t *testing.T) {
	// testWatchLog := NewWatchLog("../test")
	// defer testWatchLog.watch.Close()
	// // testProcessLog := NewProcessLog(testWatchLog.logFilePathCh)
	// testWatchLog.WatchDir()
	// for path := range testWatchLog.logFilePathCh {
	// 	testProcessLog := NewProcessLog(path)
	// 	testProcessLog.Process()
	// }

	StartProcess("../test/")
}

func TestTailLog(t *testing.T) {
	testProcessLog := NewProcessLog("../test/access.log")
	testProcessLog.Process()
	select {}
}
