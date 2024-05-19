package utils

import (
	"errors"
	"fmt"
	"github.com/pk5ls20/NekoImageWorkflow/common/log"
	"io"
	"net/http"
	"time"
)

type httpClient interface {
	Get(url string) (*http.Response, error)
	Do(req *http.Request) (*http.Response, error)
}

type HttpClient struct {
	client httpClient
}

func NewHttpClient() *HttpClient {
	return &HttpClient{
		client: &http.Client{
			Transport: &http.Transport{
				MaxIdleConns:        1000,
				MaxIdleConnsPerHost: 1000,
				MaxConnsPerHost:     1000,
			},
			Timeout: 10 * time.Second,
		},
	}
}

func (c *HttpClient) Get(url string, header map[string]string, cookie map[string]string) ([]byte, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	if header != nil {
		for k, v := range header {
			req.Header.Set(k, v)
		}
	}
	if cookie != nil {
		for k, v := range cookie {
			req.AddCookie(&http.Cookie{Name: k, Value: v})
		}
	}
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, log.ErrorWrap(err)
	}
	if resp == nil {
		return nil, log.ErrorWrap(errors.New("response is nil"))
	}
	if resp.StatusCode != 200 {
		return nil, log.ErrorWrap(errors.New(fmt.Sprintf("status code: %d", resp.StatusCode)))
	}
	defer func(Body io.ReadCloser) {
		if _err := Body.Close(); _err != nil {
			return
		}
	}(resp.Body)
	return io.ReadAll(resp.Body)
}
