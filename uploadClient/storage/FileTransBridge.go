package storage

import (
	"NekoImageWorkflowKitex/uploadClient/model"
	"context"
	"errors"
	"github.com/smallnest/chanx"
	"sync"
)

type FileTransBridge[T any] interface {
	Length() int
	Insert(val []T) error
	Pop(number int) []T
}

type BaseFileTransBridgeInstance[T any] struct {
	Channel *chanx.UnboundedChan[T]
	FileTransBridge[T]
}

func (c *BaseFileTransBridgeInstance[T]) Length() int {
	return c.Channel.Len()
}

func (c *BaseFileTransBridgeInstance[T]) Insert(val []T) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = errors.New("failed to insert: channel is closed")
		}
	}()
	for _, v := range val {
		c.Channel.In <- v
	}
	return nil
}

func (c *BaseFileTransBridgeInstance[T]) Pop(number int) ([]T, error) {
	if number <= 0 {
		return nil, errors.New("number must be positive")
	}
	tmp := make([]T, 0, number)
	for i := 0; i < number; i++ {
		select {
		case v, ok := <-c.Channel.Out:
			if !ok {
				if len(tmp) == 0 {
					return nil, errors.New("channel is closed and empty")
				}
				return tmp, errors.New("channel closed during read")
			}
			tmp = append(tmp, v)
		}
	}
	return tmp, nil
}

type PreUploadTransBridgeInstance struct {
	BaseFileTransBridgeInstance[model.PreUploadFileData]
}

type UploadTransBridgeInstance struct {
	BaseFileTransBridgeInstance[model.UploadFileData]
}

var preUploadInstance *PreUploadTransBridgeInstance
var preUploadOnce sync.Once
var uploadInstance *UploadTransBridgeInstance
var uploadOnce sync.Once

func GetPreUploadTransBridgeInstance() *PreUploadTransBridgeInstance {
	preUploadOnce.Do(func() {
		preUploadInstance = &PreUploadTransBridgeInstance{
			BaseFileTransBridgeInstance: BaseFileTransBridgeInstance[model.PreUploadFileData]{
				Channel: chanx.NewUnboundedChan[model.PreUploadFileData](context.Background(), 100),
			},
		}
	})
	return preUploadInstance
}

func GetUploadTransBridgeInstance() *UploadTransBridgeInstance {
	uploadOnce.Do(func() {
		uploadInstance = &UploadTransBridgeInstance{
			BaseFileTransBridgeInstance: BaseFileTransBridgeInstance[model.UploadFileData]{
				Channel: chanx.NewUnboundedChan[model.UploadFileData](context.Background(), 100),
			},
		}
	})
	return uploadInstance
}
