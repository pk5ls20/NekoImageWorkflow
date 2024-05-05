package model

import "time"

type SpiderTasks = []*SpiderTask

type SpiderTask struct {
	Url          string
	TotalRetries int
	Success      bool
}

type Spider interface {
	Init(fetchList []string, config *SpiderConfig) error
	Start() error
	WaitDone() ([]*SpiderTask, error)
	httpRequest(task *SpiderTask) error
}

type SpiderConfig struct {
	SingleTaskRetryDuration     time.Duration
	SingleTaskMaxRetriesTime    int
	ConcurrentTaskLimit         int
	ConcurrentTaskGroupDuration time.Duration
	AdjustLimitRate             float64
	AdjustLimitCheckTime        time.Duration
}
