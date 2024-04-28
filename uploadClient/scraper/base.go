package scraper

import (
	commomModel "github.com/pk5ls20/NekoImageWorkflow/common/model"
)

type Scraper interface {
	// PrepareData prepare raw data, designed to be run in a goroutine once
	PrepareData() error
	// ProcessData make raw data to data which client can directly post, designed to be run in a goroutine once
	ProcessData() error
}

type ScraperInstance interface {
	Scraper
	GetType() commomModel.ScraperType
}
