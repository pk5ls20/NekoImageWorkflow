package impl

import (
	"context"
	commonLog "github.com/pk5ls20/NekoImageWorkflow/common/log"
	commonModel "github.com/pk5ls20/NekoImageWorkflow/common/model"
	clientModel "github.com/pk5ls20/NekoImageWorkflow/uploadClient/client/model"
	apiModel "github.com/pk5ls20/NekoImageWorkflow/uploadClient/scraper/api/model"
	scraperModel "github.com/pk5ls20/NekoImageWorkflow/uploadClient/scraper/model"
	clientConfig "github.com/pk5ls20/NekoImageWorkflow/uploadClient/storage/config"
	"github.com/pk5ls20/NekoImageWorkflow/uploadClient/storage/msgQueue"
	"github.com/sirupsen/logrus"
)

// APIScraper TODO: local scraper can be designed like this
type APIScraper struct {
	scraperModel.BaseScraper
	InsConfig *clientConfig.APIScraperConfig
	ctx       context.Context
	cancel    context.CancelFunc
	fetcher   *APIFetcher
}

func (c *APIScraper) OnStart() error {
	logrus.Debugf("%s-%s Onstart Start!", c.ScraperID, c.GetType())
	if !c.Enable {
		logrus.Warn("APIScraper is not enabled")
		return nil
	}
	c.ctx, c.cancel = context.WithCancel(context.Background())
	c.fetcher = &APIFetcher{}
	if err := c.fetcher.Init(c.ScraperID); err != nil {
		return commonLog.ErrorWrap(err)
	}
	return nil
}

func (c *APIScraper) PrepareData() error {
	logrus.Debugf("%s-%s Start to fetch data from API", c.ScraperID, c.GetType())
	if !c.Enable {
		logrus.Warn("APIScraper is not enabled")
		return nil
	}
	var configs []*clientConfig.APIScraperSourceConfig
	for _, config := range c.InsConfig.APIScraperSource {
		configs = append(configs, &config)
	}
	// TODO: Using a slicing strategy for fetchList
	// TODO: add stop / recover mechanism
	todoTasks, err := c.fetcher.FetchList(configs)
	if err != nil {
		return commonLog.ErrorWrap(err)
	}
	for _, task := range todoTasks {
		if task.FetchData == nil {
			logrus.Warn("Task failed", task)
			continue
		}
		msgQueueData := &msgQueue.MsgQueueData{
			MsgMetaData: msgQueue.MsgMetaData{
				UploadType: commonModel.PreUploadType,
				MsgMetaID: msgQueue.MsgMetaID{
					ScraperID:   c.ScraperID,
					ScraperType: c.GetType(),
					MsgGroupID:  task.MsgGroupID,
				},
			},
			FileMetaData: &clientModel.AnyFileMetaDataModel{
				PreUploadFileMetaDataModel: &task.FetchData.PreUploadFileMetaDataModel,
			},
		}
		if _err := c.MsgQueue.AddElement(msgQueueData); _err != nil {
			return err
		}
	}
	return nil
}

func (c *APIScraper) ProcessData() error {
	logrus.Debugf("%s-%s Start to process data from API", c.ScraperID, c.GetType())
	// TODO: At present, every task will use an independent Fetcher, which is NOT elegant and efficient
	tasks := make([]*apiModel.SpiderToDoTask, 0, len(c.ScraperChanMap[c.ScraperID]))
	for itm := range c.ScraperChanMap[c.ScraperID] {
		todoTask := &apiModel.SpiderToDoTask{
			ScraperID: c.ScraperID,
			FetchData: itm,
		}
		tasks = append(tasks, todoTask)
		doneTasks, err := c.fetcher.FetchContent(tasks, c.ctx, c.cancel)
		if err != nil {
			return commonLog.ErrorWrap(err)
		}
		uploadModels := make([]*msgQueue.MsgQueueData, 0, len(doneTasks))
		for _, task := range doneTasks {
			if !task.Success {
				logrus.Warnf("Task failed: %v", task.FetchData)
				continue
			}
			msgQueueData := &msgQueue.MsgQueueData{
				MsgMetaData: msgQueue.MsgMetaData{
					UploadType: commonModel.PostUploadType,
					MsgMetaID: msgQueue.MsgMetaID{
						ScraperID:   c.ScraperID,
						ScraperType: c.GetType(),
						MsgGroupID:  task.MsgGroupID,
					},
				},
				FileMetaData: &clientModel.AnyFileMetaDataModel{
					UploadFileMetaDataModel: &task.FetchData.UploadFileMetaDataModel,
				},
			}
			uploadModels = append(uploadModels, msgQueueData)
		}
		if _err := c.MsgQueue.AddElements(uploadModels); _err != nil {
			return err
		}
		tasks = tasks[:0]
	}
	return nil
}

func (c *APIScraper) OnStop() error {
	logrus.Debugf("%s-%s Onstop Start!", c.ScraperID, c.GetType())
	return nil
}

func (c *APIScraper) GetType() commonModel.ScraperType {
	return commonModel.APIScraperType
}
