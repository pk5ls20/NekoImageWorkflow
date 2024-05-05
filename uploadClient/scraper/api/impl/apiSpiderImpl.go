package impl

import (
	"context"
	"fmt"
	"github.com/pk5ls20/NekoImageWorkflow/common/log"
	scraperModels "github.com/pk5ls20/NekoImageWorkflow/uploadClient/scraper/api/model"
	"github.com/pk5ls20/NekoImageWorkflow/uploadClient/scraper/api/utils"
	"github.com/sirupsen/logrus"
	"net/http"
	"runtime"
	"sync"
	"time"
)

type APISpiderImpl struct {
	scraperModels.Spider
	// fetch list
	fetchList []string
	// context
	ctx    context.Context
	cancel context.CancelFunc
	// protected values
	initialized                    bool
	wg                             sync.WaitGroup
	httpClient                     *http.Client
	config                         *scraperModels.SpiderConfig
	initialConcurrentTaskLimit     int
	initialSingleTaskRetryDuration time.Duration
	taskPendingChannel             chan *scraperModels.SpiderTask
	taskDoneChannel                chan *scraperModels.SpiderTask
	dynamicSemaphore               *utils.DynamicSemaphore
	// protected locked values
	fetchAllTime            utils.LockValue[int]
	fetchSuccessTime        utils.LockValue[int]
	singleTaskRetryDuration utils.RWLockValue[time.Duration]
}

func (s *APISpiderImpl) Init(fetchList []string, config *scraperModels.SpiderConfig) error {
	// init fetch list
	s.fetchList = fetchList
	// init context
	s.ctx, s.cancel = context.WithCancel(context.Background())
	// init protected values
	s.wg = sync.WaitGroup{}
	s.httpClient = &http.Client{
		Transport: &http.Transport{
			MaxIdleConns:        1000,
			MaxIdleConnsPerHost: 1000,
			MaxConnsPerHost:     1000,
		},
		Timeout: 10 * time.Second,
	}
	s.config = config
	s.initialConcurrentTaskLimit = config.ConcurrentTaskLimit
	s.initialSingleTaskRetryDuration = config.SingleTaskRetryDuration
	s.taskPendingChannel = make(chan *scraperModels.SpiderTask, len(s.fetchList))
	s.taskDoneChannel = make(chan *scraperModels.SpiderTask, len(s.fetchList))
	s.dynamicSemaphore = &utils.DynamicSemaphore{
		SetVal:     s.config.ConcurrentTaskLimit,
		CurrentVal: 0,
	}
	s.dynamicSemaphore.Cond = sync.NewCond(&s.dynamicSemaphore.Mutex)
	// init protected locked values
	s.fetchAllTime = utils.LockValue[int]{
		Value: 0,
		Lock:  &sync.Mutex{},
	}
	s.fetchSuccessTime = utils.LockValue[int]{
		Value: 0,
		Lock:  &sync.Mutex{},
	}
	s.singleTaskRetryDuration = utils.RWLockValue[time.Duration]{
		Value: config.SingleTaskRetryDuration,
		Lock:  &sync.RWMutex{},
	}
	s.initialized = true
	return nil
}

func (s *APISpiderImpl) Start() error {
	if !s.initialized {
		return log.ErrorWrap(fmt.Errorf("APIParser not initialized"))
	}
	s.wg.Add(len(s.fetchList))
	go func() {
		for _, url := range s.fetchList {
			s.taskPendingChannel <- &scraperModels.SpiderTask{Url: url, TotalRetries: 0, Success: false}
		}
	}()
	go func() {
		for pdTask := range s.taskPendingChannel {
			if err := s.dynamicSemaphore.Acquire(context.Background()); err != nil {
				logrus.Error("Failed to acquire semaphore due to err", err)
				return
			}
			go func(t *scraperModels.SpiderTask) {
				defer func() {
					s.dynamicSemaphore.Release()
				}()
				logrus.Debug("Start to fetch ", t.Url)
				s.httpRequest(t)
				time.Sleep(s.config.ConcurrentTaskGroupDuration)
			}(pdTask)
		}
	}()
	go s.dynamicChangeFailTime()
	return nil
}

func (s *APISpiderImpl) WaitDone() ([]*scraperModels.SpiderTask, error) {
	doneTasks := make([]*scraperModels.SpiderTask, 0)
	if !s.initialized {
		return doneTasks, log.ErrorWrap(fmt.Errorf("APIParser not initialized"))
	}
	s.wg.Wait()
	s.cancel()
	close(s.taskPendingChannel)
	close(s.taskDoneChannel)
	for t := range s.taskDoneChannel {
		doneTasks = append(doneTasks, t)
	}
	logrus.Info("All tasks done, total tasks: ", len(doneTasks))
	return doneTasks, nil
}

func (s *APISpiderImpl) httpRequest(task *scraperModels.SpiderTask) {
	for {
		if task.TotalRetries > s.config.SingleTaskMaxRetriesTime {
			logrus.Warning("Failed to fetch ", task.Url,
				" after ", s.config.SingleTaskMaxRetriesTime, " retries")
			s.taskDoneChannel <- task
			s.wg.Done()
			return
		}
		resp, err := s.httpClient.Get(task.Url)
		s.fetchAllTime.Set(s.fetchAllTime.Get() + 1)
		task.TotalRetries++
		if err != nil {
			logrus.Errorf("Failed to fetch %s due to err: %v", task.Url, err)
			time.Sleep(s.singleTaskRetryDuration.Get())
			continue
		}
		if resp.StatusCode != 200 {
			logrus.Errorf("Failed to fetch %s due to statusCode: %s", task.Url, resp.Status)
			time.Sleep(s.singleTaskRetryDuration.Get())
			continue
		}
		task.Success = true
		logrus.Infof("Successfully fetched %s", task.Url)
		s.taskDoneChannel <- task
		s.wg.Done()
		s.fetchSuccessTime.Set(s.fetchSuccessTime.Get() + 1)
		return
	}
}

func (s *APISpiderImpl) dynamicChangeFailTime() {
	ticker := time.NewTicker(s.config.AdjustLimitCheckTime)
	defer ticker.Stop()
	for {
		select {
		case <-s.ctx.Done():
			logrus.Warning("Stopping dynamicChangeFailTime due to context cancellation.")
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
				oriLimit := s.config.ConcurrentTaskLimit
				decrementLimit := int(float64(s.config.ConcurrentTaskLimit) * 0.9)
				setLimit := max(decrementLimit, int(float64(s.initialConcurrentTaskLimit)*0.5))
				logrus.Warningf("Fail rate is too high, decrease concurrent task limit from %d -> %d",
					oriLimit, setLimit)
				s.config.ConcurrentTaskLimit = setLimit
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
				oriLimit := s.config.ConcurrentTaskLimit
				incrementLimit := int(float64(s.config.ConcurrentTaskLimit) * 1.1)
				setLimit := min(incrementLimit, int(float64(s.initialConcurrentTaskLimit)*1.5))
				logrus.Warningf("Fail rate is normal, increase concurrent task limit from %d -> %d",
					oriLimit, setLimit)
				s.config.ConcurrentTaskLimit = setLimit
				s.dynamicSemaphore.AdjustSize(setLimit)
			}
		}
	}
}
