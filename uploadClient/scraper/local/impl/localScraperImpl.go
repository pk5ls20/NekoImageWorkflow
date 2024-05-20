package impl

import (
	commonLog "github.com/pk5ls20/NekoImageWorkflow/common/log"
	commonModel "github.com/pk5ls20/NekoImageWorkflow/common/model"
	clientModel "github.com/pk5ls20/NekoImageWorkflow/uploadClient/client/model"
	localScraperUtils "github.com/pk5ls20/NekoImageWorkflow/uploadClient/scraper/local/utils"
	scraperModel "github.com/pk5ls20/NekoImageWorkflow/uploadClient/scraper/model"
	clientConfig "github.com/pk5ls20/NekoImageWorkflow/uploadClient/storage/config"
	"github.com/sirupsen/logrus"
)

type LocalScraper struct {
	scraperModel.BaseScraper
	InsConfig *clientConfig.LocalScraperConfig
}

func (c *LocalScraper) OnStart() error {
	logrus.Debugf("%d-%s Onstart Start!", c.ScraperID, c.GetType())
	return nil
}

func (c *LocalScraper) PrepareData() error {
	logrus.Debugf("%d-%s Start to fetch data from local", c.ScraperID, c.GetType())
	err := localScraperUtils.NewWatcher(c.ScraperID, c.InsConfig.WatchFolders)
	if err != nil {
		return commonLog.ErrorWrap(err)
	}
	return nil
}

func (c *LocalScraper) ProcessData() error {
	logrus.Debugf("%d-%s Start to process data from local", c.ScraperID, c.GetType())
	// actually do nothing, just transform PreUploadFileDataModel to UploadFileDataModel
	for itm := range c.ScraperChanMap[c.ScraperID] {
		model := clientModel.NewScraperUploadFileData(itm)
		uploadModels := []*clientModel.UploadFileDataModel{model}
		if err := c.UploadQueue.Insert(uploadModels); err != nil {
			return err
		}
	}
	return nil
}

func (c *LocalScraper) OnStop() error {
	logrus.Debugf("%d-%s Onstop Start!", c.ScraperID, c.GetType())
	return nil
}

func (c *LocalScraper) GetType() commonModel.ScraperType {
	return commonModel.LocalScraperType
}
