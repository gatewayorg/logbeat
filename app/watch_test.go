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

func TestWatchDir(t *testing.T) {
	testWatchDir := NewWatchDir("/Users/renzihao/project/github.com/gatewayorg/logbeat/test/run/containerd/io.containerd.runtime.v2.task/k8s.io")
	testWatchDir.WatchDir()
}
