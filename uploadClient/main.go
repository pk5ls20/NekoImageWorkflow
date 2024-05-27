package main

import (
	"context"
	"fmt"
	kitexClient "github.com/cloudwego/kitex/client"
	kitexTransport "github.com/cloudwego/kitex/transport"
	_ "github.com/pk5ls20/NekoImageWorkflow/common/log"
	kitexUploadService "github.com/pk5ls20/NekoImageWorkflow/kitex_gen/proto/clientTransform/fileuploadservice"
	clientImpl "github.com/pk5ls20/NekoImageWorkflow/uploadClient/client/impl"
	"github.com/pk5ls20/NekoImageWorkflow/uploadClient/pprof"
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
	// 0. paste flag
	parseFlags()
	// 1. load pprof
	if *pprofEnable {
		go pprof.RunPprof(*pprofAddr)
	}
	// 2. init clientImpl
	client = clientImpl.Client{}
	if err := client.OnInit(); err != nil {
		logrus.Fatal("OnInit error:", err)
	}
	// 3. load client uuid
	clientUUID, _ := storageSqlite.LoadClientUUID()
	logrus.Debug("Client uuid: ", clientUUID.String())
	// 4. init kitex client
	clientName := fmt.Sprintf("uploadClient-%s", clientUUID.String())
	kitexClientImpl := kitexUploadService.MustNewClient(
		clientName,
		kitexClient.WithTransportProtocol(kitexTransport.GRPC),
		kitexClient.WithHostPorts(client.ClientConfig.KitexServerAddress),
	)
	if err := client.OnStart(clientName, clientUUID); err != nil {
		logrus.Error("OnStart error:", err)
	}
	// 5. register signal handle, trigger client.OnStop when receive SIGINT
	RegisterSignalHandle()
	// 6. handle runtime
	if *runTime > 0 {
		logrus.Warning("Will send SIGINT after ", *runTime, " seconds")
		time.AfterFunc(time.Duration(*runTime)*time.Second, func() {
			signalChan <- syscall.SIGINT
			logrus.Warn("Sending SIGINT after ", *runTime, " seconds")
		})
	}
	// 7. start client upload
	ctx := context.Background()
	wg.Add(2)
	go func() {
		for {
			if err := client.HandleFilePreUpload(ctx, kitexClientImpl); err != nil {
				logrus.Error("PreUpload error:", err)
				time.Sleep(time.Duration(int64(client.ClientConfig.UploadFailWaitSecond*1000)) * time.Millisecond)
			}
			time.Sleep(time.Duration(int64(client.ClientConfig.UploadWaitSecond*1000)) * time.Millisecond)
		}
	}()
	go func() {
		for {
			if err := client.HandleFilePostUpload(ctx, kitexClientImpl); err != nil {
				logrus.Error("PostUpload error:", err)
				time.Sleep(time.Duration(int64(client.ClientConfig.UploadFailWaitSecond*1000)) * time.Millisecond)
			}
			time.Sleep(time.Duration(int64(client.ClientConfig.UploadWaitSecond*1000)) * time.Millisecond)
		}
	}()
	// 8. block main goroutine
	wg.Wait()
}
