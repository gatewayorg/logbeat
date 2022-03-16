package app

import (
	"io/ioutil"
	"os"
	"time"

	"github.com/gatewayorg/logbeat/share"
	"go.uber.org/zap"
)

func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

type WatchDir struct {
	dirPath string
	pathMap map[string]bool
}

func NewWatchDir(dirPath string) *WatchDir {
	return &WatchDir{
		dirPath: dirPath,
		pathMap: make(map[string]bool),
	}
}

func (w *WatchDir) WatchDir(pubMetrics *PubMetrics) {
	for {
		log.Info("Watch", zap.Any("Start watch this dir", w.dirPath))

		files, _ := ioutil.ReadDir(w.dirPath)
		for _, f := range files {
			if f.IsDir() && len(f.Name()) == 64 {
				pathRes := w.dirPath + "/" + f.Name() + share.LOG_FILE_PATH_SUFFIX
				existsFlag, _ := PathExists(pathRes)
				if existsFlag {
					if !w.pathMap[pathRes] {
						w.pathMap[pathRes] = true
						processLog := NewProcessLog(pathRes, pubMetrics)
						processLog.Process()
						log.Info("Watch", zap.Any("start process this log", pathRes))
					} else {
						log.Info("Watch", zap.Any("processed this log", pathRes))
					}
				} else {
					log.Info("Watch", zap.Any("this log not exist", pathRes))
				}
			}
		}
		time.Sleep(time.Second * 100)
	}

}
