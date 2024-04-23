package local

import (
	"NekoImageWorkflowKitex/common"
	"NekoImageWorkflowKitex/uploadClient/model"
	"NekoImageWorkflowKitex/uploadClient/scraper"
	"NekoImageWorkflowKitex/uploadClient/storage"
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/sirupsen/logrus"
	"io/fs"
	"path/filepath"
)

type LocalScraperInstance struct {
	scraper.ScraperInstance
}

func walkDir(rootPath string) ([]string, error) {
	var dirs []string
	dirs = append(dirs, rootPath)
	err := filepath.WalkDir(rootPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			dirs = append(dirs, path)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return dirs, nil
}

func (c *LocalScraperInstance) PrepareData() error {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		logrus.Fatal(err)
	}
	defer func(watcher *fsnotify.Watcher) {
		err := watcher.Close()
		if err != nil {
			logrus.Fatal(err)
		}
	}(watcher)
	go func() {
		preUploadBridgeIns := storage.GetPreUploadTransBridgeInstance()
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				logrus.Debug("event:", event)
				if event.Has(fsnotify.Create) {
					logrus.Debug("create file:", event.Name)
					uuid, _ := common.GenerateFileUUID(event.Name)
					tmpSlice := []model.PreUploadFileData{{ResourceUUID: uuid, ResourceUri: event.Name}}
					err := preUploadBridgeIns.Insert(tmpSlice)
					if err != nil {
						logrus.Error("Failed to insert file into preUploadBridgeIns:", err)
					}
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				logrus.Fatal("error:", err)
			}
		}
	}()
	config := storage.GetConfig()
	allRecursiveDirs := make([]string, 0)
	for _, dir := range config.ScraperConfig.LocalScraperConfig.WatchFolders {
		recursiveDirs, _ := walkDir(dir)
		allRecursiveDirs = append(allRecursiveDirs, recursiveDirs...)
	}
	for _, dir := range allRecursiveDirs {
		err = watcher.Add(dir)
		if err != nil {
			logrus.Fatal(err)
		}
	}
	if err != nil {
		logrus.Fatal(err)
	}
	logrus.Info(
		fmt.Sprintf("Successfully started watching %d folders under %s",
			len(allRecursiveDirs), config.ScraperConfig.LocalScraperConfig.WatchFolders),
	)
	// Block main goroutine forever.
	<-make(chan struct{})
	return nil
}

func (c *LocalScraperInstance) ProcessData() error {
	// TODO: make it really work
	return nil
}
