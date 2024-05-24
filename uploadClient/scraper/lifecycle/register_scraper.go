package lifecycle

import (
	"github.com/mitchellh/mapstructure"
	commonModel "github.com/pk5ls20/NekoImageWorkflow/common/model"
	clientModel "github.com/pk5ls20/NekoImageWorkflow/uploadClient/client/model"
	apiScraper "github.com/pk5ls20/NekoImageWorkflow/uploadClient/scraper/api/impl"
	localScraper "github.com/pk5ls20/NekoImageWorkflow/uploadClient/scraper/local/impl"
	scraperModel "github.com/pk5ls20/NekoImageWorkflow/uploadClient/scraper/model"
	configModel "github.com/pk5ls20/NekoImageWorkflow/uploadClient/storage/config"
	"github.com/pk5ls20/NekoImageWorkflow/uploadClient/storage/msgQueue"
	"github.com/sirupsen/logrus"
)

func RegisterScraper(ScraperList configModel.ScraperList) []scraperModel.Scraper {
	queue := msgQueue.NewMessageQueue()
	ins := make([]scraperModel.Scraper, 0)
	id := 0
	chanMap := make(scraperModel.ScraperChanMap)
	for scraperType, instances := range ScraperList {
		switch scraperType {
		case commonModel.LocalScraperType:
			for _, instance := range instances {
				var config configModel.LocalScraperConfig
				if err := mapstructure.Decode(instance, &config); err != nil {
					logrus.Error("Error decoding LocalScraperConfig: ", err)
					continue
				}
				chanMap[id] = make(chan *clientModel.PreUploadFileDataModel)
				ins = append(ins, &localScraper.LocalScraper{
					InsConfig: &config,
					BaseScraper: scraperModel.BaseScraper{
						ScraperID:      id,
						Enable:         config.Enable,
						MsgQueue:       queue,
						ScraperChanMap: chanMap,
					},
				})
				id += 1
			}
		case commonModel.APIScraperType:
			for _, instance := range instances {
				var config configModel.APIScraperConfig
				if err := mapstructure.Decode(instance, &config); err != nil {
					logrus.Error("Error decoding APIScraperConfig: ", err)
					continue
				}
				chanMap[id] = make(chan *clientModel.PreUploadFileDataModel)
				ins = append(ins, &apiScraper.APIScraper{
					InsConfig: &config,
					BaseScraper: scraperModel.BaseScraper{
						ScraperID:      id,
						Enable:         config.Enable,
						MsgQueue:       queue,
						ScraperChanMap: chanMap,
					},
				})
				id += 1
			}
		}
	}
	return ins
}
