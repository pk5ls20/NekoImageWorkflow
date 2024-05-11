package model

import (
	"context"
	"github.com/pk5ls20/NekoImageWorkflow/uploadClient/kitex_gen/protoFile/fileuploadservice"
)

type ClientImpl interface {
	// OnInit call after logger init, load client self config, init database, write database data to fileQueue
	OnInit() error
	// OnStart call after kitex's MustNewClient, aka after kitex client start and before kitex transport start
	OnStart() error
	// HandleFilePreUpload upload preUploadData to kitex server
	HandleFilePreUpload(ctx context.Context, cli fileuploadservice.Client) error
	// HandleFilePostUpload upload uploadData to kitex server
	HandleFilePostUpload(ctx context.Context, cli fileuploadservice.Client) error
	// OnStop call on program exit, write fileQueue data to database
	// TODO: maybe we can add database write logic here
	OnStop() error
}
