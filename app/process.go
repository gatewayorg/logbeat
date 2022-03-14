package app

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/gatewayorg/logbeat/share"
	"github.com/hpcloud/tail"
	"go.uber.org/zap"
)

var (
	rgx         = regexp.MustCompile(`\{(.*?)\}`)
	nginxFormat = `$remote_addr - $remote_user  $time_local  "$request" '
	'$status $body_bytes_sent "$http_referer" '
	'"$http_user_agent" "$http_x_forwarded_for"'
	'$upstream_addr  $request_time "$upstream_response_time"'
			  '"$upstream_cache_status" "$upstream_addr"' '$request_body`
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
		if len(resList) >= 15 {
			for i, v := range resList {
				fmt.Printf("log index is %d, log is %s \n", i, v)
			}
			log.Info("Process", zap.Any("log remoteAddr is", resList[0]))
			log.Info("Process", zap.Any("log remoteUser is", resList[3]))
			log.Info("Process", zap.Any("log timeLocal is", resList[4]+resList[5]))
			log.Info("Process", zap.Any("log request is", resList[7]+" "+resList[8]+" "+resList[9]))
			log.Info("Process", zap.Any("log status is", resList[10]))
		}
		responseList := rgx.FindStringSubmatch(line.Text)
		if len(responseList) == 1 {
			log.Info("Process", zap.Any("log request body is", responseList[0]))
		}

		// reader := gonx.NewReader(strings.NewReader(line.Text), nginxFormat)
		// res, err := reader.Read()
		// if err != nil{
		// 	log.Info()
		// }
		// log.Info("Process", zap.Any("res", res))

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
