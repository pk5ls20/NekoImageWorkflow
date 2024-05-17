package model

import (
	"context"
	"github.com/pk5ls20/NekoImageWorkflow/uploadClient/client/model"
	"time"
)

// SpiderTask is the base task for spider
// NOTE: SpiderTasks are not goroutine-safe to write to
// which means that they should not be changed once they have been initialised.
type SpiderTask struct {
	Url     string
	Headers map[string]string
	Cookies map[string]string
}

// SpiderToDoTask is the task to do
type SpiderToDoTask struct {
	*SpiderTask
}

// SpiderDoTask completes its init in func (s *APISpider) Init and returns after WaitDone()
type SpiderDoTask struct {
	*SpiderTask
	TotalRetries int
	Success      bool
	FetchData    model.UploadFileDataModel
}

type Spider interface {
	Init(fetchTaskList []*SpiderToDoTask, config *SpiderConfig, ctx context.Context, cancel context.CancelFunc) error
	Start() error
	Cancel() error
	WaitDone() ([]*SpiderDoTask, error)
	httpRequest(task *SpiderDoTask)
}

type SpiderConfig struct {
	SingleTaskRetryDuration     time.Duration
	SingleTaskMaxRetriesTime    int
	ConcurrentTaskLimit         int
	ConcurrentTaskGroupDuration time.Duration
	AdjustLimitRate             float64
	AdjustLimitCheckTime        time.Duration
}
