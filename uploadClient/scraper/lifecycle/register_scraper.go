package lifecycle

import (
	"fmt"
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
	localID := 0
	apiID := 0
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
				id := fmt.Sprintf("%d-local", localID)
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
				localID += 1
			}
		case commonModel.APIScraperType:
			for _, instance := range instances {
				var config configModel.APIScraperConfig
				if err := mapstructure.Decode(instance, &config); err != nil {
					logrus.Error("Error decoding APIScraperConfig: ", err)
					continue
				}
				id := fmt.Sprintf("%d-api", apiID)
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
				apiID += 1
			}
		}
	}
	return ins
}
