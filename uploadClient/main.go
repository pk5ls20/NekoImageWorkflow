package main

import (
	"context"
	"fmt"
	kitexClient "github.com/cloudwego/kitex/client"
	kitexTransport "github.com/cloudwego/kitex/transport"
	_ "github.com/pk5ls20/NekoImageWorkflow/common/log"
	kitexUploadService "github.com/pk5ls20/NekoImageWorkflow/kitex_gen/proto/clientTransform/fileuploadservice"
	clientImpl "github.com/pk5ls20/NekoImageWorkflow/uploadClient/client/impl"
	storageSqlite "github.com/pk5ls20/NekoImageWorkflow/uploadClient/storage/sqlite"
	"github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

var (
	client     clientImpl.Client
	signalChan chan os.Signal
	wg         sync.WaitGroup
)

func RegisterSignalHandle() {
	signalChan = make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-signalChan
		logrus.Warn("Received shutdown signal")
		if err := client.OnStop(); err != nil {
			logrus.Error("OnStop error:", err)
		}
		os.Exit(0)
	}()
}

func main() {
	// 1. init clientImpl
	client = clientImpl.Client{}
	if err := client.OnInit(); err != nil {
		logrus.Fatal("OnInit error:", err)
	}
	// TODO: 2. load client uuid
	clientUUID, _ := storageSqlite.LoadClientUUID()
	logrus.Debug("Client uuid: ", clientUUID.String())
	// TODO: 3. init kitex client
	clientName := fmt.Sprintf("uploadClient-%s", clientUUID.String())
	kitexClientImpl := kitexUploadService.MustNewClient(
		clientName,
		kitexClient.WithTransportProtocol(kitexTransport.GRPC),
		kitexClient.WithHostPorts("127.0.0.1:8888"),
	)
	if err := client.OnStart(clientName, clientUUID); err != nil {
		logrus.Error("OnStart error:", err)
	}
	// TODO: debug after 30s
	//time.AfterFunc(30*time.Second, func() {
	//	signalChan <- syscall.SIGINT
	//	logrus.Warn("Sending SIGINT after 300 seconds")
	//})
	// 4. register signal handle, trigger client.OnStop when receive SIGINT
	RegisterSignalHandle()
	// 5. start client upload
	ctx := context.Background()
	wg.Add(2)
	go func() {
		for {
			if err := client.HandleFilePreUpload(ctx, kitexClientImpl); err != nil {
				logrus.Error("PreUpload error:", err)
				time.Sleep(time.Duration(int64(client.ClientConfig.PostUploadPeriod*1000)) * time.Millisecond)
			}
		}
	}()
	go func() {
		for {
			if err := client.HandleFilePostUpload(ctx, kitexClientImpl); err != nil {
				logrus.Error("PostUpload error:", err)
				time.Sleep(time.Duration(int64(client.ClientConfig.PostUploadPeriod*1000)) * time.Millisecond)
			}
		}
	}()
	wg.Wait()
}
