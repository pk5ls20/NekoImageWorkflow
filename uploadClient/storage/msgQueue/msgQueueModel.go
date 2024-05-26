package msgQueue

import (
	"context"
	commonModel "github.com/pk5ls20/NekoImageWorkflow/common/model"
	clientModel "github.com/pk5ls20/NekoImageWorkflow/uploadClient/client/model"
	"sync"
)

type msgQueueType int

const (
	ActivateQueue msgQueueType = iota
	DeadQueue
)

type msgQueueData interface {
	// Commit removes the message from the msgQueue
	Commit() error
	// GoDead removes the message from the msgQueue and stores it in the dead msgQueue
	GoDead() error
}

// MsgMetaData identifies a message msgQueue
// ScraperID is the ID of the scraper
// MsgGroupID is the ID of the message group
// For example, in LocalScraper, ScraperID is the timestamp when collection started
// (because collection is non-repetitive for each start)
// but in APIScraper, ScraperID is simply the APIAddress itself
type MsgMetaData struct {
	UploadType commonModel.UploadType
	MsgMetaID
}

type MsgMetaID struct {
	ScraperType commonModel.ScraperType
	ScraperID   string
	MsgGroupID  string
}

// MsgQueueData represents a message in the msgQueue
type MsgQueueData struct {
	msgQueueData
	MsgMetaData
	FileMetaData *clientModel.AnyFileMetaDataModel
}

type msgPureData struct {
	MsgMetaID      MsgMetaID
	PreUploadModel clientModel.PreUploadFileMetaDataModel
	UploadModel    clientModel.UploadFileMetaDataModel
}

type messageQueue interface {
	AddElement(data *MsgQueueData) error
	AddElements(dataSlice []*MsgQueueData) error
	ListenUploadType(ctx context.Context, uploadType commonModel.UploadType) (<-chan *MsgQueueData, error)
	ListenMsgMetaData(data MsgMetaData) (<-chan *MsgQueueData, error)
	PopData(data MsgMetaData) ([]*MsgQueueData, error)
	PopAll(queueType msgQueueType) ([]*MsgQueueData, error)
}

type MessageQueue struct {
	messageQueue
	initialize bool
	lock       sync.Mutex
	// activateQueue stores the active message
	activateQueue sync.Map
	// deadQueue stores the dead message
	deadQueue sync.Map
	// uploadTypeElementChan is used to notify the ListenScpID that a new message is added
	uploadTypeElementChan sync.Map
	// uploadTypeListenChan is used to store the channel for whole ScraperID
	uploadTypeListenChan sync.Map
}
