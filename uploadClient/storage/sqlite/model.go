package sqlite

import (
	"github.com/pk5ls20/NekoImageWorkflow/uploadClient/model"
)

const (
	UUIDTag keyTag = iota
	DbDataTag
)

// dbDataModel represents an interface that implements any model.BaseDataModel.
type dbDataModel interface {
	model.BaseDataModel
}

// dbData is a data structure that acts as an intermediary between the actual dbDataStoredModel stored in the sqlite
// and the arbitrary implementation of model.BaseDataModel used by the client.
type dbData struct {
	Tag  keyTag
	Data dbDataModel
}

// dbDataStoredModel The actual structure stored in sqlite
type dbDataStoredModel struct {
	Tag  keyTag
	Data []byte
}
