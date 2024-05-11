package main

import (
	"context"
	kitexClient "github.com/cloudwego/kitex/client"
	kitexTransport "github.com/cloudwego/kitex/transport"
	_ "github.com/pk5ls20/NekoImageWorkflow/common/log"
	clientImpl "github.com/pk5ls20/NekoImageWorkflow/uploadClient/client/impl"
	"github.com/pk5ls20/NekoImageWorkflow/uploadClient/kitex_gen/protoFile/fileuploadservice"
	"github.com/pk5ls20/NekoImageWorkflow/uploadClient/scraper/model"
	clientModel "github.com/pk5ls20/NekoImageWorkflow/uploadClient/storage/config"
	"github.com/pk5ls20/NekoImageWorkflow/uploadClient/storage/queue"
	"github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var clientEntry clientImpl.ClientInstance
var signalChan chan os.Signal

func RegisterSignalHandle() {
	signalChan = make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-signalChan
		logrus.Warn("Received shutdown signal")
		if err := clientEntry.OnStop(); err != nil {
			logrus.Error("OnStop error:", err)
		}
		os.Exit(0)
	}()
}

func main() {
	// 1. init clientImpl
	clientEntry = clientImpl.ClientInstance{
		ClientInfo:     &clientModel.ClientConfig{},
		Scrapers:       make([]model.ScraperInstance, 0),
		PreUploadQueue: queue.GetPreUploadQueueInstance(),
		UploadQueue:    queue.GetUploadQueueInstance(),
	}
	if err := clientEntry.OnInit(); err != nil {
		logrus.Fatal("OnInit error:", err)
	}
	// TODO: 2. load client uuid
	//uuid, _ := sqlite.LoadClientUUID()
	//logrus.Debug("Client uuid: ", uuid.String())
	// TODO: 3. init kitex client
	kitexClientImpl := fileuploadservice.MustNewClient(
		clientEntry.ClientInfo.DestServiceName,
		kitexClient.WithTransportProtocol(kitexTransport.GRPC),
		kitexClient.WithHostPorts("127.0.0.1:8888"),
	)
	if err := clientEntry.OnStart(); err != nil {
		logrus.Error("OnStart error:", err)
	}
	// TODO: debug after 30s
	//time.AfterFunc(30*time.Second, func() {
	//	signalChan <- syscall.SIGINT
	//	logrus.Warn("Sending SIGINT after 300 seconds")
	//})
	// 4. register signal handle, trigger clientEntry.OnStop when receive SIGINT
	RegisterSignalHandle()
	// 5. start client upload
	ctx := context.Background()
	for {
		if err := clientEntry.HandleFilePreUpload(ctx, kitexClientImpl); err != nil {
			logrus.Error("PreUpload error:", err)
		}
		if err := clientEntry.HandleFilePostUpload(ctx, kitexClientImpl); err != nil {
			logrus.Error("PostUpload error:", err)
		}
		time.Sleep(time.Second * time.Duration(clientEntry.ClientInfo.PostUploadPeriod))
	}
}
