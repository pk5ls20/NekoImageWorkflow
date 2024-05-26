package model

import (
	"context"
	"github.com/pk5ls20/NekoImageWorkflow/uploadClient/storage/config"
)

type FetcherTaskList struct {
	APIAddress    string
	PasteFilePath string
	Headers       map[string]string
	Cookies       map[string]string
}

// Fetcher is the interface for fetching data from API, include parser
type Fetcher interface {
	// Init initialize the fetcher
	Init(scID string) error
	FetchList(cf []*config.APIScraperSourceConfig) (task []*SpiderToDoTask, err error)
	FetchContent(task []*SpiderToDoTask, ctx context.Context, cancel context.CancelFunc) ([]*SpiderDoTask, error)
}
