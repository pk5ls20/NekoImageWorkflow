package lifecycle

import (
	"github.com/mitchellh/mapstructure"
	commonModel "github.com/pk5ls20/NekoImageWorkflow/common/model"
	apiScp "github.com/pk5ls20/NekoImageWorkflow/uploadClient/scraper/api/impl"
	localScp "github.com/pk5ls20/NekoImageWorkflow/uploadClient/scraper/local/impl"
	"github.com/pk5ls20/NekoImageWorkflow/uploadClient/scraper/model"
	clientModel "github.com/pk5ls20/NekoImageWorkflow/uploadClient/storage/config"
	"github.com/sirupsen/logrus"
)

func RegisterScraper(ScraperList clientModel.ScraperList) []model.Scraper {
	ins := make([]model.Scraper, 0)
	for scraperType, instances := range ScraperList {
		switch scraperType {
		case commonModel.LocalScraperType:
			for _, instance := range instances {
				var config clientModel.LocalScraperConfig
				if err := mapstructure.Decode(instance, &config); err != nil {
					logrus.Error("Error decoding LocalScraperConfig: ", err)
					continue
				}
				ins = append(ins, &localScp.LocalScraper{
					InsConfig: &config,
				})
			}
		case commonModel.APIScraperType:
			for _, instance := range instances {
				var config clientModel.APIScraperConfig
				if err := mapstructure.Decode(instance, &config); err != nil {
					logrus.Error("Error decoding APIScraperConfig: ", err)
					continue
				}
				ins = append(ins, &apiScp.APIScraper{
					InsConfig: &config,
				})
			}
		}
	}
	return ins
}
