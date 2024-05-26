package msgQueue

import (
	"context"
	"errors"
	mapset "github.com/deckarep/golang-set/v2"
	commonLog "github.com/pk5ls20/NekoImageWorkflow/common/log"
	commonModel "github.com/pk5ls20/NekoImageWorkflow/common/model"
	"github.com/sirupsen/logrus"
	"sync"
)

// TODO: more precise lock
var (
	messageQueueInstance = &MessageQueue{}
	once                 sync.Once
	commitLock           sync.Mutex
	goDeadLock           sync.Mutex
)

func NewMessageQueue() *MessageQueue {
	once.Do(func() {
		messageQueueInstance = &MessageQueue{
			initialize: true,
		}
	})
	return messageQueueInstance
}

func (mq *MessageQueue) AddElement(data *MsgQueueData) error {
	if !mq.initialize {
		return commonLog.ErrorWrap(errors.New("MessageQueue not initialized"))
	}
	// Add to the activateQueue
	activateQueueEntry, _ := mq.activateQueue.LoadOrStore(data.MsgMetaData, mapset.NewSet[*MsgQueueData]())
	activateQueueEntry.(msgQueueDataSet).Add(data)
	// Add MsgMetaData to the listener
	uploadTypeEntry, _ := mq.uploadTypeElementChan.LoadOrStore(data.UploadType, make(chan *MsgQueueData, 10000))
	go func() {
		uploadTypeEntry.(chan *MsgQueueData) <- data
	}()
	return nil
}

// AddElements TODO: Use builtin methods
func (mq *MessageQueue) AddElements(dataSlice []*MsgQueueData) error {
	if !mq.initialize {
		return commonLog.ErrorWrap(errors.New("MessageQueue not initialized"))
	}
	for _, data := range dataSlice {
		if err := mq.AddElement(data); err != nil {
			return err
		}
	}
	return nil
}

func (mq *MessageQueue) ListenUploadType(ctx context.Context, uploadType commonModel.UploadType) (<-chan *MsgQueueData, error) {
	mq.lock.Lock()
	defer mq.lock.Unlock()
	if !mq.initialize {
		return nil, commonLog.ErrorWrap(errors.New("MessageQueue not initialized"))
	}
	// check if already exists
	if existChan, ok := mq.uploadTypeListenChan.Load(uploadType); ok {
		logrus.Debugf("uploadTypeListenChan for uploadType %s already exists...", uploadType)
		return existChan.(chan *MsgQueueData), nil
	}
	mainChan := make(chan *MsgQueueData)
	// handle ctx cancel
	go func() {
		<-ctx.Done()
		close(mainChan)
		mq.uploadTypeListenChan.Delete(uploadType)
		logrus.Debugf("Context Done for ScraperType %s", uploadType)
	}()
	// listen further uploadTypeElementChan - UploadType
	go func(uploadType commonModel.UploadType) {
		logrus.Debugf("Start to listen further uploadTypeElementChan for uploadType %s", uploadType)
		entry, exist := mq.uploadTypeElementChan.LoadOrStore(uploadType, make(chan *MsgQueueData, 10000))
		logrus.Debugf("uploadTypeElementChan for uploadType %s exists: %t", uploadType, exist)
		mp := make(map[MsgQueueData]struct{})
		var itm *MsgQueueData
		var ok bool
		for {
			select {
			case <-ctx.Done():
				logrus.Debug("Context cancelled, exiting goroutine")
				return
			case itm, ok = <-entry.(chan *MsgQueueData):
				if !ok {
					logrus.Debug("Channel closed, performing cleanup")
					return
				}
				if _, exists := mp[*itm]; exists {
					logrus.Debugf("Element already exists in the map...")
					continue
				}
				mp[*itm] = struct{}{}
				mainChan <- itm
			}
		}
	}(uploadType)
	mq.uploadTypeListenChan.Store(uploadType, mainChan)
	return mainChan, nil
}

// ListenMsgMetaData
// NOTE: Due to the design of mapset itself (encapsulated as a map underneath),
// it can only iterate through the content of the set that existed prior to the call.
func (mq *MessageQueue) ListenMsgMetaData(data MsgMetaData) (<-chan *MsgQueueData, error) {
	if !mq.initialize {
		return nil, commonLog.ErrorWrap(errors.New("MessageQueue not initialized"))
	}
	entry, ok := mq.activateQueue.Load(data)
	if !ok {
		return nil, commonLog.ErrorWrap(errors.New("no such message group"))
	}
	return entry.(msgQueueDataSet).Iter(), nil
}

func (mq *MessageQueue) PopData(data MsgMetaData) (*MsgQueueData, error) {
	mq.lock.Lock()
	defer mq.lock.Unlock()
	if !mq.initialize {
		return &MsgQueueData{}, commonLog.ErrorWrap(errors.New("MessageQueue not initialized"))
	}
	entry, ok := mq.activateQueue.Load(data)
	if !ok {
		return &MsgQueueData{}, commonLog.ErrorWrap(errors.New("no such message group"))
	}
	element, ok := entry.(msgQueueDataSet).Pop()
	if !ok {
		return &MsgQueueData{}, commonLog.ErrorWrap(errors.New("no element in the message group"))
	}
	mq.activateQueue.Delete(data)
	return element, nil
}

// PopAll TODO: lock?
func (mq *MessageQueue) PopAll(queueType msgQueueType) ([]*MsgQueueData, error) {
	mq.lock.Lock()
	defer mq.lock.Unlock()
	if !mq.initialize {
		return nil, commonLog.ErrorWrap(errors.New("MessageQueue not initialized"))
	}
	var elements []*MsgQueueData
	var currentQueueMap *sync.Map
	switch queueType {
	case ActivateQueue:
		currentQueueMap = &mq.activateQueue
	case DeadQueue:
		currentQueueMap = &mq.deadQueue
	default:
		return nil, commonLog.ErrorWrap(errors.New("invalid msgQueue type"))
	}
	currentQueueMap.Range(func(k, v interface{}) bool {
		elements = append(elements, v.(msgQueueDataSet).ToSlice()...)
		v.(msgQueueDataSet).Clear()
		return true
	})
	return elements, nil
}

func (mq *MsgQueueData) Commit() error {
	commitLock.Lock()
	defer commitLock.Unlock()
	if !messageQueueInstance.initialize {
		return commonLog.ErrorWrap(errors.New("MessageQueue not initialized"))
	}
	q, ok := messageQueueInstance.activateQueue.Load(mq.MsgMetaData)
	if !ok {
		return commonLog.ErrorWrap(errors.New("MessageQueue not found"))
	}
	q.(msgQueueDataSet).Remove(mq)
	return nil
}

func (mq *MsgQueueData) GoDead() error {
	goDeadLock.Lock()
	defer goDeadLock.Unlock()
	if !messageQueueInstance.initialize {
		return commonLog.ErrorWrap(errors.New("MessageQueue not initialized"))
	}
	if err := mq.Commit(); err != nil {
		return commonLog.ErrorWrap(err)
	}
	entry, _ := messageQueueInstance.deadQueue.LoadOrStore(mq.MsgMetaData, mapset.NewSet[*MsgQueueData]())
	entry.(msgQueueDataSet).Add(mq)
	return nil
}