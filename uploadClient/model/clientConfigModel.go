package model

import "NekoImageWorkflowKitex/common"

type ConfigWrapper struct {
	ClientConfig ClientConfig `mapstructure:"clientconfig" json:"clientconfig"`
}

type ClientConfig struct {
	ClientID              string               `mapstructure:"ClientID" json:"ClientID"`
	ClientName            string               `mapstructure:"ClientName" json:"ClientName"`
	DestServiceName       string               `mapstructure:"DestServiceName" json:"DestServiceName"`
	ClientRegisterAddress string               `mapstructure:"ClientRegisterAddress" json:"ClientRegisterAddress"`
	ConsulAddress         string               `mapstructure:"ConsulAddress" json:"ConsulAddress"`
	PostUploadPeriod      int                  `mapstructure:"PostUploadPeriod" json:"PostUploadPeriod"`
	ScraperList           []common.ScraperType `mapstructure:"ScraperList" json:"ScraperList"`
	ScraperConfig         ScraperConfig        `mapstructure:"ScraperConfig" json:"ScraperConfig"`
}

type ScraperConfig struct {
	LocalScraperConfig LocalScraperConfig `mapstructure:"LocalScraperConfig"`
	APIScraperConfig   APIScraperConfig   `mapstructure:"APIScraperConfig"`
}

type LocalScraperConfig struct {
	WatchFolders []string `mapstructure:"WatchFolders"`
}

type APIScraperConfig struct {
	APIScraperSource []APIScraperSourceConfig `mapstructure:"APIScraperSource"`
}

type APIScraperSourceConfig struct {
	APIAddress           string `mapstructure:"APIAddress"`
	ParserJavaScriptFile string `mapstructure:"ParserJavaScriptFile"`
}
