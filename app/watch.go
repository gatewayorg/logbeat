package app

import (
	"os"
	"path/filepath"

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
			return err
		}
		if info.IsDir() {
			path, err := filepath.Abs(path)
			if err != nil {
				return err
			}
			err = w.watch.Add(path)
			if err != nil {
				return err
			}
			log.Info("watch", zap.Any("watch add path is", path))
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
							w.watch.Add(ev.Name)
							log.Info("watch", zap.Any("watch add path is", ev.Name))
						}
					}
					if ev.Op&fsnotify.Write == fsnotify.Write {
						log.Info("watch", zap.Any("write is", ev.Name))
						if !w.pathMap[ev.Name] {
							w.pathMap[ev.Name] = true
							w.logFilePathCh <- ev.Name
						}
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
