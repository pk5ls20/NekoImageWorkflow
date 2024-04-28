package local

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/pk5ls20/NekoImageWorkflow/common/log"
	commonModel "github.com/pk5ls20/NekoImageWorkflow/common/model"
	"github.com/pk5ls20/NekoImageWorkflow/common/uuid"
	"github.com/pk5ls20/NekoImageWorkflow/uploadClient/model"
	"github.com/pk5ls20/NekoImageWorkflow/uploadClient/scraper"
	"github.com/pk5ls20/NekoImageWorkflow/uploadClient/storage/bridge"
	"github.com/sirupsen/logrus"
	"io/fs"
	"path/filepath"
)

type LocalScraperInstance struct {
	scraper.ScraperInstance
	InsConfig model.LocalScraperConfig
}

func walkDir(rootPath string) (*[]string, error) {
	dirs := make([]string, 0)
	dirs = append(dirs, rootPath)
	walkDirErr := filepath.WalkDir(rootPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return log.ErrorWrap(err)
		}
		if d.IsDir() {
			dirs = append(dirs, path)
		}
		return nil
	})
	if walkDirErr != nil {
		return nil, log.ErrorWrap(walkDirErr)
	}
	return &dirs, nil
}

func (c *LocalScraperInstance) PrepareData() error {
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
		preUploadBridgeIns := bridge.GetPreUploadTransBridgeInstance()
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				logrus.Debug("event:", event)
				if event.Has(fsnotify.Create) {
					logrus.Debug("create file:", event.Name)
					uuid, _ := uuid.GenerateFileUUID(event.Name)
					tmpSlice := []*model.ScraperPreUploadFileDataModel{{ResourceUUID: uuid, ResourceUri: event.Name}}
					if err_ := preUploadBridgeIns.Insert(tmpSlice); err_ != nil {
						logrus.Error("Failed to insert file into preUploadBridgeIns:", err_)
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
	for _, dir := range c.InsConfig.WatchFolders {
		recursiveDirs, err := walkDir(dir)
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
			len(allRecursiveDirs), c.InsConfig.WatchFolders),
	)
	// Block main goroutine forever.
	<-make(chan struct{})
	return nil
}

func (c *LocalScraperInstance) ProcessData() error {
	// TODO: make it really work
	return nil
}

func (c *LocalScraperInstance) GetType() commonModel.ScraperType {
	return commonModel.LocalScraperType
}
