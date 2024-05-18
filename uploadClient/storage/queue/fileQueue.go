package queue

import (
	"context"
	"errors"
	commonLog "github.com/pk5ls20/NekoImageWorkflow/common/log"
	clientModel "github.com/pk5ls20/NekoImageWorkflow/uploadClient/client/model"
	"github.com/smallnest/chanx"
	"sync"
)

type fileQueue[T clientModel.BaseFileDataModel] interface {
	Length() int
	Insert(val []*T) error
	Pop() (*T, error)
	PopNum(number int) ([]*T, error)
	PopAll() ([]*T, error)
	closeInputChannel() error
}

type baseFileQueue[T clientModel.BaseFileDataModel] struct {
	channel *chanx.UnboundedChan[*T]
	fileQueue[T]
}

func (c *baseFileQueue[T]) Length() int {
	return c.channel.Len()
}

func (c *baseFileQueue[T]) Insert(val []*T) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = commonLog.ErrorWrap(errors.New("failed to insert: channel is closed or channel is nil"))
		}
	}()
	for _, v := range val {
		c.channel.In <- v
	}
	return nil
}

func (c *baseFileQueue[T]) Pop() (*T, error) {
	select {
	case v, ok := <-c.channel.Out:
		if !ok {
			return nil, commonLog.ErrorWrap(errors.New("channel is closed and empty"))
		}
		return v, nil
	}
}

func (c *baseFileQueue[T]) PopNum(number int) ([]*T, error) {
	if number < 0 {
		return nil, commonLog.ErrorWrap(errors.New("pop number should be positive"))
	}
	tmp := make([]*T, 0, number)
	for i := 0; i < number; i++ {
		data, err := c.Pop()
		if err != nil {
			return nil, commonLog.ErrorWrap(err)
		}
		tmp = append(tmp, data)
	}
	return tmp, nil
}

func (c *baseFileQueue[T]) PopAll() ([]*T, error) {
	if err := c.closeInputChannel(); err != nil {
		return nil, commonLog.ErrorWrap(err)
	}
	chanLen := c.channel.Len()
	pop, err := c.PopNum(chanLen)
	if err != nil {
		return nil, commonLog.ErrorWrap(err)
	}
	return pop, nil
}

func (c *baseFileQueue[T]) closeInputChannel() (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = commonLog.ErrorWrap(errors.New("failed to close: channel is closed or channel is nil"))
		}
	}()
	close(c.channel.In)
	return nil
}

type PreUploadQueue struct {
	baseFileQueue[clientModel.PreUploadFileDataModel]
}

type UploadQueue struct {
	baseFileQueue[clientModel.UploadFileDataModel]
}

var (
	preUploadQueueInstance *PreUploadQueue
	preUploadQueueOnce     sync.Once
	uploadQueueInstance    *UploadQueue
	uploadQueueOnce        sync.Once
)

const initCap = 100

func GetPreUploadQueue() *PreUploadQueue {
	preUploadQueueOnce.Do(func() {
		preUploadQueueInstance = &PreUploadQueue{
			baseFileQueue: baseFileQueue[clientModel.PreUploadFileDataModel]{
				channel: chanx.NewUnboundedChan[*clientModel.PreUploadFileDataModel](context.Background(), initCap),
			},
		}
	})
	return preUploadQueueInstance
}

func GetUploadQueue() *UploadQueue {
	uploadQueueOnce.Do(func() {
		uploadQueueInstance = &UploadQueue{
			baseFileQueue: baseFileQueue[clientModel.UploadFileDataModel]{
				channel: chanx.NewUnboundedChan[*clientModel.UploadFileDataModel](context.Background(), initCap),
			},
		}
	})
	return uploadQueueInstance
}
