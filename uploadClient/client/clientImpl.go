package client

import (
	"context"
	"github.com/pk5ls20/NekoImageWorkflow/common/log"
	uploadClient "github.com/pk5ls20/NekoImageWorkflow/uploadClient/kitex_gen/protoFile"
	"github.com/pk5ls20/NekoImageWorkflow/uploadClient/kitex_gen/protoFile/fileuploadservice"
	"github.com/pk5ls20/NekoImageWorkflow/uploadClient/model"
	"github.com/pk5ls20/NekoImageWorkflow/uploadClient/scraper"
	"github.com/pk5ls20/NekoImageWorkflow/uploadClient/storage/bridge"
	"github.com/pk5ls20/NekoImageWorkflow/uploadClient/storage/config"
	"github.com/sirupsen/logrus"
)

type ClientImpl interface {
	// OnInit call after logger init, load client self config and before data
	OnInit() error
	// OnStart call after MustNewClient
	OnStart() error
	// HandleFilePreUpload report pre upload data
	HandleFilePreUpload(ctx context.Context, cli fileuploadservice.Client) error
	// HandleFilePostUpload report post upload data
	HandleFilePostUpload(ctx context.Context, cli fileuploadservice.Client) error
	// OnStop call on program exit
	OnStop() error
}

type ClientInstance struct {
	ClientImpl
	ClientInfo        *model.ClientConfig
	Scrapers          []scraper.ScraperInstance
	PreUploadBridge   *bridge.PreUploadTransBridgeInstance
	UploadTransBridge *bridge.UploadTransBridgeInstance
}

// OnInit load client self config and before data, then init Scrapers
func (ci *ClientInstance) OnInit() error {
	// init
	logrus.Debug("ClientInstance OnInit start")
	ci.ClientInfo = config.GetConfig()
	return nil
}

// OnStart currently do nothing
func (ci *ClientInstance) OnStart() error {
	// TODO: make it really work, such as write disk to flush channel data
	logrus.Debug("ClientInstance OnStart start")
	return nil
}

// HandleFilePreUpload report pre upload data
func (ci *ClientInstance) HandleFilePreUpload(ctx context.Context, cli fileuploadservice.Client) error {
	// TODO: make it really work
	logrus.Debug("ClientInstance PreUpload start")
	req := &uploadClient.FilePreRequest{}
	if _, err := cli.HandleFilePreUpload(ctx, req); err != nil {
		return log.ErrorWrap(err)
	}
	return nil
}

// HandleFilePostUpload report post upload data
func (ci *ClientInstance) HandleFilePostUpload(ctx context.Context, cli fileuploadservice.Client) error {
	// TODO: make it really work
	logrus.Debug("ClientInstance PostUpload start")
	req := &uploadClient.FilePostRequest{}
	if _, err := cli.HandleFilePostUpload(ctx, req); err != nil {
		return log.ErrorWrap(err)
	}
	return nil
}

// OnStop write PreUploadBridge data and UploadTransBridge data to disk
func (ci *ClientInstance) OnStop() error {
	// TODO: make it really work, such as write flush channel data to disk
	logrus.Debug("ClientInstance OnStop start")
	return nil
}
