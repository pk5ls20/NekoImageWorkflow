package sqlite

import (
	clientModel "github.com/pk5ls20/NekoImageWorkflow/uploadClient/client/model"
)

const (
	UUIDTag keyTag = iota
	QueueDataTag
)

// dbDataModel interface should contain all the structs expected to be stored in the database
// These structs **must be registered by gob**
type dbDataModel interface {
	clientModel.BaseFileDataModel
}

// dbData is a data structure that acts as an intermediary between the actual dbDataStoredModel stored in the sqlite
// and the arbitrary implementation of model.FileDataModel used by the client.
type dbData[T dbDataModel] struct {
	Tag  keyTag
	Data T
}

// dbDataStoredModel The actual structure stored in sqlite
type dbDataStoredModel struct {
	Tag  keyTag
	Data []byte
}
