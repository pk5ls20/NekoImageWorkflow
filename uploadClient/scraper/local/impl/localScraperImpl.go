package impl

import (
	"github.com/pk5ls20/NekoImageWorkflow/common/log"
	commonModel "github.com/pk5ls20/NekoImageWorkflow/common/model"
	"github.com/pk5ls20/NekoImageWorkflow/uploadClient/scraper/local/utils"
	"github.com/pk5ls20/NekoImageWorkflow/uploadClient/scraper/model"
	clientModel "github.com/pk5ls20/NekoImageWorkflow/uploadClient/storage/config"
	"github.com/sirupsen/logrus"
)

type LocalScraperInstance struct {
	model.ScraperInstance
	InsConfig *clientModel.LocalScraperConfig
}

func (c *LocalScraperInstance) OnStart() error {
	logrus.Debugf("%s Onstart Start!", c.GetType())
	return nil
}

func (c *LocalScraperInstance) PrepareData() error {
	logrus.Debugf("Start to fetch data from local")
	err := utils.NewWatcher(c.InsConfig.WatchFolders)
	if err != nil {
		return log.ErrorWrap(err)
	}
	return nil
}

func (c *LocalScraperInstance) ProcessData() error {
	logrus.Debugf("Start to process data from local")
	// TODO: make it really work
	return nil
}

func (c *LocalScraperInstance) OnStop() error {
	logrus.Debugf("%s Onstop Start!", c.GetType())
	return nil
}

func (c *LocalScraperInstance) GetType() commonModel.ScraperType {
	return commonModel.LocalScraperType
}