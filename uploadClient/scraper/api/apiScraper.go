package api

import (
	"NekoImageWorkflowKitex/common"
	"NekoImageWorkflowKitex/uploadClient/scraper"
)

type APIScraperInstance struct {
	scraper.ScraperInstance
}

func (c *APIScraperInstance) PrepareData() error {
	// TODO: make it really work
	return nil
}

func (c *APIScraperInstance) ProcessData() error {
	// TODO: make it really work
	return nil
}

func (c *APIScraperInstance) GetType() common.ScraperType {
	return common.APIScraperType
}
