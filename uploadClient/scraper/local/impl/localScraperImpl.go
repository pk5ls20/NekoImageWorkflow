package impl

import (
	"github.com/pk5ls20/NekoImageWorkflow/common/log"
	commonModel "github.com/pk5ls20/NekoImageWorkflow/common/model"
	"github.com/pk5ls20/NekoImageWorkflow/uploadClient/scraper/impl"
	"github.com/pk5ls20/NekoImageWorkflow/uploadClient/scraper/local/utils"
	clientModel "github.com/pk5ls20/NekoImageWorkflow/uploadClient/storage/config"
)

type LocalScraperInstance struct {
	impl.ScraperInstance
	InsConfig *clientModel.LocalScraperConfig
}

func (c *LocalScraperInstance) PrepareData() error {
	err := utils.NewWatcher(c.InsConfig.WatchFolders)
	if err != nil {
		return log.ErrorWrap(err)
	}
	return nil
}

func (c *LocalScraperInstance) ProcessData() error {
	// TODO: make it really work
	return nil
}

func (c *LocalScraperInstance) GetType() commonModel.ScraperType {
	return commonModel.LocalScraperType
}
