package model

import (
	"github.com/pk5ls20/NekoImageWorkflow/uploadClient/client/model"
	"time"
)

type SpiderTask struct {
	Url          string
	TotalRetries int
	Success      bool
	FetchData    model.UploadFileDataModel
}

type Spider interface {
	Init(fetchList []string, config *SpiderConfig) error
	Start() error
	WaitDone() ([]*SpiderTask, error)
	httpRequest(task *SpiderTask)
}

type SpiderConfig struct {
	SingleTaskRetryDuration     time.Duration
	SingleTaskMaxRetriesTime    int
	ConcurrentTaskLimit         int
	ConcurrentTaskGroupDuration time.Duration
	AdjustLimitRate             float64
	AdjustLimitCheckTime        time.Duration
}
