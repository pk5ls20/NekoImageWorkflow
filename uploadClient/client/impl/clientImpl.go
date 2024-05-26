package impl

import (
	"context"
	"github.com/google/uuid"
	commonLog "github.com/pk5ls20/NekoImageWorkflow/common/log"
	commonModel "github.com/pk5ls20/NekoImageWorkflow/common/model"
	kitexUploadClient "github.com/pk5ls20/NekoImageWorkflow/kitex_gen/proto/clientTransform"
	kitexUploadService "github.com/pk5ls20/NekoImageWorkflow/kitex_gen/proto/clientTransform/fileuploadservice"
	clientModel "github.com/pk5ls20/NekoImageWorkflow/uploadClient/client/model"
	scraperLifeCycle "github.com/pk5ls20/NekoImageWorkflow/uploadClient/scraper/lifecycle"
	scraperModel "github.com/pk5ls20/NekoImageWorkflow/uploadClient/scraper/model"
	storageConfig "github.com/pk5ls20/NekoImageWorkflow/uploadClient/storage/config"
	storageQueue "github.com/pk5ls20/NekoImageWorkflow/uploadClient/storage/msgQueue"
	storageSqlite "github.com/pk5ls20/NekoImageWorkflow/uploadClient/storage/sqlite"
	"github.com/sirupsen/logrus"
)

type client interface {
	// OnInit call after logger init, load client self config, init database, write database data to fileQueue
	OnInit() error
	// OnStart call after kitex's MustNewClient, aka after kitex client start and before kitex transport start
	OnStart(clientName string, clientUUID uuid.UUID) error
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
	KitexClientInfo *kitexUploadClient.ClientInfo
	ClientConfig    *storageConfig.ClientConfig
	MsgQueue        *storageQueue.MessageQueue
	Scrapers        []scraperModel.Scraper
	ScraperChanMap  scraperModel.ScraperChanMap
}

// OnInit load client self config and before data, then init Scrapers
func (ci *Client) OnInit() error {
	// init
	storageSqlite.InitSqlite()
	logrus.Debug("Client OnInit start")
	ci.MsgQueue = storageQueue.NewMessageQueue()
	ci.ClientConfig = storageConfig.GetConfig()
	return nil
}

// OnStart currently init KitexClientInfo and Scrapers
func (ci *Client) OnStart(clientName string, clientUUID uuid.UUID) error {
	logrus.Debug("Client OnStart start")
	ci.KitexClientInfo = &kitexUploadClient.ClientInfo{
		ClientUUID: clientUUID.String(),
		ClientName: clientName,
	}
	ci.Scrapers = scraperLifeCycle.RegisterScraper(ci.ClientConfig.Scraper)
	go scraperLifeCycle.StartScraper(ci.Scrapers)
	return nil
}

// HandleFilePreUpload report pre upload data
// TODO: Need to store filedata that failed to upload
func (ci *Client) HandleFilePreUpload(ctx context.Context, cli kitexUploadService.Client) error {
	logrus.Debug("Client PreUpload start")
	ch, err := ci.MsgQueue.ListenUploadType(ctx, commonModel.PreUploadType)
	if err != nil {
		return commonLog.ErrorWrap(err)
	}
	// TODO: make it really work
	select {
	case <-ctx.Done():
		return nil
	case msg := <-ch:
		logrus.Debug("Client PreUpload msg: ", msg)
		m := clientModel.NewPreUploadFileDataFromMeta(msg.FileMetaData.PreUploadFileMetaDataModel)
		mSlice := []*kitexUploadClient.PreUploadFileData{{
			ScraperType:  kitexUploadClient.ScraperType(commonModel.PasteScraperTypeToInt(m.ScraperType)),
			ResourceUUID: m.ResourceUUID.String(),
			ResourceUri:  m.ResourceUri,
		}}
		req := &kitexUploadClient.FilePreRequest{
			ClientInfo: ci.KitexClientInfo,
			Data:       mSlice,
		}
		resp, _err := cli.HandleFilePreUpload(ctx, req)
		if _err != nil {
			return commonLog.ErrorWrap(err)
		}
		logrus.Debug("Client PreUpload resp: ", resp)
		// TODO: handle HandleFilePreUpload response
		if err = m.FinishUpload(); err != nil {
			return commonLog.ErrorWrap(err)
		}
		if err = msg.Commit(); err != nil {
			return commonLog.ErrorWrap(err)
		}
		return nil
	}
}

// HandleFilePostUpload report post upload data
// TODO: Need to store filedata that failed to upload
// TODO: make it really work
func (ci *Client) HandleFilePostUpload(ctx context.Context, cli kitexUploadService.Client) error {
	logrus.Debug("Client PostUpload start")
	ch, err := ci.MsgQueue.ListenUploadType(ctx, commonModel.PostUploadType)
	if err != nil {
		return commonLog.ErrorWrap(err)
	}
	// TODO: make it really work
	select {
	case <-ctx.Done():
		return nil
	case msg := <-ch:
		logrus.Debug("Client PostUpload msg: ", msg)
		m := clientModel.NewUploadFileDataFromMeta(msg.FileMetaData.UploadFileMetaDataModel)
		fileData, _err := m.GetFileContent()
		if _err != nil {
			return commonLog.ErrorWrap(_err)
		}
		mSlice := []*kitexUploadClient.UploadFileData{{
			ScraperType: kitexUploadClient.ScraperType(commonModel.PasteScraperTypeToInt(m.ScraperType)),
			FileUUID:    m.FileUUID.String(),
			FileContent: fileData,
		}}
		req := &kitexUploadClient.FilePostRequest{
			ClientInfo: ci.KitexClientInfo,
			Data:       mSlice,
		}
		resp, _err := cli.HandleFilePostUpload(ctx, req)
		if _err != nil {
			return commonLog.ErrorWrap(err)
		}
		// TODO: handle HandleFilePreUpload response
		logrus.Debug("Client PostUpload resp: ", resp)
		if err = m.FinishUpload(); err != nil {
			return commonLog.ErrorWrap(err)
		}
		if err = msg.Commit(); err != nil {
			return commonLog.ErrorWrap(err)
		}
		return nil
	}
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
