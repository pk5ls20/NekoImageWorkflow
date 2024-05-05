package model

import (
	"context"
	"github.com/pk5ls20/NekoImageWorkflow/uploadClient/kitex_gen/protoFile/fileuploadservice"
)

type ClientImpl interface {
	// OnInit call after logger init, load client self config, init database, write data to dataBridge
	OnInit() error
	// OnStart call after kitex's MustNewClient, aka after kitex client start and before kitex transport start
	OnStart() error
	// HandleFilePreUpload report pre upload data
	HandleFilePreUpload(ctx context.Context, cli fileuploadservice.Client) error
	// HandleFilePostUpload report post upload data
	HandleFilePostUpload(ctx context.Context, cli fileuploadservice.Client) error
	// OnStop call on program exit
	// TODO: maybe we can add database write logic here
	OnStop() error
}
