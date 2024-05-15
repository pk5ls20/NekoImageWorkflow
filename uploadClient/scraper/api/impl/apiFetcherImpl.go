package impl

import (
	"github.com/pk5ls20/NekoImageWorkflow/common/log"
	"github.com/pk5ls20/NekoImageWorkflow/uploadClient/scraper/api/model"
	scraperModels "github.com/pk5ls20/NekoImageWorkflow/uploadClient/scraper/api/model"
	"github.com/pk5ls20/NekoImageWorkflow/uploadClient/scraper/api/utils"
	"github.com/pk5ls20/NekoImageWorkflow/uploadClient/storage/config"
	"github.com/sirupsen/logrus"
	"time"
)

type APIFetcher struct {
	model.Fetcher
	parser        APIParser
	spider        APISpider
	spiderConfig  *scraperModels.SpiderConfig
	httpClient    *utils.HttpClient
	apiImplConfig *[]config.APIScraperSourceConfig
	fetchTaskList *[]model.FetcherTaskList
}

func (a *APIFetcher) Init(cf *[]config.APIScraperSourceConfig) error {
	logrus.Debug("APIFetcher Init start")
	a.apiImplConfig = cf
	a.parser.Init()
	a.httpClient = utils.NewHttpClient()
	a.fetchTaskList = &[]model.FetcherTaskList{}
	for _, c := range *cf {
		if err := a.parser.RegisterParser(c.ParserJavaScriptFile); err != nil {
			return log.ErrorWrap(err)
		}
		*a.fetchTaskList = append(*a.fetchTaskList, model.FetcherTaskList{
			APIAddress:    c.APIAddress,
			PasteFilePath: c.ParserJavaScriptFile,
			Headers:       c.OptionalHeaders,
			Cookies:       c.OptionalCookies,
		})
	}
	logrus.Debug("APIFetcher initialized")
	return nil
}

func (a *APIFetcher) FetchList() ([]*scraperModels.SpiderToDoTask, error) {
	logrus.Debug("APIFetcher Fetch start")
	var infoList []*scraperModels.SpiderToDoTask
	for _, task := range *a.fetchTaskList {
		response, err := a.httpClient.Get(task.APIAddress, task.Headers, task.Cookies)
		if err != nil {
			return nil, log.ErrorWrap(err)
		}
		rawJson := string(response)
		parsedUrls, err := a.parser.ParseJson(rawJson, task.PasteFilePath)
		if err != nil {
			return nil, log.ErrorWrap(err)
		}
		for _, url := range parsedUrls {
			newTask := &scraperModels.SpiderToDoTask{
				SpiderTask: &scraperModels.SpiderTask{
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

func (a *APIFetcher) FetchContent(task []*scraperModels.SpiderToDoTask) ([]*model.SpiderDoTask, error) {
	logrus.Debug("APIFetcher FetchContent start")
	a.spiderConfig = &scraperModels.SpiderConfig{
		SingleTaskMaxRetriesTime:    5,
		SingleTaskRetryDuration:     1000 * time.Millisecond,
		ConcurrentTaskLimit:         10,
		ConcurrentTaskGroupDuration: 300 * time.Millisecond,
		AdjustLimitRate:             0.3,
		AdjustLimitCheckTime:        500 * time.Millisecond,
	}
	if err := a.spider.Init(task, a.spiderConfig); err != nil {
		return nil, log.ErrorWrap(err)
	}
	if err := a.spider.Start(); err != nil {
		return nil, log.ErrorWrap(err)
	}
	rs, err := a.spider.WaitDone()
	if err != nil {
		return nil, log.ErrorWrap(err)
	}
	return rs, nil
}
