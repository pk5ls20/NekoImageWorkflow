package impl

import (
	"context"
	commonLog "github.com/pk5ls20/NekoImageWorkflow/common/log"
	clientModel "github.com/pk5ls20/NekoImageWorkflow/uploadClient/client/model"
	kitexUploadClient "github.com/pk5ls20/NekoImageWorkflow/uploadClient/kitex_gen/protoFile"
	kitexUploadService "github.com/pk5ls20/NekoImageWorkflow/uploadClient/kitex_gen/protoFile/fileuploadservice"
	scraperLifeCycle "github.com/pk5ls20/NekoImageWorkflow/uploadClient/scraper/lifecycle"
	scraperModel "github.com/pk5ls20/NekoImageWorkflow/uploadClient/scraper/model"
	storageConfig "github.com/pk5ls20/NekoImageWorkflow/uploadClient/storage/config"
	storageQueue "github.com/pk5ls20/NekoImageWorkflow/uploadClient/storage/queue"
	storageSqlite "github.com/pk5ls20/NekoImageWorkflow/uploadClient/storage/sqlite"
	"github.com/sirupsen/logrus"
)

type client interface {
	// OnInit call after logger init, load client self config, init database, write database data to fileQueue
	OnInit() error
	// OnStart call after kitex's MustNewClient, aka after kitex client start and before kitex transport start
	OnStart() error
	// HandleFilePreUpload upload preUploadData to kitex server
	HandleFilePreUpload(ctx context.Context, cli kitexUploadService.Client) error
	// HandleFilePostUpload upload uploadData to kitex server
	HandleFilePostUpload(ctx context.Context, cli kitexUploadService.Client) error
	// OnStop call on program exit, write fileQueue data to database
	// TODO: maybe we can add database write logic here
	OnStop() error
}

type Client struct {
	client
	ClientInfo     *storageConfig.ClientConfig
	Scrapers       []scraperModel.Scraper
	PreUploadQueue *storageQueue.PreUploadQueue
	UploadQueue    *storageQueue.UploadQueue
	ScraperChanMap scraperModel.ScraperChanMap
}

// OnInit load client self config and before data, then init Scrapers
func (ci *Client) OnInit() error {
	// init
	storageSqlite.InitSqlite()
	logrus.Debug("Client OnInit start")
	ci.PreUploadQueue = storageQueue.GetPreUploadQueue()
	ci.UploadQueue = storageQueue.GetUploadQueue()
	ci.ClientInfo = storageConfig.GetConfig()
	return nil
}

// OnStart currently do nothing
func (ci *Client) OnStart() error {
	// TODO: make it really work
	logrus.Debug("Client OnStart start")
	ci.Scrapers = scraperLifeCycle.RegisterScraper(ci.ClientInfo.Scraper)
	go scraperLifeCycle.StartScraper(ci.Scrapers)
	return nil
}

// HandleFilePreUpload report pre upload data
// TODO: Need to store filedata that failed to upload
func (ci *Client) HandleFilePreUpload(ctx context.Context, cli kitexUploadService.Client) error {
	// TODO: make it really work
	logrus.Debug("Client PreUpload start")
	err := ci.PreUploadQueue.Iterate(func(fileData *clientModel.PreUploadFileDataModel) error {
		req := &kitexUploadClient.FilePreRequest{}
		if _, err := cli.HandleFilePreUpload(ctx, req); err != nil {
			return commonLog.ErrorWrap(err)
		}
		return nil
	})
	if err != nil {
		return commonLog.ErrorWrap(err)
	}
	return nil
}

// HandleFilePostUpload report post upload data
// TODO: Need to store filedata that failed to upload
func (ci *Client) HandleFilePostUpload(ctx context.Context, cli kitexUploadService.Client) error {
	// TODO: make it really work
	logrus.Debug("Client PostUpload start")
	err := ci.UploadQueue.Iterate(func(fileData *clientModel.UploadFileDataModel) error {
		req := &kitexUploadClient.FilePostRequest{}
		if _, err := cli.HandleFilePostUpload(ctx, req); err != nil {
			return commonLog.ErrorWrap(err)
		}
		return nil
	})
	if err != nil {
		return commonLog.ErrorWrap(err)
	}
	return nil
}

// OnStop write PreUploadQueue data and UploadQueue data to disk
func (ci *Client) OnStop() error {
	// TODO: make it really work, such as write flush channel data to disk
	// TODO: dynamic write sqlite data
	logrus.Debug("Client OnStop start")
	scraperLifeCycle.StopScraper(ci.Scrapers)
	storageSqlite.CloseSqlite()
	return nil
}
