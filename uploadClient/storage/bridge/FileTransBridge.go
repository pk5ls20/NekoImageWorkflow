package bridge

import (
	"context"
	"errors"
	"github.com/pk5ls20/NekoImageWorkflow/common/log"
	"github.com/pk5ls20/NekoImageWorkflow/uploadClient/model"
	"github.com/smallnest/chanx"
	"sync"
)

type fileTransBridge[T model.AnyDataModel] interface {
	Length() int
	Insert(val []*T) error
	Pop(number int) ([]*T, error)
	PopAll() ([]*T, error)
	closeInputChannel() error
}

type baseFileTransBridgeInstance[T model.AnyDataModel] struct {
	channel *chanx.UnboundedChan[*T]
	fileTransBridge[T]
}

func (c *baseFileTransBridgeInstance[T]) Length() int {
	return c.channel.Len()
}

func (c *baseFileTransBridgeInstance[T]) Insert(val []*T) (err error) {
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

func (c *baseFileTransBridgeInstance[T]) Pop(number int) ([]*T, error) {
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

func (c *baseFileTransBridgeInstance[T]) PopAll() ([]*T, error) {
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

func (c *baseFileTransBridgeInstance[T]) closeInputChannel() (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = log.ErrorWrap(errors.New("failed to close: channel is closed or channel is nil"))
		}
	}()
	close(c.channel.In)
	return nil
}

type PreUploadTransBridgeInstance struct {
	baseFileTransBridgeInstance[model.ScraperPreUploadFileDataModel]
}

type UploadTransBridgeInstance struct {
	baseFileTransBridgeInstance[model.ScraperPostUploadFileDataModel]
}

var preUploadInstance *PreUploadTransBridgeInstance
var preUploadOnce sync.Once
var uploadInstance *UploadTransBridgeInstance
var uploadOnce sync.Once

const initCap = 100

func GetPreUploadTransBridgeInstance() *PreUploadTransBridgeInstance {
	preUploadOnce.Do(func() {
		preUploadInstance = &PreUploadTransBridgeInstance{
			baseFileTransBridgeInstance: baseFileTransBridgeInstance[model.ScraperPreUploadFileDataModel]{
				channel: chanx.NewUnboundedChan[*model.ScraperPreUploadFileDataModel](context.Background(), initCap),
			},
		}
	})
	return preUploadInstance
}

func GetUploadTransBridgeInstance() *UploadTransBridgeInstance {
	uploadOnce.Do(func() {
		uploadInstance = &UploadTransBridgeInstance{
			baseFileTransBridgeInstance: baseFileTransBridgeInstance[model.ScraperPostUploadFileDataModel]{
				channel: chanx.NewUnboundedChan[*model.ScraperPostUploadFileDataModel](context.Background(), initCap),
			},
		}
	})
	return uploadInstance
}
