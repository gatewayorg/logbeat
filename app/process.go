package app

import (
	"fmt"
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
		log.Info("Process", zap.Any("log is", line.Text))
		resList := strings.Split(line.Text, " ")
		for i, v := range resList {
			fmt.Printf("log index is %d, log is %s \n", i, v)
		}
		log.Info("Process", zap.Any("log remoteAddr is", resList[0]))
		log.Info("Process", zap.Any("log remoteUser is", resList[3]))
		log.Info("Process", zap.Any("log timeLocal is", resList[4]+resList[5]))
		log.Info("Process", zap.Any("log request is", resList[7]+" "+resList[8]+" "+resList[9]))
		log.Info("Process", zap.Any("log status is", resList[10]))
		log.Info("Process", zap.Any("log bodyBytesSent is", resList[11]))
		log.Info("Process", zap.Any("log requestTime is", resList[16]))
		log.Info("Process", zap.Any("log requestBody is", resList[18]))
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
