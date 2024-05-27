package config

import (
	commonModel "github.com/pk5ls20/NekoImageWorkflow/common/model"
)

type ConfigWrapper struct {
	ClientConfig ClientConfig `mapstructure:"clientconfig" json:"clientconfig"`
}

type ScraperList map[commonModel.ScraperType][]Scraper

type ClientConfig struct {
	KitexServerAddress   string      `mapstructure:"KitexServerAddress" json:"KitexServerAddress"`
	UploadWaitSecond     float64     `mapstructure:"UploadWaitSecond" json:"UploadWaitSecond"`
	UploadFailWaitSecond float64     `mapstructure:"UploadFailWaitSecond" json:"UploadFailWaitSecond"`
	Scraper              ScraperList `mapstructure:"scraper" json:"scraper"`
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
	OptionalHeaders      map[string]string `mapstructure:"OptionalHeaders"`
	OptionalCookies      map[string]string `mapstructure:"OptionalCookies"`
}
