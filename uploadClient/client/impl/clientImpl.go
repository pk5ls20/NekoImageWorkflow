package impl

import (
	"context"
	"github.com/pk5ls20/NekoImageWorkflow/common/log"
	uploadClient "github.com/pk5ls20/NekoImageWorkflow/uploadClient/kitex_gen/protoFile"
	"github.com/pk5ls20/NekoImageWorkflow/uploadClient/kitex_gen/protoFile/fileuploadservice"
	ScraperLifeCycle "github.com/pk5ls20/NekoImageWorkflow/uploadClient/scraper/lifecycle"
	"github.com/pk5ls20/NekoImageWorkflow/uploadClient/scraper/model"
	"github.com/pk5ls20/NekoImageWorkflow/uploadClient/storage/config"
	"github.com/pk5ls20/NekoImageWorkflow/uploadClient/storage/queue"
	"github.com/pk5ls20/NekoImageWorkflow/uploadClient/storage/sqlite"
	"github.com/sirupsen/logrus"
)

type client interface {
	// OnInit call after logger init, load client self config, init database, write database data to fileQueue
	OnInit() error
	// OnStart call after kitex's MustNewClient, aka after kitex client start and before kitex transport start
	OnStart() error
	// HandleFilePreUpload upload preUploadData to kitex server
	HandleFilePreUpload(ctx context.Context, cli fileuploadservice.Client) error
	// HandleFilePostUpload upload uploadData to kitex server
	HandleFilePostUpload(ctx context.Context, cli fileuploadservice.Client) error
	// OnStop call on program exit, write fileQueue data to database
	// TODO: maybe we can add database write logic here
	OnStop() error
}

type Client struct {
	client
	ClientInfo     *config.ClientConfig
	Scrapers       []model.Scraper
	PreUploadQueue *queue.PreUploadQueue
	UploadQueue    *queue.UploadQueue
}

// OnInit load client self config and before data, then init Scrapers
func (ci *Client) OnInit() error {
	// init
	sqlite.InitSqlite()
	logrus.Debug("Client OnInit start")
	ci.ClientInfo = config.GetConfig()
	return nil
}

// OnStart currently do nothing
func (ci *Client) OnStart() error {
	// TODO: make it really work
	logrus.Debug("Client OnStart start")
	ci.Scrapers = ScraperLifeCycle.RegisterScraper(ci.ClientInfo.Scraper)
	go ScraperLifeCycle.StartScraper(ci.Scrapers)
	return nil
}

// HandleFilePreUpload report pre upload data
func (ci *Client) HandleFilePreUpload(ctx context.Context, cli fileuploadservice.Client) error {
	// TODO: make it really work
	logrus.Debug("Client PreUpload start")
	req := &uploadClient.FilePreRequest{}
	if _, err := cli.HandleFilePreUpload(ctx, req); err != nil {
		return log.ErrorWrap(err)
	}
	return nil
}

// HandleFilePostUpload report post upload data
func (ci *Client) HandleFilePostUpload(ctx context.Context, cli fileuploadservice.Client) error {
	// TODO: make it really work
	logrus.Debug("Client PostUpload start")
	req := &uploadClient.FilePostRequest{}
	if _, err := cli.HandleFilePostUpload(ctx, req); err != nil {
		return log.ErrorWrap(err)
	}
	return nil
}

// OnStop write PreUploadQueue data and UploadQueue data to disk
func (ci *Client) OnStop() error {
	// TODO: make it really work, such as write flush channel data to disk
	// TODO: dynamic write sqlite data
	logrus.Debug("Client OnStop start")
	ScraperLifeCycle.StopScraper(ci.Scrapers)
	sqlite.CloseSqlite()
	return nil
}
