package impl

import (
	"context"
	"github.com/pk5ls20/NekoImageWorkflow/common/log"
	clientImplModel "github.com/pk5ls20/NekoImageWorkflow/uploadClient/client/model"
	uploadClient "github.com/pk5ls20/NekoImageWorkflow/uploadClient/kitex_gen/protoFile"
	"github.com/pk5ls20/NekoImageWorkflow/uploadClient/kitex_gen/protoFile/fileuploadservice"
	"github.com/pk5ls20/NekoImageWorkflow/uploadClient/scraper/impl"
	"github.com/pk5ls20/NekoImageWorkflow/uploadClient/storage/bridge"
	"github.com/pk5ls20/NekoImageWorkflow/uploadClient/storage/config"
	"github.com/pk5ls20/NekoImageWorkflow/uploadClient/storage/sqlite"
	"github.com/sirupsen/logrus"
)

type ClientInstance struct {
	clientImplModel.ClientImpl
	ClientInfo        *config.ClientConfig
	Scrapers          []impl.ScraperInstance
	PreUploadBridge   *bridge.PreUploadTransBridgeInstance
	UploadTransBridge *bridge.PostUploadTransBridgeInstance
}

// OnInit load client self config and before data, then init Scrapers
func (ci *ClientInstance) OnInit() error {
	// init
	sqlite.InitSqlite()
	logrus.Debug("ClientInstance OnInit start")
	ci.ClientInfo = config.GetConfig()
	return nil
}

// OnStart currently do nothing
func (ci *ClientInstance) OnStart() error {
	// TODO: make it really work
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
	// TODO: dynamic write sqlite data
	logrus.Debug("ClientInstance OnStop start")
	sqlite.CloseSqlite()
	return nil
}
