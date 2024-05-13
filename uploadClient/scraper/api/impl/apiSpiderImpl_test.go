package impl

import (
	_ "github.com/pk5ls20/NekoImageWorkflow/common/log"
	scraperModels "github.com/pk5ls20/NekoImageWorkflow/uploadClient/scraper/api/model"
	"github.com/sirupsen/logrus"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"
)

func TestAPISpiderWithMockServer(t *testing.T) {
	var count int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		randVal := rand.Float64()
		if randVal < 0.1 {
			return
		}
		if atomic.LoadInt32(&count) < 200 && randVal < 0.5 {
			atomic.AddInt32(&count, 1)
			w.WriteHeader(http.StatusRequestTimeout)
			return
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()
	spider := &APISpider{}
	urls := make([]string, 1000)
	for i := 0; i < 1000; i++ {
		urls[i] = server.URL
	}
	config := &scraperModels.SpiderConfig{
		SingleTaskMaxRetriesTime:    1,
		SingleTaskRetryDuration:     100 * time.Millisecond,
		ConcurrentTaskLimit:         10,
		ConcurrentTaskGroupDuration: 200 * time.Millisecond,
		AdjustLimitRate:             0.3,
		AdjustLimitCheckTime:        100 * time.Millisecond,
	}
	// first, not init
	if err := spider.Start(); err == nil {
		t.Error("Should not Start without init")
	}
	if _, err := spider.WaitDone(); err == nil {
		t.Error("Should not WaitDone without init")
	}
	if err := spider.Init(urls, config); err != nil {
		t.Error(err)
	}
	if err := spider.Start(); err != nil {
		t.Error(err)
	}
	rs, err := spider.WaitDone()
	if err != nil {
		t.Error(err)
	}
	expectedResultCount := len(urls)
	actualResultCount := len(rs)
	if actualResultCount != expectedResultCount {
		t.Errorf("Expected %d successful tasks, but got %d", expectedResultCount, actualResultCount)
	}
	logrus.Info("Total success tasks: ", actualResultCount)
}
