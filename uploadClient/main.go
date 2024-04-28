package main

import (
	"context"
	kitexClient "github.com/cloudwego/kitex/client"
	kitexTransport "github.com/cloudwego/kitex/transport"
	_ "github.com/pk5ls20/NekoImageWorkflow/common/log"
	"github.com/pk5ls20/NekoImageWorkflow/uploadClient/client"
	"github.com/pk5ls20/NekoImageWorkflow/uploadClient/kitex_gen/protoFile/fileuploadservice"
	"github.com/pk5ls20/NekoImageWorkflow/uploadClient/model"
	"github.com/pk5ls20/NekoImageWorkflow/uploadClient/scraper"
	ScraperLifeCycle "github.com/pk5ls20/NekoImageWorkflow/uploadClient/scraper/lifecycle"
	"github.com/pk5ls20/NekoImageWorkflow/uploadClient/storage/bridge"
	"github.com/pk5ls20/NekoImageWorkflow/uploadClient/storage/sqlite"
	"github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var clientImpl client.ClientInstance

func RegisterSignalHandle() {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigChan
		logrus.Warn("Received shutdown signal")
		sqlite.CloseSqlite()
		os.Exit(0)
	}()
}

func init() {
	// 1. init db
	sqlite.InitSqlite()
	// 2. handle signal
	RegisterSignalHandle()
	// 3. init clientImpl
	clientImpl = client.ClientInstance{
		ClientInfo:        &model.ClientConfig{},
		Scrapers:          make([]scraper.ScraperInstance, 0),
		PreUploadBridge:   bridge.GetPreUploadTransBridgeInstance(),
		UploadTransBridge: bridge.GetUploadTransBridgeInstance(),
	}
	if err := clientImpl.OnInit(); err != nil {
		logrus.Fatal("OnInit error:", err)
	}
	//uuid, _ := sqlite.LoadClientUUID()
	//logrus.Debug("Client uuid: ", uuid.String())
}

func main() {
	// init kitex client
	// TODO: make it really work
	// TODO: do we really need etcd?
	kitexClientImpl := fileuploadservice.MustNewClient(
		clientImpl.ClientInfo.DestServiceName,
		kitexClient.WithTransportProtocol(kitexTransport.GRPC),
		kitexClient.WithHostPorts("127.0.0.1:8888"),
	)
	if err := clientImpl.OnStart(); err != nil {
		logrus.Error("OnStart error:", err)
	}
	defer func() {
		if err := clientImpl.OnStop(); err != nil {
			logrus.Error("OnStop error:", err)
		}
	}()
	// TODO:
	//time.AfterFunc(30*time.Second, func() {
	//	logrus.Warn("Sending SIGINT after 300 seconds")
	//	sigChan <- syscall.SIGINT
	//})
	// start scrapers
	clientImpl.Scrapers = ScraperLifeCycle.RegisterScraper(clientImpl.ClientInfo.ScraperInstance)
	go ScraperLifeCycle.StartScraper(clientImpl.Scrapers)
	// start client upload
	ctx := context.Background()
	for {
		if err := clientImpl.HandleFilePreUpload(ctx, kitexClientImpl); err != nil {
			logrus.Error("PreUpload error:", err)
		}
		if err := clientImpl.HandleFilePostUpload(ctx, kitexClientImpl); err != nil {
			logrus.Error("PostUpload error:", err)
		}
		time.Sleep(time.Second * time.Duration(clientImpl.ClientInfo.PostUploadPeriod))
	}
}
