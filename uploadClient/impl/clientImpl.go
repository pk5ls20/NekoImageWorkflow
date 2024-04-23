package impl

import (
	"NekoImageWorkflowKitex/uploadClient/kitex_gen/uploadClient"
	"NekoImageWorkflowKitex/uploadClient/kitex_gen/uploadClient/fileuploadservice"
	"NekoImageWorkflowKitex/uploadClient/model"
	"NekoImageWorkflowKitex/uploadClient/scraper"
	"NekoImageWorkflowKitex/uploadClient/storage"
	"context"
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
	Scrapers          *[]scraper.ScraperInstance
	PreUploadBridge   *storage.PreUploadTransBridgeInstance
	UploadTransBridge *storage.UploadTransBridgeInstance
}

// OnInit load client self config and before data, then init Scrapers
func (ci *ClientInstance) OnInit() error {
	// init
	logrus.Debug("ClientInstance OnInit start")
	ci.ClientInfo = storage.GetConfig()
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
	_, err := cli.HandleFilePreUpload(ctx, req)
	if err != nil {
		return err
	}
	return nil
}

// HandleFilePostUpload report post upload data
func (ci *ClientInstance) HandleFilePostUpload(ctx context.Context, cli fileuploadservice.Client) error {
	// TODO: make it really work
	logrus.Debug("ClientInstance PostUpload start")
	req := &uploadClient.FilePostRequest{}
	_, err := cli.HandleFilePostUpload(ctx, req)
	if err != nil {
		return err
	}
	return nil
}

// OnStop write PreUploadBridge data and UploadTransBridge data to disk
func (ci *ClientInstance) OnStop() error {
	// TODO: make it really work, such as write flush channel data to disk
	logrus.Debug("ClientInstance OnStop start")
	return nil
}
