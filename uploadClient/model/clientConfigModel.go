package model

import (
	"github.com/pk5ls20/NekoImageWorkflow/common/model"
)

type ConfigWrapper struct {
	ClientConfig ClientConfig `mapstructure:"clientconfig" json:"clientconfig"`
}

type ScraperInstanceList map[model.ScraperType][]ScraperInstance

type ClientConfig struct {
	ClientID              string              `mapstructure:"ClientID" json:"ClientID"`
	ClientName            string              `mapstructure:"ClientName" json:"ClientName"`
	DestServiceName       string              `mapstructure:"DestServiceName" json:"DestServiceName"`
	ClientRegisterAddress string              `mapstructure:"ClientRegisterAddress" json:"ClientRegisterAddress"`
	ConsulAddress         string              `mapstructure:"ConsulAddress" json:"ConsulAddress"`
	PostUploadPeriod      int                 `mapstructure:"PostUploadPeriod" json:"PostUploadPeriod"`
	ScraperInstance       ScraperInstanceList `mapstructure:"ScraperInstance" json:"ScraperInstance"`
}

type ScraperInstance interface {
}

type LocalScraperConfig struct {
	_            ScraperInstance `mapstructure:"LocalScraper"`
	Enable       bool            `mapstructure:"Enable"`
	WatchFolders []string        `mapstructure:"WatchFolders"`
}

type APIScraperConfig struct {
	_                ScraperInstance          `mapstructure:"APIScraper"`
	Enable           bool                     `mapstructure:"Enable"`
	APIScraperSource []APIScraperSourceConfig `mapstructure:"APIScraperSource"`
}

type APIScraperSourceConfig struct {
	APIAddress           string `mapstructure:"APIAddress"`
	ParserJavaScriptFile string `mapstructure:"ParserJavaScriptFile"`
}
