package model

import (
	commomModel "github.com/pk5ls20/NekoImageWorkflow/common/model"
)

type Scraper interface {
	// OnStart is called when scraper is started
	OnStart() error
	// PrepareData prepare raw data, designed to be run in a goroutine once
	PrepareData() error
	// ProcessData make raw data to data which client can directly post, designed to be run in a goroutine once
	ProcessData() error
	// OnStop is called when program stopped
	OnStop() error
	// GetType return the type of the scraper
	GetType() commomModel.ScraperType
}
