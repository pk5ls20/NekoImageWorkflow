package model

import (
	"errors"
	"github.com/pk5ls20/NekoImageWorkflow/common/log"
)

type ScraperType string

const (
	LocalScraperType ScraperType = "localscraper"
	APIScraperType   ScraperType = "apiscraper"
)

func ParseScraperType(s string) (ScraperType, error) {
	switch s {
	case string(LocalScraperType):
		return LocalScraperType, nil
	case string(APIScraperType):
		return APIScraperType, nil
	default:
		return "", log.ErrorWrap(errors.New("invalid ScraperType"))
	}
}
