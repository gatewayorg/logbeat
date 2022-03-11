package app

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/fsnotify/fsnotify"
	"go.uber.org/zap"
)

type WatchLog struct {
	dirPath       string
	watch         *fsnotify.Watcher
	logFilePathCh chan string
	pathMap       map[string]bool
}

func NewWatchLog(dirPath string) *WatchLog {
	watch, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal("NewWatcher err", zap.Error(err))
	}
	return &WatchLog{
		dirPath:       dirPath,
		watch:         watch,
		logFilePathCh: make(chan string),
		pathMap:       make(map[string]bool),
	}
}

func (w *WatchLog) WatchDir() {
	filepath.Walk(w.dirPath, func(path string, info os.FileInfo, err error) error {
		if info == nil {
			return nil
		}
		if info.IsDir() {
			path, err := filepath.Abs(path)
			if err != nil {
				return err
			}
			if w.checkPath(path) {
				err = w.watch.Add(path)
				if err != nil {
					return err
				}
				log.Info("watch", zap.Any("init watch add path is", path))
			}

		}
		return nil
	})
	go func() {
		for {
			select {
			case ev := <-w.watch.Events:
				{
					if ev.Op&fsnotify.Create == fsnotify.Create {
						log.Info("watch", zap.Any("create is", ev.Name))
						fi, err := os.Stat(ev.Name)
						if err == nil && fi.IsDir() {
							if w.checkPath(ev.Name) {
								w.watch.Add(ev.Name)
								log.Info("watch", zap.Any("watch add path is", ev.Name))
							}
						}
					}
					if ev.Op&fsnotify.Write == fsnotify.Write {
						log.Info("watch", zap.Any("write is", ev.Name))
						// if !w.pathMap[ev.Name] {
						// 	w.pathMap[ev.Name] = true
						// 	w.logFilePathCh <- ev.Name
						// }
					}
					if ev.Op&fsnotify.Remove == fsnotify.Remove {
						log.Info("watch", zap.Any("remove is", ev.Name))
						fi, err := os.Stat(ev.Name)
						if err == nil && fi.IsDir() {
							w.watch.Remove(ev.Name)
							log.Info("watch", zap.Any("watch remove is", ev.Name))
						}
					}
					if ev.Op&fsnotify.Rename == fsnotify.Rename {
						log.Info("watch", zap.Any("rename is", ev.Name))
						w.watch.Remove(ev.Name)
					}
				}
			case err := <-w.watch.Errors:
				{
					log.Error("watch err is", zap.Error(err))
					return
				}
			}
		}
	}()
}

func (w *WatchLog) checkPath(path string) bool {

	// fmt.Println(path)
	// fmt.Println(w.dirPath)

	if path == w.dirPath {
		return true
	}

	pathList := strings.Split(path, "/")

	if len(pathList[len(pathList)-1]) == 64 {
		pathRes := w.dirPath + "/" + pathList[len(pathList)-1] + "/rootfs/root/logs/access.log"
		fmt.Println(pathRes)
		existsFlag, _ := PathExists(pathRes)
		fmt.Println("existsFlag", existsFlag)
		if existsFlag {
			if !w.pathMap[pathRes] {
				w.pathMap[pathRes] = true
				processLog := NewProcessLog(pathRes)
				processLog.Process()
				// w.logFilePathCh <- pathRes
			}
		}
		return true
	}

	return false
}

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
