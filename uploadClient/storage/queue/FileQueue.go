package queue

import (
	"context"
	"errors"
	"github.com/pk5ls20/NekoImageWorkflow/common/log"
	clientModel "github.com/pk5ls20/NekoImageWorkflow/uploadClient/client/model"
	"github.com/smallnest/chanx"
	"sync"
)

type fileQueue[T clientModel.BaseFileDataModel] interface {
	Length() int
	Insert(val []*T) error
	Pop(number int) ([]*T, error)
	PopAll() ([]*T, error)
	closeInputChannel() error
}

type baseFileQueueInstance[T clientModel.BaseFileDataModel] struct {
	channel *chanx.UnboundedChan[*T]
	fileQueue[T]
}

func (c *baseFileQueueInstance[T]) Length() int {
	return c.channel.Len()
}

func (c *baseFileQueueInstance[T]) Insert(val []*T) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = log.ErrorWrap(errors.New("failed to insert: channel is closed or channel is nil"))
		}
	}()
	for _, v := range val {
		c.channel.In <- v
	}
	return nil
}

func (c *baseFileQueueInstance[T]) Pop(number int) ([]*T, error) {
	if number < 0 {
		return nil, log.ErrorWrap(errors.New("pop number should be positive"))
	}
	tmp := make([]*T, 0, number)
	for i := 0; i < number; i++ {
		select {
		case v, ok := <-c.channel.Out:
			if !ok {
				if len(tmp) == 0 {
					return nil, log.ErrorWrap(errors.New("channel is closed and empty"))
				}
				return tmp, log.ErrorWrap(errors.New("channel closed during read"))
			}
			tmp = append(tmp, v)
		}
	}
	return tmp, nil
}

func (c *baseFileQueueInstance[T]) PopAll() ([]*T, error) {
	if err := c.closeInputChannel(); err != nil {
		return nil, log.ErrorWrap(err)
	}
	chanLen := c.channel.Len()
	pop, err := c.Pop(chanLen)
	if err != nil {
		return nil, log.ErrorWrap(err)
	}
	return pop, nil
}

func (c *baseFileQueueInstance[T]) closeInputChannel() (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = log.ErrorWrap(errors.New("failed to close: channel is closed or channel is nil"))
		}
	}()
	close(c.channel.In)
	return nil
}

type PreUploadQueueInstance struct {
	baseFileQueueInstance[clientModel.PreUploadFileDataModel]
}

type UploadQueueInstance struct {
	baseFileQueueInstance[clientModel.UploadFileDataModel]
}

var (
	preUploadQueueInstance *PreUploadQueueInstance
	preUploadQueueOnce     sync.Once
	uploadQueueInstance    *UploadQueueInstance
	uploadQueueOnce        sync.Once
)

const initCap = 100

func GetPreUploadQueueInstance() *PreUploadQueueInstance {
	preUploadQueueOnce.Do(func() {
		preUploadQueueInstance = &PreUploadQueueInstance{
			baseFileQueueInstance: baseFileQueueInstance[clientModel.PreUploadFileDataModel]{
				channel: chanx.NewUnboundedChan[*clientModel.PreUploadFileDataModel](context.Background(), initCap),
			},
		}
	})
	return preUploadQueueInstance
}

func GetUploadQueueInstance() *UploadQueueInstance {
	uploadQueueOnce.Do(func() {
		uploadQueueInstance = &UploadQueueInstance{
			baseFileQueueInstance: baseFileQueueInstance[clientModel.UploadFileDataModel]{
				channel: chanx.NewUnboundedChan[*clientModel.UploadFileDataModel](context.Background(), initCap),
			},
		}
	})
	return uploadQueueInstance
}
