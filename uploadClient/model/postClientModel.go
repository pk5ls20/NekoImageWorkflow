// Package model TODO: update these model to protobuf
package model

import (
	"github.com/google/uuid"
	"github.com/pk5ls20/NekoImageWorkflow/common/model"
)

// BaseDataModel Requires specific transport DataModel (just below) implementation
type BaseDataModel interface {
}

// ScraperPreUploadFileDataModel
// ResourceUUID Used to uniquely identify the resource
// ResourceUri Path to the local ClientImpl for the resource
type ScraperPreUploadFileDataModel struct {
	BaseDataModel
	ResourceUUID uuid.UUID
	ResourceUri  string
}

// ScraperPostUploadFileDataModel
// FileUUID Used to uniquely identify the uploaded file
// FileContent Path to the local ClientImpl for the uploaded file
type ScraperPostUploadFileDataModel struct {
	BaseDataModel
	FileUUID    uuid.UUID
	FileContent []byte
}

// PreTransformDataModel is the model for pre-transform data
type PreTransformDataModel struct {
	BaseDataModel
	ScraperType       model.ScraperType
	PreUploadFileData []*ScraperPreUploadFileDataModel
}

// PostTransformDataModel is the actual upload file data
type PostTransformDataModel struct {
	BaseDataModel
	ScraperType        model.ScraperType
	PostUploadFileData []*ScraperPostUploadFileDataModel
}

// AnyDataModel used to limit the total number of total types that have BaseDataModel
type AnyDataModel interface {
	ScraperPreUploadFileDataModel | ScraperPostUploadFileDataModel | PreTransformDataModel | PostTransformDataModel
}
