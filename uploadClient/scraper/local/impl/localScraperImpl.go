package impl

import (
	commonLog "github.com/pk5ls20/NekoImageWorkflow/common/log"
	commonModel "github.com/pk5ls20/NekoImageWorkflow/common/model"
	clientModel "github.com/pk5ls20/NekoImageWorkflow/uploadClient/client/model"
	localScraperUtils "github.com/pk5ls20/NekoImageWorkflow/uploadClient/scraper/local/utils"
	scraperModel "github.com/pk5ls20/NekoImageWorkflow/uploadClient/scraper/model"
	clientConfig "github.com/pk5ls20/NekoImageWorkflow/uploadClient/storage/config"
	"github.com/pk5ls20/NekoImageWorkflow/uploadClient/storage/msgQueue"
	"github.com/sirupsen/logrus"
)

type LocalScraper struct {
	scraperModel.BaseScraper
	InsConfig *clientConfig.LocalScraperConfig
}

func (c *LocalScraper) OnStart() error {
	logrus.Debugf("%d-%s Onstart Start!", c.ScraperID, c.GetType())
	return nil
}

func (c *LocalScraper) PrepareData() error {
	logrus.Debugf("%d-%s Start to fetch data from local", c.ScraperID, c.GetType())
	err := localScraperUtils.NewWatcher(c.ScraperID, c.InsConfig.WatchFolders)
	if err != nil {
		return commonLog.ErrorWrap(err)
	}
	return nil
}

func (c *LocalScraper) ProcessData() error {
	logrus.Debugf("%d-%s Start to process data from local", c.ScraperID, c.GetType())
	// actually do nothing, just transform PreUploadFileDataModel to UploadFileDataModel
	queue := msgQueue.NewMessageQueue()
	for itm := range c.ScraperChanMap[c.ScraperID] {
		oriData := clientModel.NewUploadFileData(itm)
		model := msgQueue.MsgQueueData{
			MsgMetaData: msgQueue.MsgMetaData{
				UploadType: commonModel.PostUploadType,
				MsgMetaID: msgQueue.MsgMetaID{
					ScraperType: commonModel.LocalScraperType,
					ScraperID:   c.ScraperID,
					MsgGroupID:  0, //TODO:
				},
			},
			FileMetaData: &clientModel.AnyFileMetaDataModel{
				UploadFileMetaDataModel: &oriData.UploadFileMetaDataModel,
			},
		}
		if err := queue.AddElement(&model); err != nil {
			return err
		}
	}
	return nil
}

func (c *LocalScraper) OnStop() error {
	logrus.Debugf("%d-%s Onstop Start!", c.ScraperID, c.GetType())
	return nil
}

func (c *LocalScraper) GetType() commonModel.ScraperType {
	return commonModel.LocalScraperType
}
