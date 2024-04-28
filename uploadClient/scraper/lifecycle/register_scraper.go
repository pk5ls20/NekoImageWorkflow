package lifecycle

import (
	"github.com/mitchellh/mapstructure"
	commonModel "github.com/pk5ls20/NekoImageWorkflow/common/model"
	"github.com/pk5ls20/NekoImageWorkflow/uploadClient/model"
	"github.com/pk5ls20/NekoImageWorkflow/uploadClient/scraper"
	apiScp "github.com/pk5ls20/NekoImageWorkflow/uploadClient/scraper/api"
	localScp "github.com/pk5ls20/NekoImageWorkflow/uploadClient/scraper/local"
	"github.com/sirupsen/logrus"
)

func RegisterScraper(ScraperInstanceList model.ScraperInstanceList) []scraper.ScraperInstance {
	ins := make([]scraper.ScraperInstance, 0)
	for scraperType, instances := range ScraperInstanceList {
		switch scraperType {
		case commonModel.LocalScraperType:
			for _, instance := range instances {
				var config model.LocalScraperConfig
				if err := mapstructure.Decode(instance, &config); err != nil {
					logrus.Error("Error decoding LocalScraperConfig: ", err)
					continue
				}
				ins = append(ins, &localScp.LocalScraperInstance{
					InsConfig: config,
				})
			}
		case commonModel.APIScraperType:
			for _, instance := range instances {
				var config model.APIScraperConfig
				if err := mapstructure.Decode(instance, &config); err != nil {
					logrus.Error("Error decoding APIScraperConfig: ", err)
					continue
				}
				ins = append(ins, &apiScp.APIScraperInstance{
					InsConfig: config,
				})
			}
		}
	}
	return ins
}
