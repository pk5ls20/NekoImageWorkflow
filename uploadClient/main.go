package main

import (
	"NekoImageWorkflowKitex/uploadClient/impl"
	"NekoImageWorkflowKitex/uploadClient/kitex_gen/uploadClient/fileuploadservice"
	FinalLog "NekoImageWorkflowKitex/uploadClient/log"
	"NekoImageWorkflowKitex/uploadClient/model"
	"NekoImageWorkflowKitex/uploadClient/scraper"
	"NekoImageWorkflowKitex/uploadClient/scraper/local"
	"NekoImageWorkflowKitex/uploadClient/storage"
	"context"
	kitexClient "github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/pkg/klog"
	kitexTransport "github.com/cloudwego/kitex/transport"
	kitexlogrus "github.com/kitex-contrib/obs-opentelemetry/logging/logrus"
	"github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	ctx := context.Background()
	// TODO: make it graceful?
	// init logrus
	logrus.SetOutput(os.Stdout)
	logrus.SetLevel(logrus.DebugLevel)
	logrus.SetReportCaller(true)
	logrus.SetFormatter(&FinalLog.CustomFormatter{})
	// init klog
	logger := kitexlogrus.NewLogger()
	logger.SetOutput(os.Stdout)
	logger.SetLevel(klog.LevelDebug)
	logger.Logger().SetReportCaller(true)
	logger.Logger().SetFormatter(&FinalLog.CustomFormatter{})
	klog.SetLogger(logger)
	// init clientImpl
	clientImpl := impl.ClientInstance{
		ClientInfo:        &model.ClientConfig{},
		Scrapers:          new([]scraper.ScraperInstance),
		PreUploadBridge:   storage.GetPreUploadTransBridgeInstance(),
		UploadTransBridge: storage.GetUploadTransBridgeInstance(),
	}
	// OnInit will load config
	if err := clientImpl.OnInit(); err != nil {
		logrus.Error("OnInit error:", err)
	}
	// init kitex client
	// TODO: make it really work
	// TODO: do we really need etcd?
	client := fileuploadservice.MustNewClient(
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
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigChan
		logrus.Warn("Received shutdown signal")
		os.Exit(0)
	}()
	// start scrapers
	var localScraper = &local.LocalScraperInstance{}
	*clientImpl.Scrapers = append(*clientImpl.Scrapers, localScraper)
	go func() {
		for _, scraperInstance := range *clientImpl.Scrapers {
			go func() {
				err := scraperInstance.PrepareData()
				if err != nil {
					logrus.Error("PrepareData error:", err)
				}
			}()
			anotherScraperInstance := scraperInstance
			go func() {
				err := anotherScraperInstance.ProcessData()
				if err != nil {
					logrus.Error("ProcessData error:", err)
				}
			}()
		}
	}()
	// start client upload
	for {
		if err := clientImpl.HandleFilePreUpload(ctx, client); err != nil {
			logrus.Error("PreUpload error:", err)
		}
		if err := clientImpl.HandleFilePostUpload(ctx, client); err != nil {
			logrus.Error("PostUpload error:", err)
		}
		time.Sleep(time.Second * time.Duration(clientImpl.ClientInfo.PostUploadPeriod))
	}
}
