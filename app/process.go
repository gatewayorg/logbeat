package app

import (
	"os"

	"github.com/gatewayorg/logbeat/share"
	"github.com/hpcloud/tail"
	"go.uber.org/zap"
)

type ProcessLog struct {
	path       string
	tailConfig tail.Config
	pubMetrics *PubMetrics
	filter     map[string]bool
}

func NewProcessLog(path string, pubMetrics *PubMetrics, filter []string) *ProcessLog {
	seek := &tail.SeekInfo{}
	seek.Offset = 0
	seek.Whence = os.SEEK_END
	tailConfig := tail.Config{}
	tailConfig.Follow = true
	tailConfig.Location = seek
	filterMap := make(map[string]bool)
	if len(filter) != 0 {
		for _, k := range filter {
			filterMap[k] = true
		}
	}
	return &ProcessLog{
		path:       path,
		tailConfig: tailConfig,
		pubMetrics: pubMetrics,
		filter:     filterMap,
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
		message := TransferMetricsToProtobuf(line.Text, p.filter)
		if message != nil {
			err := p.pubMetrics.ProducerPub(message)
			if err != nil {
				log.Error("Process", zap.Error(err))
			}
		}
	}
}

func StartProcess(dir string, pubMetrics *PubMetrics, filter []string) {

	WatchDir := NewWatchDir(dir)
	WatchDir.WatchDir(pubMetrics, filter)
	select {}
}
