package app

import (
	"os"
	"strings"

	"github.com/gatewayorg/logbeat/share"
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
	if p.path[len(p.path)-28:] == share.LOG_FILE_PATH_SUFFIX {
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
		resList := strings.Split(line.Text, " ")
		for _, v := range resList {
			log.Info("Process", zap.Any("log is", v))
		}
	}
}

func StartProcess(dir string) {
	// watchLog := NewWatchLog(dir)
	// defer watchLog.watch.Close()
	// watchLog.WatchDir()
	// go watchLog.CheckNotExist()
	// select {}

	WatchDir := NewWatchDir(dir)
	WatchDir.WatchDir()
	select {}
}
