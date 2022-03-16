package app

import (
	"fmt"
	"math/rand"
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

	StartProcess("../test/", nil)
}

func TestTailLog(t *testing.T) {
	testProcessLog := NewProcessLog("../test/access.log", nil)
	testProcessLog.Process()
	select {}
}

func TestWatchDir(t *testing.T) {
	testWatchDir := NewWatchDir("/Users/renzihao/project/github.com/gatewayorg/logbeat/test/run/containerd/io.containerd.runtime.v2.task/k8s.io")
	testWatchDir.WatchDir(nil)
}

func TestTransferMetricsToProtobuf(t *testing.T) {
	logText := "220.197.189.56 - -  14/Mar/2022:13:53:08 +0000  \"POST /d2852ab01af54c51be3a6575edfe3a97/04f8bc9c328003b6cfd4ee8ab435fe58/binance/full/main HTTP/1.1\" 200 44 \"-\" \"curl/7.77.0\" \"-\"10.98.126.191:8545  0.004 \"0.000\"\"-\" \"10.98.126.191:8545\"{\\x22jsonrpc\\x22:\\x222.0\\x22,\\x22method\\x22:\\x22eth_blockNumber\\x22,\\x22params\\x22:[],\\x22id\\x22:1}"
	t.Log(logText)
	TransferMetricsToProtobuf(logText)
}

func TestRandomGetAccess(t *testing.T) {

	for i := 0; i < 100; i++ {
		r := rand.Intn(3)
		fmt.Println(r)

	}

}
