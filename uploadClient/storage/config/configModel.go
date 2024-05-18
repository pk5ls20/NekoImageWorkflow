package config

import (
	commonModel "github.com/pk5ls20/NekoImageWorkflow/common/model"
)

type ConfigWrapper struct {
	ClientConfig ClientConfig `mapstructure:"clientconfig" json:"clientconfig"`
}

type ScraperList map[commonModel.ScraperType][]Scraper

type ClientConfig struct {
	ClientRegisterAddress string      `mapstructure:"ClientRegisterAddress" json:"ClientRegisterAddress"`
	ConsulAddress         string      `mapstructure:"ConsulAddress" json:"ConsulAddress"`
	PostUploadPeriod      int         `mapstructure:"PostUploadPeriod" json:"PostUploadPeriod"`
	Scraper               ScraperList `mapstructure:"scraper" json:"scraper"`
}

type Scraper interface {
}

type LocalScraperConfig struct {
	_            Scraper  `mapstructure:"LocalScraper"`
	Enable       bool     `mapstructure:"Enable"`
	WatchFolders []string `mapstructure:"WatchFolders"`
}

type APIScraperConfig struct {
	_                Scraper                  `mapstructure:"APIScraper"`
	Enable           bool                     `mapstructure:"Enable"`
	APIScraperSource []APIScraperSourceConfig `mapstructure:"APIScraperSource"`
}

type APIScraperSourceConfig struct {
	APIAddress           string            `mapstructure:"APIAddress"`
	ParserJavaScriptFile string            `mapstructure:"ParserJavaScriptFile"`
	OptionalHeaders      map[string]string `mapstructure:"OptionalHeader"`
	OptionalCookies      map[string]string `mapstructure:"OptionalCookies"`
}
