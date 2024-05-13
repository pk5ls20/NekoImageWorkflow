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
	header string
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
		// TODO: separate this to a config file as a demo
		header: "Mozilla/5.0 (Windows NT 10.0; Win64; x64) " +
			"AppleWebKit/537.36 (KHTML, like Gecko) Chrome/124.0.0.0 Safari/537.36 Edg/124.0.0.0",
	}
}

func (c *HttpClient) Get(url string) ([]byte, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", c.header)
	resp, err := c.client.Do(req)
	if resp.StatusCode != 200 {
		return nil, log.ErrorWrap(errors.New(fmt.Sprintf("status code: %d", resp.StatusCode)))
	}
	if err != nil {
		return nil, log.ErrorWrap(err)
	}
	defer func(Body io.ReadCloser) {
		if _err := Body.Close(); _err != nil {
			return
		}
	}(resp.Body)
	return io.ReadAll(resp.Body)
}
