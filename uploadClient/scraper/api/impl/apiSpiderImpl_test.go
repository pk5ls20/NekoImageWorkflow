package impl

import (
	"context"
	"errors"
	"github.com/google/uuid"
	_ "github.com/pk5ls20/NekoImageWorkflow/common/log"
	commonUUID "github.com/pk5ls20/NekoImageWorkflow/common/uuid"
	scraperModels "github.com/pk5ls20/NekoImageWorkflow/uploadClient/scraper/api/model"
	"github.com/sirupsen/logrus"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"strconv"
	"sync/atomic"
	"testing"
	"time"
	"unsafe"
)

// uuid must contain, and contain only once
func containsOnce[T comparable](slice []T, a T) bool {
	countMap := make(map[T]int)
	for _, element := range slice {
		countMap[element]++
	}
	return countMap[a] == 1
}

func realTest(t *testing.T, spider *APISpider, tasks []*scraperModels.SpiderToDoTask,
	config *scraperModels.SpiderConfig, expectedUUID []uuid.UUID, expTime time.Duration) {
	ctx, cancel := context.WithCancel(context.Background())
	if err := spider.Init(tasks, config, ctx, cancel); err != nil {
		t.Error(err)
	}
	if err := spider.Start(); err != nil {
		t.Error(err)
	}
	// cancel after expTime
	go func() {
		time.Sleep(expTime)
		if err := spider.Cancel(); err != nil {
			return
		}
	}()
	rs, err := spider.WaitDone()
	if err != nil {
		t.Error(err)
	}
	expectedResultCount := len(tasks)
	actualResultCount := len(rs)
	if actualResultCount != expectedResultCount {
		t.Errorf("Expected %d successful tasks, but got %d", expectedResultCount, actualResultCount)
	}
	for i, result := range rs {
		if result.Success {
			val := reflect.ValueOf(result.FetchData).Elem()
			fileUUIDField := val.FieldByName("fileUUID")
			if !fileUUIDField.IsValid() {
				t.Errorf("Task %d: field 'fileUUID' not found", i)
				continue
			}
			// TODO: Since the final form of the UploadFileDataModel is undetermined,
			//  we'll first force it to get its fileUUID as an under-exported field here
			fileUUID := reflect.NewAt(fileUUIDField.Type(),
				unsafe.Pointer(fileUUIDField.UnsafeAddr())).Elem().Interface().(uuid.UUID)
			if !containsOnce(expectedUUID, fileUUID) {
				t.Errorf("Task %d: fileUUID not found in expectedUUID or contain > 1", i)
			}
		}
	}
	logrus.Info("Total success tasks: ", actualResultCount)
}

func TestAPISpiderWithMockServer(t *testing.T) {
	if _, err := os.Stat("_tmp"); errors.Is(err, os.ErrNotExist) {
		if err := os.Mkdir("_tmp", 0755); err != nil {
			t.Fatalf("Failed to create temp directory: %v", err)
		}
	}
	defer func() {
		if err := os.RemoveAll("_tmp"); err != nil {
			logrus.Error("Error removing _tmp folder: ", err)
		}
		logrus.Info("tmp folder removed")
	}()
	var count int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		randVal := rand.Float64()
		if randVal < 0.1 {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if atomic.LoadInt32(&count) < 200 && randVal < 0.5 {
			atomic.AddInt32(&count, 1)
			w.WriteHeader(http.StatusRequestTimeout)
			return
		}
		userAgent := r.Header.Get("User-Agent")
		testCookie, err := r.Cookie("TestCookie")
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		responseContent := "User-Agent: " + userAgent + ", TestCookie: " + testCookie.Value
		if _, err := w.Write([]byte(responseContent)); err != nil {
			return
		}
	}))
	defer server.Close()
	spider := &APISpider{}
	tasks := make([]*scraperModels.SpiderToDoTask, 1000)
	expectedUUID := make([]uuid.UUID, 1000)
	for i := 0; i < 1000; i++ {
		userAgent := "TestUserAgent" + strconv.Itoa(i)
		testCookie := "TestValue" + strconv.Itoa(i)
		tasks[i] = &scraperModels.SpiderToDoTask{
			SpiderTask: &scraperModels.SpiderTask{
				Url: server.URL,
				Headers: map[string]string{
					"User-Agent": userAgent,
				},
				Cookies: map[string]string{
					"TestCookie": testCookie,
				},
			},
		}
		fileContent := "User-Agent: " + userAgent + ", TestCookie: " + testCookie
		expectedUUID[i] = commonUUID.GenerateStrUUID(fileContent)
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
		t.Error("Should not Started without init")
	}
	if err := spider.Cancel(); err == nil {
		t.Error("Should not Canceled without init")
	}
	if _, err := spider.WaitDone(); err == nil {
		t.Error("Should not WaitDone without init")
	}
	// test without timeout
	realTest(t, &APISpider{}, tasks, config, expectedUUID, 114514*time.Hour) // never timeout
	// test with timeout
	realTest(t, &APISpider{}, tasks, config, expectedUUID, 1*time.Microsecond)
	realTest(t, &APISpider{}, tasks, config, expectedUUID, 1*time.Millisecond)
	realTest(t, &APISpider{}, tasks, config, expectedUUID, 1*time.Second)
	realTest(t, &APISpider{}, tasks, config, expectedUUID, 10*time.Second)
	realTest(t, &APISpider{}, tasks, config, expectedUUID, 30*time.Second)
	realTest(t, &APISpider{}, tasks, config, expectedUUID, 60*time.Second)
}
