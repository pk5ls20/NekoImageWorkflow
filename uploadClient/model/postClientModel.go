// Package model TODO: update these model to protobuf
package model

import (
	"NekoImageWorkflowKitex/common"
	"github.com/google/uuid"
)

// PreUploadFileData
// ResourceUUID Used to uniquely identify the resource
// ResourceUri Path to the local ClientImpl for the resource
type PreUploadFileData struct {
	ResourceUUID uuid.UUID
	ResourceUri  string
}

// UploadFileData
// FileUUID Used to uniquely identify the uploaded file
// FileContent Path to the local ClientImpl for the uploaded file
type UploadFileData struct {
	FileUUID    uuid.UUID
	FileContent []byte
}

// PreTransformDataModel is the model for pre-transform data,
// if FileUUID in server, then don't upload
type PreTransformDataModel struct {
	common.ScraperType
	PreUploadFileData []*PreUploadFileData
}

// PostTransformDataModel is the actual upload file data
type PostTransformDataModel struct {
	common.ScraperType
	PostUploadFileData []*UploadFileData
}
