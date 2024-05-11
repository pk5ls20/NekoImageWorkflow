package utils

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/pk5ls20/NekoImageWorkflow/common/log"
	commonModel "github.com/pk5ls20/NekoImageWorkflow/common/model"
	clientModel "github.com/pk5ls20/NekoImageWorkflow/uploadClient/client/model"
	"github.com/pk5ls20/NekoImageWorkflow/uploadClient/storage/queue"
	"github.com/sirupsen/logrus"
)

func NewWatcher(watchFolders []string) error {
	watcher, watcherErr := fsnotify.NewWatcher()
	if watcherErr != nil {
		return log.ErrorWrap(watcherErr)
	}
	defer func(watcher *fsnotify.Watcher) {
		if err_ := watcher.Close(); err_ != nil {
			logrus.Error(err_)
		}
	}(watcher)
	go func() {
		preUploadQueue := queue.GetPreUploadQueueInstance()
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				logrus.Debug("event:", event)
				if event.Has(fsnotify.Create) {
					logrus.Debug("create file:", event.Name)
					d, err := clientModel.NewScraperPreUploadFileData(commonModel.LocalScraperType, event.Name)
					if err != nil {
						logrus.Errorf("Failed to create ScraperPreUploadFileData: %v", err)
						continue
					}
					if _err := preUploadQueue.Insert([]*clientModel.PreUploadFileDataModel{d}); _err != nil {
						logrus.Errorf("Failed to insert file into preUploadQueue: %v", err)
					}
				}
			case err_, ok := <-watcher.Errors:
				if !ok {
					logrus.Error("watcher.Errors not ok")
					return
				}
				logrus.Error("error:", err_)
			}
		}
	}()
	allRecursiveDirs := make([]string, 0)
	for _, dir := range watchFolders {
		recursiveDirs, err := WalkDir(dir)
		if err != nil {
			logrus.Error(err)
			continue
		}
		allRecursiveDirs = append(allRecursiveDirs, *recursiveDirs...)
	}
	for _, dir := range allRecursiveDirs {
		if err := watcher.Add(dir); err != nil {
			logrus.Error(err)
		}
	}
	logrus.Info(
		fmt.Sprintf("Successfully started watching %d folders under %s",
			len(allRecursiveDirs), watchFolders),
	)
	// Block main goroutine forever.
	<-make(chan struct{})
	return nil
}
