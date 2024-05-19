package impl

import (
	"context"
	"errors"
	"fmt"
	"github.com/pk5ls20/NekoImageWorkflow/uploadClient/storage/config"
	"github.com/sirupsen/logrus"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"sync"
	"testing"
)

var jsonGen = func(st string, p string, l1 string, l2 string) string {
	return fmt.Sprintf(`{
	"%s": [
		{"link": "%s/%s"},
		{"link": "%s/%s"}
	]
	}`, st, p, l1, p, l2)
}

const jsCode1 = `
function pasteJson(json) {
	var data = JSON.parse(json);
	var urls = [];
	for (var i = 0; i < data.websites.length; i++) {
		urls.push(data.websites[i].link);
	}
	return urls;
}
`

const jsCode2 = `
function pasteJson(json) {
	var data = JSON.parse(json);
	var urls = [];
	for (var i = 0; i < data.urls.length; i++) {
		urls.push(data.urls[i].link);
	}
	return urls;
}
`

var routeMap = func(prefix string) map[string]string {
	return map[string]string{
		"/json1": jsonGen("websites", prefix, "1", "2"),
		"/json2": jsonGen("urls", prefix, "3", "4"),
		"/1":     "1",
		"/2":     "2",
		"/3":     "3",
		"/4":     "4",
	}
}

// doFetch TODO: check the fetched contents
func doFetch(t *testing.T, fetcher *APIFetcher, scraperID int, cf []*config.APIScraperSourceConfig) error {
	var err error
	tasks, err := fetcher.FetchList(cf)
	if err != nil {
		err = fmt.Errorf("failed to fetch list: %v", err)
	}
	// FetchContent
	ctx, cancel := context.WithCancel(context.Background())
	contents, err := fetcher.FetchContent(tasks, ctx, cancel)
	if err != nil {
		err = fmt.Errorf("failed to fetch content: %v", err)
	}
	// check scraper id
	for _, content := range contents {
		if content.ScraperID != scraperID {
			t.Errorf("scraper id not match: %d, %d", content.ScraperID, scraperID)
		}
	}
	logrus.Info("scraperID #", scraperID, " Contents: ", contents)
	return err
}

func TestAPIFetcherImpl_FetchList(t *testing.T) {
	// init tempdir
	tempDir, _err := os.MkdirTemp("", "test")
	if _err != nil {
		t.Fatalf("Failed to create temp directory: %v", _err)
	}
	defer func(path string) {
		if err := os.RemoveAll(path); err != nil {
			logrus.Error("Failed to remove temp directory: ", err)
		}
	}(tempDir)
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
	// init mock httpserver
	var ts *httptest.Server
	ts = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if data, ok := routeMap(ts.URL)[r.URL.Path]; ok {
			w.WriteHeader(http.StatusOK)
			if _, err := fmt.Fprintln(w, data); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		} else {
			http.NotFound(w, r)
		}
	}))
	defer ts.Close()
	// init jsFile
	jsFile1 := filepath.Join(tempDir, "1.js")
	if err := os.WriteFile(jsFile1, []byte(jsCode1), 0644); err != nil {
		t.Fatalf("Failed to write js1 to temp file: %v", err)
	}
	jsFile2 := filepath.Join(tempDir, "2.js")
	if err := os.WriteFile(jsFile2, []byte(jsCode2), 0644); err != nil {
		t.Fatalf("Failed to write js2 to temp file: %v", err)
	}
	// init config
	var cf []*config.APIScraperSourceConfig
	cf = append(cf, &config.APIScraperSourceConfig{
		APIAddress:           ts.URL + "/json1",
		ParserJavaScriptFile: jsFile1,
		OptionalHeaders:      map[string]string{"User-Agent": "test"},
		OptionalCookies:      map[string]string{"Cookie": "test"},
	})
	cf = append(cf, &config.APIScraperSourceConfig{
		APIAddress:           ts.URL + "/json2",
		ParserJavaScriptFile: jsFile2,
		OptionalHeaders:      map[string]string{"User-Agent": "test"},
		OptionalCookies:      map[string]string{"Cookie": "test"},
	})
	// init fetcher
	fetcher := &APIFetcher{}
	if err := doFetch(t, fetcher, 0, cf); err == nil {
		t.Fatalf("expected error, got nil")
	}
	var wg sync.WaitGroup
	var scraperLen = 100
	var scraperIDList = make([]int, scraperLen)
	for i := 0; i < scraperLen; i++ {
		scraperIDList[i] = i
	}
	wg.Add(len(scraperIDList))
	for _, scpID := range scraperIDList {
		go func(id int) {
			defer wg.Done()
			ft := &APIFetcher{}
			if err := ft.Init(id); err != nil {
				t.Errorf("Failed to init fetcher: %v", err)
				return
			}
			if err := doFetch(t, ft, id, cf); err != nil {
				t.Errorf("don't expect error, got: %v", err)
				return
			}
		}(scpID)
	}
	wg.Wait()
}
