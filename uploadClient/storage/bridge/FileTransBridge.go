package bridge

import (
	"context"
	"errors"
	"github.com/pk5ls20/NekoImageWorkflow/common/log"
	clientModel "github.com/pk5ls20/NekoImageWorkflow/uploadClient/client/model"
	"github.com/smallnest/chanx"
	"sync"
)

type fileTransBridge[T clientModel.BaseBridgeDataModel] interface {
	Length() int
	Insert(val []*T) error
	Pop(number int) ([]*T, error)
	PopAll() ([]*T, error)
	closeInputChannel() error
}

type baseFileTransBridgeInstance[T clientModel.BaseBridgeDataModel] struct {
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
	baseFileTransBridgeInstance[clientModel.ScraperPreUploadFileDataModel]
}

type PostUploadTransBridgeInstance struct {
	baseFileTransBridgeInstance[clientModel.ScraperPostUploadFileDataModel]
}

type PreTransformTransBridgeInstance struct {
	baseFileTransBridgeInstance[clientModel.PreTransformDataModel]
}

type PostTransformTransBridgeInstance struct {
	baseFileTransBridgeInstance[clientModel.PostTransformDataModel]
}

var (
	scraperPreUploadInstance  *PreUploadTransBridgeInstance
	scraperPreUploadOnce      sync.Once
	scraperPostUploadInstance *PostUploadTransBridgeInstance
	scraperPostUploadOnce     sync.Once
	preTransformInstance      *PreTransformTransBridgeInstance
	preTransformOnce          sync.Once
	postTransformInstance     *PostTransformTransBridgeInstance
	postTransformOnce         sync.Once
)

const initCap = 100

func GetPreUploadTransBridgeInstance() *PreUploadTransBridgeInstance {
	scraperPreUploadOnce.Do(func() {
		scraperPreUploadInstance = &PreUploadTransBridgeInstance{
			baseFileTransBridgeInstance: baseFileTransBridgeInstance[clientModel.ScraperPreUploadFileDataModel]{
				channel: chanx.NewUnboundedChan[*clientModel.ScraperPreUploadFileDataModel](context.Background(), initCap),
			},
		}
	})
	return scraperPreUploadInstance
}

func GetUploadTransBridgeInstance() *PostUploadTransBridgeInstance {
	scraperPostUploadOnce.Do(func() {
		scraperPostUploadInstance = &PostUploadTransBridgeInstance{
			baseFileTransBridgeInstance: baseFileTransBridgeInstance[clientModel.ScraperPostUploadFileDataModel]{
				channel: chanx.NewUnboundedChan[*clientModel.ScraperPostUploadFileDataModel](context.Background(), initCap),
			},
		}
	})
	return scraperPostUploadInstance
}

func GetPreTransformTransBridgeInstance() *PreTransformTransBridgeInstance {
	preTransformOnce.Do(func() {
		preTransformInstance = &PreTransformTransBridgeInstance{
			baseFileTransBridgeInstance: baseFileTransBridgeInstance[clientModel.PreTransformDataModel]{
				channel: chanx.NewUnboundedChan[*clientModel.PreTransformDataModel](context.Background(), initCap),
			},
		}
	})
	return preTransformInstance
}

func GetPostTransformTransBridgeInstance() *PostTransformTransBridgeInstance {
	postTransformOnce.Do(func() {
		postTransformInstance = &PostTransformTransBridgeInstance{
			baseFileTransBridgeInstance: baseFileTransBridgeInstance[clientModel.PostTransformDataModel]{
				channel: chanx.NewUnboundedChan[*clientModel.PostTransformDataModel](context.Background(), initCap),
			},
		}
	})
	return postTransformInstance
}
