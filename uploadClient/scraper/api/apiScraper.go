package api

import (
	commonModel "github.com/pk5ls20/NekoImageWorkflow/common/model"
	"github.com/pk5ls20/NekoImageWorkflow/uploadClient/model"
	"github.com/pk5ls20/NekoImageWorkflow/uploadClient/scraper"
)

type APIScraperInstance struct {
	scraper.ScraperInstance
	InsConfig model.APIScraperConfig
}

func (c *APIScraperInstance) PrepareData() error {
	// TODO: make it really work
	return nil
}

func (c *APIScraperInstance) ProcessData() error {
	// TODO: make it really work
	return nil
}

func (c *APIScraperInstance) GetType() commonModel.ScraperType {
	return commonModel.APIScraperType
}
