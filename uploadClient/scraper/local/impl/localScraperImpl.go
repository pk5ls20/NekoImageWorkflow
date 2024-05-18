package impl

import (
	commonLog "github.com/pk5ls20/NekoImageWorkflow/common/log"
	commonModel "github.com/pk5ls20/NekoImageWorkflow/common/model"
	localScraperUtils "github.com/pk5ls20/NekoImageWorkflow/uploadClient/scraper/local/utils"
	scraperModel "github.com/pk5ls20/NekoImageWorkflow/uploadClient/scraper/model"
	clientModel "github.com/pk5ls20/NekoImageWorkflow/uploadClient/storage/config"
	"github.com/sirupsen/logrus"
)

type LocalScraper struct {
	scraperModel.Scraper
	ScraperID int
	InsConfig *clientModel.LocalScraperConfig
}

func (c *LocalScraper) OnStart() error {
	logrus.Debugf("%s Onstart Start!", c.GetType())
	return nil
}

func (c *LocalScraper) PrepareData() error {
	logrus.Debugf("Start to fetch data from local")
	err := localScraperUtils.NewWatcher(c.ScraperID, c.InsConfig.WatchFolders)
	if err != nil {
		return commonLog.ErrorWrap(err)
	}
	return nil
}

func (c *LocalScraper) ProcessData() error {
	logrus.Debugf("Start to process data from local")
	// TODO: make it really work
	return nil
}

func (c *LocalScraper) OnStop() error {
	logrus.Debugf("%s Onstop Start!", c.GetType())
	return nil
}

func (c *LocalScraper) GetType() commonModel.ScraperType {
	return commonModel.LocalScraperType
}

func (c *LocalScraper) GetID() int {
	return c.ScraperID
}
