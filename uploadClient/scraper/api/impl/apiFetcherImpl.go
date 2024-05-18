package impl

import (
	"context"
	"errors"
	"github.com/pk5ls20/NekoImageWorkflow/common/log"
	"github.com/pk5ls20/NekoImageWorkflow/common/utils"
	apiModel "github.com/pk5ls20/NekoImageWorkflow/uploadClient/scraper/api/model"
	scraperModel "github.com/pk5ls20/NekoImageWorkflow/uploadClient/scraper/api/model"
	"github.com/pk5ls20/NekoImageWorkflow/uploadClient/storage/config"
	"github.com/sirupsen/logrus"
	"time"
)

type APIFetcher struct {
	apiModel.Fetcher
	fetchListParser     APIParser
	fetchListHttpClient *utils.HttpClient
	initialized         bool
}

func (a *APIFetcher) Init() error {
	logrus.Debug("APIFetcher Init start")
	a.fetchListParser.Init()
	a.fetchListHttpClient = utils.NewHttpClient()
	a.initialized = true
	logrus.Debug("APIFetcher initialized")
	return nil
}

func (a *APIFetcher) FetchList(cf []*config.APIScraperSourceConfig) ([]*scraperModel.SpiderToDoTask, error) {
	if !a.initialized {
		return nil, log.ErrorWrap(errors.New("APIFetcher not initialized"))
	}
	logrus.Debug("APIFetcher Fetch start")
	var infoList []*scraperModel.SpiderToDoTask
	var fetchTaskList []*apiModel.FetcherTaskList
	for _, c := range cf {
		if err := a.fetchListParser.RegisterParser(c.ParserJavaScriptFile); err != nil {
			return infoList, log.ErrorWrap(err)
		}
		fetchTaskList = append(fetchTaskList, &apiModel.FetcherTaskList{
			APIAddress:    c.APIAddress,
			PasteFilePath: c.ParserJavaScriptFile,
			Headers:       c.OptionalHeaders,
			Cookies:       c.OptionalCookies,
		})
	}
	for _, task := range fetchTaskList {
		response, err := a.fetchListHttpClient.Get(task.APIAddress, task.Headers, task.Cookies)
		if err != nil {
			return nil, log.ErrorWrap(err)
		}
		rawJson := string(response)
		parsedUrls, err := a.fetchListParser.ParseJson(rawJson, task.PasteFilePath)
		if err != nil {
			return nil, log.ErrorWrap(err)
		}
		for _, url := range parsedUrls {
			newTask := &scraperModel.SpiderToDoTask{
				SpiderTask: &scraperModel.SpiderTask{
					Url:     url,
					Headers: task.Headers,
					Cookies: task.Cookies,
				},
			}
			infoList = append(infoList, newTask)
		}
	}
	return infoList, nil
}

func (a *APIFetcher) FetchContent(task []*scraperModel.SpiderToDoTask) ([]*scraperModel.SpiderDoTask, error) {
	if !a.initialized {
		return nil, log.ErrorWrap(errors.New("APIFetcher not initialized"))
	}
	logrus.Debug("APIFetcher FetchContent start")
	spider := &APISpider{}
	spiderConfig := &scraperModel.SpiderConfig{
		SingleTaskMaxRetriesTime:    5,
		SingleTaskRetryDuration:     1000 * time.Millisecond,
		ConcurrentTaskLimit:         10,
		ConcurrentTaskGroupDuration: 300 * time.Millisecond,
		AdjustLimitRate:             0.3,
		AdjustLimitCheckTime:        500 * time.Millisecond,
	}
	ctx, cancel := context.WithCancel(context.Background())
	if err := spider.Init(task, spiderConfig, ctx, cancel); err != nil {
		return nil, log.ErrorWrap(err)
	}
	if err := spider.Start(); err != nil {
		return nil, log.ErrorWrap(err)
	}
	rs, err := spider.WaitDone()
	if err != nil {
		return nil, log.ErrorWrap(err)
	}
	return rs, nil
}
