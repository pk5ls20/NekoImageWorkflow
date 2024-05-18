package impl

import (
	"context"
	commonLog "github.com/pk5ls20/NekoImageWorkflow/common/log"
	commonModel "github.com/pk5ls20/NekoImageWorkflow/common/model"
	clientModel "github.com/pk5ls20/NekoImageWorkflow/uploadClient/client/model"
	scraperImpl "github.com/pk5ls20/NekoImageWorkflow/uploadClient/scraper/model"
	clientConfig "github.com/pk5ls20/NekoImageWorkflow/uploadClient/storage/config"
	storageQueue "github.com/pk5ls20/NekoImageWorkflow/uploadClient/storage/queue"
	"github.com/sirupsen/logrus"
)

// APIScraper TODO: local scraper can be designed like this
type APIScraper struct {
	scraperImpl.Scraper
	ScraperID      int
	InsConfig      *clientConfig.APIScraperConfig
	enable         bool
	ctx            context.Context
	cancel         context.CancelFunc
	fetcher        *APIFetcher
	preUploadQueue *storageQueue.PreUploadQueue
	uploadQueue    *storageQueue.UploadQueue
}

func (c *APIScraper) OnStart() error {
	logrus.Debugf("%s Onstart Start!", c.GetType())
	c.enable = c.InsConfig.Enable
	if !c.enable {
		logrus.Warn("APIScraper is not enabled")
		return nil
	}
	c.ctx, c.cancel = context.WithCancel(context.Background())
	c.fetcher = &APIFetcher{}
	c.preUploadQueue = storageQueue.GetPreUploadQueue()
	c.uploadQueue = storageQueue.GetUploadQueue()
	if err := c.fetcher.Init(c.ScraperID); err != nil {
		return commonLog.ErrorWrap(err)
	}
	return nil
}

func (c *APIScraper) PrepareData() error {
	logrus.Debugf("Start to fetch data from API")
	if !c.enable {
		logrus.Warn("APIScraper is not enabled")
		return nil
	}
	var configs []*clientConfig.APIScraperSourceConfig
	for _, config := range c.InsConfig.APIScraperSource {
		configs = append(configs, &config)
	}
	// TODO: Using a slicing strategy for fetchList
	// TODO: add stop / recover mechanism
	todoTasks, err := c.fetcher.FetchList(configs)
	if err != nil {
		return commonLog.ErrorWrap(err)
	}
	for _, task := range todoTasks {
		if task.FetchData == nil {
			logrus.Warn("Task failed", task)
			continue
		}
		if _err := c.preUploadQueue.Insert([]*clientModel.PreUploadFileDataModel{task.FetchData}); _err != nil {
			return _err
		}
	}
	return nil
}

// ProcessData TODO: make it really work
func (c *APIScraper) ProcessData() error {
	logrus.Debugf("Start to process data from API")
	// TODO: Using a slicing strategy for todoTasks
	// TODO: use channel to sync???
	//doneTasks, err := c.fetcher.FetchContent(c.todoTasks, c.ctx, c.cancel)
	//if err != nil {
	//	return log.ErrorWrap(err)
	//}
	//for _, task := range doneTasks {
	//	if !task.Success {
	//		logrus.Warn("Task failed", task.FetchData)
	//		continue
	//	}
	//}
	//// TODO: do we need retry more failed data?
	//preUploadModels := make([]*clientModel.UploadFileDataModel, 0)
	//for _, task := range doneTasks {
	//	preUploadModels = append(preUploadModels, task.FetchData)
	//}
	return nil
}

func (c *APIScraper) OnStop() error {
	logrus.Debugf("%s Onstop Start!", c.GetType())
	return nil
}

func (c *APIScraper) GetType() commonModel.ScraperType {
	return commonModel.APIScraperType
}

func (c *APIScraper) GetID() int {
	return c.ScraperID
}
