package main

import (
	"context"
	kitexServer "github.com/cloudwego/kitex/server"
	_ "github.com/pk5ls20/NekoImageWorkflow/common/log"
	kitexClientTransform "github.com/pk5ls20/NekoImageWorkflow/kitex_gen/proto/clientTransform"
	kitexUploadService "github.com/pk5ls20/NekoImageWorkflow/kitex_gen/proto/clientTransform/fileuploadservice"
	"github.com/sirupsen/logrus"
)

type ServiceImpl struct{}

func (s *ServiceImpl) HandleFilePreUpload(ctx context.Context, req *kitexClientTransform.FilePreRequest) (resp *kitexClientTransform.FilePreResponse, err error) {
	logrus.Info("ChatA called, req: ", req)
	resp = new(kitexClientTransform.FilePreResponse)
	resp.Message = "hello" + req.ClientInfo.ClientName
	return
}

func (s *ServiceImpl) HandleFilePostUpload(ctx context.Context, req *kitexClientTransform.FilePostRequest) (resp *kitexClientTransform.FilePostResponse, err error) {
	logrus.Info("ChatB called, req: ", req)
	resp = new(kitexClientTransform.FilePostResponse)
	resp.Message = "hello " + req.ClientInfo.ClientName
	return
}

func main() {
	svr := kitexServer.NewServer()
	if err := kitexUploadService.RegisterService(svr, &ServiceImpl{}); err != nil {
		logrus.Fatal(err)
	}
	if err := svr.Run(); err != nil {
		logrus.Warning("server stopped with error:", err)
	} else {
		logrus.Info("server stopped")
	}
}
