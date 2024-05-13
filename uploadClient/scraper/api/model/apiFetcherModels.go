package model

import "github.com/pk5ls20/NekoImageWorkflow/uploadClient/storage/config"

type FetcherTask struct {
	APIAddress    string
	PasteFilePath string
}

// Fetcher is the interface for fetching data from API, include parser
type Fetcher interface {
	// Init initialize the fetcher
	Init(cf *[]config.APIScraperSourceConfig) error
	FetchList() (url []string, err error)
	FetchContent(url []string) ([]*SpiderTask, error)
}
