package impl

import (
	"context"
	"fmt"
	"github.com/pk5ls20/NekoImageWorkflow/common/log"
	commonModel "github.com/pk5ls20/NekoImageWorkflow/common/model"
	"github.com/pk5ls20/NekoImageWorkflow/uploadClient/client/model"
	scraperModels "github.com/pk5ls20/NekoImageWorkflow/uploadClient/scraper/api/model"
	"github.com/pk5ls20/NekoImageWorkflow/uploadClient/scraper/api/utils"
	"github.com/sirupsen/logrus"
	"runtime"
	"sync"
	"time"
)

type APISpider struct {
	scraperModels.Spider
	// fetch list
	fetchList []*scraperModels.SpiderToDoTask
	// context
	ctx    context.Context
	cancel context.CancelFunc
	// protected values
	initialized                    bool
	wg                             sync.WaitGroup
	httpClient                     *utils.HttpClient
	config                         *scraperModels.SpiderConfig
	initialConcurrentTaskLimit     int
	initialSingleTaskRetryDuration time.Duration
	taskPendingChannel             chan *scraperModels.SpiderDoTask
	taskDoneChannel                chan *scraperModels.SpiderDoTask
	dynamicSemaphore               *utils.DynamicSemaphore
	// protected locked values
	fetchAllTime            *utils.LockValue[int]
	fetchSuccessTime        *utils.LockValue[int]
	concurrentTaskLimit     *utils.RWLockValue[int]
	singleTaskRetryDuration *utils.RWLockValue[time.Duration]
}

func (s *APISpider) Init(fetchTaskList []*scraperModels.SpiderToDoTask, config *scraperModels.SpiderConfig,
	ctx context.Context, cancel context.CancelFunc) error {
	s.fetchList = fetchTaskList
	// init context
	s.ctx, s.cancel = ctx, cancel
	// init protected values
	s.wg = sync.WaitGroup{}
	s.httpClient = utils.NewHttpClient()
	s.config = config
	s.initialConcurrentTaskLimit = config.ConcurrentTaskLimit
	s.initialSingleTaskRetryDuration = config.SingleTaskRetryDuration
	s.taskPendingChannel = make(chan *scraperModels.SpiderDoTask, len(s.fetchList))
	s.taskDoneChannel = make(chan *scraperModels.SpiderDoTask, len(s.fetchList))
	s.dynamicSemaphore = utils.NewDynamicSemaphore(config.ConcurrentTaskLimit)
	// init protected locked values
	s.fetchAllTime = utils.NewLockValue[int](0)
	s.fetchSuccessTime = utils.NewLockValue[int](0)
	s.concurrentTaskLimit = utils.NewRWLockValue[int](s.initialConcurrentTaskLimit)
	s.singleTaskRetryDuration = utils.NewRWLockValue[time.Duration](s.initialSingleTaskRetryDuration)
	s.initialized = true
	return nil
}

func (s *APISpider) Start() error {
	if !s.initialized {
		return log.ErrorWrap(fmt.Errorf("apiParser not initialized"))
	}
	s.wg.Add(len(s.fetchList))
	go func() {
		for _, task := range s.fetchList {
			s.taskPendingChannel <- &scraperModels.SpiderDoTask{
				TotalRetries: 0,
				Success:      false,
				SpiderTask:   task.SpiderTask,
			}
		}
	}()
	go func() {
		for pdTask := range s.taskPendingChannel {
			select {
			case <-s.ctx.Done():
				logrus.Debug("Stopping Start due to context cancellation.")
				s.taskDoneChannel <- pdTask
				s.wg.Done()
			default:
				if err := s.dynamicSemaphore.Acquire(s.ctx); err != nil {
					logrus.Error("Failed to acquire semaphore due to err", err)
					return
				}
				go func(t *scraperModels.SpiderDoTask) {
					defer s.dynamicSemaphore.Release()
					logrus.Debug("Start to fetch ", t.Url)
					s.httpRequest(t)
					time.Sleep(s.config.ConcurrentTaskGroupDuration)
				}(pdTask)
			}
		}
	}()
	go s.dynamicChangeFailTime()
	return nil
}

// Cancel will directly cancel all tasks, means WaitDone stops being blocked and return all results
func (s *APISpider) Cancel() error {
	if !s.initialized {
		return log.ErrorWrap(fmt.Errorf("apiParser not initialized"))
	}
	logrus.Infof("Task Done triggered, having %d tasks left, %d tasks done",
		len(s.taskPendingChannel), len(s.taskDoneChannel))
	s.cancel()
	return nil
}

// WaitDone return all results, will block until all tasks are done except manually Cancel
func (s *APISpider) WaitDone() ([]*scraperModels.SpiderDoTask, error) {
	doneTasks := make([]*scraperModels.SpiderDoTask, 0)
	if !s.initialized {
		return doneTasks, log.ErrorWrap(fmt.Errorf("apiParser not initialized"))
	}
	s.wg.Wait()
	if err := s.Cancel(); err != nil {
		return nil, err
	}
	close(s.taskPendingChannel)
	close(s.taskDoneChannel)
	for t := range s.taskDoneChannel {
		doneTasks = append(doneTasks, t)
	}
	return doneTasks, nil
}

func (s *APISpider) httpRequest(task *scraperModels.SpiderDoTask) {
	for {
		select {
		case <-s.ctx.Done():
			logrus.Debug("Stopping httpRequest due to context cancellation.")
			s.taskDoneChannel <- task
			s.wg.Done()
			return
		default:
			if task.TotalRetries > s.config.SingleTaskMaxRetriesTime {
				logrus.Warning("Failed to fetch ", task.Url,
					" after ", s.config.SingleTaskMaxRetriesTime, " retries")
				s.taskDoneChannel <- task
				s.wg.Done()
				return
			}
			resData, resError := s.httpClient.Get(task.Url, task.Headers, task.Cookies)
			s.fetchAllTime.Set(s.fetchAllTime.Get() + 1)
			task.TotalRetries++
			if resError != nil {
				logrus.Errorf("Failed to fetch %s due to err: %v", task.Url, resError)
				time.Sleep(s.singleTaskRetryDuration.Get())
				continue
			}
			if fetchData, err := model.NewScraperUploadTempFileData(commonModel.APIScraperType, resData); err != nil {
				logrus.Errorf("Failed to create temp file data due to err: %v", err)
				time.Sleep(s.singleTaskRetryDuration.Get())
				continue
			} else {
				task.FetchData = *fetchData
			}
			task.Success = true
			logrus.Infof("Successfully fetched %s", task.Url)
			s.taskDoneChannel <- task
			s.wg.Done()
			s.fetchSuccessTime.Set(s.fetchSuccessTime.Get() + 1)
			return
		}
	}
}

// TODO: Use more flexible flow control algorithms
func (s *APISpider) dynamicChangeFailTime() {
	ticker := time.NewTicker(s.config.AdjustLimitCheckTime)
	defer ticker.Stop()
	for {
		select {
		case <-s.ctx.Done():
			logrus.Debug("Stopping dynamicChangeFailTime due to context cancellation.")
			return
		case <-ticker.C:
			logrus.Info("Initial goroutines:", runtime.NumGoroutine())
			failRate := float64(s.fetchAllTime.Get()-s.fetchSuccessTime.Get()) / float64(s.fetchAllTime.Get())
			logrus.Infof("Fetch Success / Fetch all: %d / %d", s.fetchSuccessTime.Get(), s.fetchAllTime.Get())
			logrus.Infof("Fail rate: %.2f", failRate)
			if failRate > s.config.AdjustLimitRate {
				// set fail task wait time
				ori := s.singleTaskRetryDuration.Get()
				increment := time.Duration(float64(ori) * 1.1)
				set := min(increment, 2*time.Duration(s.initialSingleTaskRetryDuration))
				logrus.Warningf("Fail rate is too high, increase retry duration from %s -> %s",
					ori.String(), set.String())
				s.singleTaskRetryDuration.Set(set)
				// set concurrent task limit
				oriLimit := s.concurrentTaskLimit.Get()
				decrementLimit := int(float64(oriLimit) * 0.9)
				setLimit := max(decrementLimit, int(float64(s.initialConcurrentTaskLimit)*0.5))
				logrus.Warningf("Fail rate is too high, decrease concurrent task limit from %d -> %d",
					oriLimit, setLimit)
				s.concurrentTaskLimit.Set(setLimit)
				s.dynamicSemaphore.AdjustSize(setLimit)
			} else {
				// set fail task wait time
				ori := s.singleTaskRetryDuration.Get()
				increment := time.Duration(float64(ori) * 0.9)
				set := max(increment, s.initialSingleTaskRetryDuration)
				logrus.Warningf("Fail rate is normal, decrease retry duration from %s -> %s",
					ori.String(), set.String())
				s.singleTaskRetryDuration.Set(set)
				// set concurrent task limit
				oriLimit := s.concurrentTaskLimit.Get()
				incrementLimit := int(float64(oriLimit) * 1.1)
				setLimit := min(incrementLimit, int(float64(s.initialConcurrentTaskLimit)*1.5))
				logrus.Warningf("Fail rate is normal, increase concurrent task limit from %d -> %d",
					oriLimit, setLimit)
				s.concurrentTaskLimit.Set(setLimit)
				s.dynamicSemaphore.AdjustSize(setLimit)
			}
		}
	}
}
