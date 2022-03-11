package app

import (
	"fmt"
	"os"

	"github.com/hpcloud/tail"
	"go.uber.org/zap"
)

type ProcessLog struct {
	path       string
	tailConfig tail.Config
}

func NewProcessLog(path string) *ProcessLog {
	seek := &tail.SeekInfo{}
	seek.Offset = 0
	seek.Whence = os.SEEK_END
	tailConfig := tail.Config{}
	tailConfig.Follow = true
	tailConfig.Location = seek
	return &ProcessLog{
		path:       path,
		tailConfig: tailConfig,
	}
}

func (p *ProcessLog) Process() {
	fmt.Println("Process ", p.path[len(p.path)-28:])
	if p.path[len(p.path)-28:] == "/rootfs/root/logs/access.log" {
		go func(filePath string) {
			p.TailLog(filePath)
		}(p.path)
	}
}

func (p *ProcessLog) TailLog(filePath string) {
	t, err := tail.TailFile(filePath, p.tailConfig)
	if err != nil {
		log.Error("tail file fail", zap.Error(err))
		return
	}
	for line := range t.Lines {
		log.Info("Process", zap.Any("log is", line.Text))
	}
}

func StartProcess(dir string) {
	watchLog := NewWatchLog(dir)
	defer watchLog.watch.Close()
	// testProcessLog := NewProcessLog(testWatchLog.logFilePathCh)
	watchLog.WatchDir()
	// for path := range watchLog.logFilePathCh {
	// 	ProcessLog := NewProcessLog(path)
	// 	ProcessLog.Process()
	// }
	select {}
}
