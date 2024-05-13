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
	fetchTaskList *[]model.FetcherTask
}

func (a *APIFetcher) Init(cf *[]config.APIScraperSourceConfig) error {
	// TODO: load config cookie & header config
	logrus.Debug("APIFetcher Init start")
	a.apiImplConfig = cf
	a.parser.Init()
	a.httpClient = utils.NewHttpClient()
	a.fetchTaskList = &[]model.FetcherTask{}
	for _, c := range *cf {
		if err := a.parser.RegisterParser(c.ParserJavaScriptFile); err != nil {
			return log.ErrorWrap(err)
		}
		*a.fetchTaskList = append(*a.fetchTaskList, model.FetcherTask{
			APIAddress:    c.APIAddress,
			PasteFilePath: c.ParserJavaScriptFile,
		})
	}
	logrus.Debug("APIFetcher initialized")
	return nil
}

func (a *APIFetcher) FetchList() (url []string, err error) {
	logrus.Debug("APIFetcher Fetch start")
	infoList := make([]string, 0)
	for _, t := range *a.fetchTaskList {
		res, _err := a.httpClient.Get(t.APIAddress)
		if _err != nil {
			return nil, log.ErrorWrap(_err)
		}
		rawJson := string(res)
		parsed, _err := a.parser.ParseJson(rawJson, t.PasteFilePath)
		if _err != nil {
			return nil, log.ErrorWrap(_err)
		}
		infoList = append(infoList, parsed...)
	}
	return infoList, nil
}

func (a *APIFetcher) FetchContent(url []string) ([]*model.SpiderTask, error) {
	logrus.Debug("APIFetcher FetchContent start")
	a.spiderConfig = &scraperModels.SpiderConfig{
		SingleTaskMaxRetriesTime:    5,
		SingleTaskRetryDuration:     1000 * time.Millisecond,
		ConcurrentTaskLimit:         10,
		ConcurrentTaskGroupDuration: 300 * time.Millisecond,
		AdjustLimitRate:             0.3,
		AdjustLimitCheckTime:        500 * time.Millisecond,
	}
	if err := a.spider.Init(url, a.spiderConfig); err != nil {
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
