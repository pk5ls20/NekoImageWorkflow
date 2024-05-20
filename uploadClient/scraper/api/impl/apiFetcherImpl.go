package impl

import (
	"context"
	"errors"
	commonLog "github.com/pk5ls20/NekoImageWorkflow/common/log"
	commonModel "github.com/pk5ls20/NekoImageWorkflow/common/model"
	commonUtils "github.com/pk5ls20/NekoImageWorkflow/common/utils"
	clientModel "github.com/pk5ls20/NekoImageWorkflow/uploadClient/client/model"
	apiModel "github.com/pk5ls20/NekoImageWorkflow/uploadClient/scraper/api/model"
	"github.com/pk5ls20/NekoImageWorkflow/uploadClient/storage/config"
	"github.com/sirupsen/logrus"
	"time"
)

// APIFetcher is the implementation of Fetcher, init every time you use it
type APIFetcher struct {
	apiModel.Fetcher
	scraperID           int
	fetchListParser     APIParser
	fetchListHttpClient *commonUtils.HttpClient
	initialized         bool
}

func (a *APIFetcher) Init(scID int) error {
	a.scraperID = scID
	logrus.Debugf("%d-APIFetcher Init start", a.scraperID)
	a.fetchListParser.Init()
	a.fetchListHttpClient = commonUtils.NewHttpClient()
	a.initialized = true
	logrus.Debugf("%d-APIFetcher initialized", a.scraperID)
	return nil
}

func (a *APIFetcher) FetchList(cf []*config.APIScraperSourceConfig) ([]*apiModel.SpiderToDoTask, error) {
	if !a.initialized {
		return nil, commonLog.ErrorWrap(errors.New("APIFetcher not initialized"))
	}
	logrus.Debugf("%d-APIFetcher Fetch start", a.scraperID)
	var infoList []*apiModel.SpiderToDoTask
	var fetchTaskList []*apiModel.FetcherTaskList
	for _, c := range cf {
		if err := a.fetchListParser.RegisterParser(c.ParserJavaScriptFile); err != nil {
			return infoList, commonLog.ErrorWrap(err)
		}
		fetchTaskList = append(fetchTaskList, &apiModel.FetcherTaskList{
			APIAddress:    c.APIAddress,
			PasteFilePath: c.ParserJavaScriptFile,
			Headers:       c.OptionalHeaders,
			Cookies:       c.OptionalCookies,
		})
	}
	// TODO: add retry mechanism
	for _, task := range fetchTaskList {
		response, err := a.fetchListHttpClient.Get(task.APIAddress, task.Headers, task.Cookies)
		if err != nil {
			return nil, commonLog.ErrorWrap(err)
		}
		rawJson := string(response)
		parsedUrls, err := a.fetchListParser.ParseJson(rawJson, task.PasteFilePath)
		if err != nil {
			return nil, commonLog.ErrorWrap(err)
		}
		for _, url := range parsedUrls {
			data, _err := clientModel.NewScraperPreUploadFileData(commonModel.APIScraperType, a.scraperID, url)
			if _err != nil {
				return nil, commonLog.ErrorWrap(_err)
			}
			newTask := &apiModel.SpiderToDoTask{
				SpiderTask: &apiModel.SpiderTask{
					Url:     url,
					Headers: task.Headers,
					Cookies: task.Cookies,
				},
				ScraperID: a.scraperID,
				FetchData: data,
			}
			infoList = append(infoList, newTask)
		}
	}
	return infoList, nil
}

func (a *APIFetcher) FetchContent(task []*apiModel.SpiderToDoTask,
	ctx context.Context, cancel context.CancelFunc) ([]*apiModel.SpiderDoTask, error) {
	if !a.initialized {
		return nil, commonLog.ErrorWrap(errors.New("APIFetcher not initialized"))
	}
	logrus.Debugf("%d-APIFetcher FetchContent start", a.scraperID)
	spider := &APISpider{}
	spiderConfig := &apiModel.SpiderConfig{
		SingleTaskMaxRetriesTime:    5,
		SingleTaskRetryDuration:     1000 * time.Millisecond,
		ConcurrentTaskLimit:         10,
		ConcurrentTaskGroupDuration: 300 * time.Millisecond,
		AdjustLimitRate:             0.3,
		AdjustLimitCheckTime:        500 * time.Millisecond,
	}
	if err := spider.Init(task, spiderConfig, ctx, cancel); err != nil {
		return nil, commonLog.ErrorWrap(err)
	}
	if err := spider.Start(); err != nil {
		return nil, commonLog.ErrorWrap(err)
	}
	rs, err := spider.WaitDone()
	if err != nil {
		return nil, commonLog.ErrorWrap(err)
	}
	return rs, nil
}
