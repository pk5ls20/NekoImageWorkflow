package impl

import (
	commonModel "github.com/pk5ls20/NekoImageWorkflow/common/model"
	scraperImpl "github.com/pk5ls20/NekoImageWorkflow/uploadClient/scraper/model"
	clientModel "github.com/pk5ls20/NekoImageWorkflow/uploadClient/storage/config"
	"github.com/sirupsen/logrus"
)

type APIScraperWorker struct {
	spider *APISpiderImpl
}

type APIScraperInstance struct {
	scraperImpl.ScraperInstance
	InsConfig      *clientModel.APIScraperConfig
	fetchAPISpider []*APIScraperWorker
}

func (c *APIScraperInstance) OnStart() error {
	logrus.Debugf("%s Onstart Start!", c.GetType())
	return nil
}

// PrepareData TODO: implement this
func (c *APIScraperInstance) PrepareData() error {
	logrus.Debugf("Start to fetch data from API")
	return nil
}

// ProcessData TODO: make it really work
func (c *APIScraperInstance) ProcessData() error {
	logrus.Debugf("Start to process data from API")
	return nil
}

func (c *APIScraperInstance) OnStop() error {
	logrus.Debugf("%s Onstop Start!", c.GetType())
	return nil
}

func (c *APIScraperInstance) GetType() commonModel.ScraperType {
	return commonModel.APIScraperType
}
