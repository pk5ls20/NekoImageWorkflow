// Package model TODO: update these model to protobuf
package model

import (
	"github.com/google/uuid"
	"github.com/pk5ls20/NekoImageWorkflow/common/log"
	"github.com/pk5ls20/NekoImageWorkflow/common/model"
	uuidTool "github.com/pk5ls20/NekoImageWorkflow/common/uuid"
)

type BaseFileDataModel interface {
}

// FileDataModel Requires specific transport DataModel (just below) implementation
type FileDataModel interface {
	BaseFileDataModel
	// calculateUUID is a function that calculates the UUID of the resourceUri / file
	// only call in the constructor
	calculateUUID() error
	// PrepareUpload is a function that prepares the data for upload, wait to implement
	PrepareUpload() error
	// FinishUpload use to clean up temp files (if exists) after successful upload
	FinishUpload() error
}

// PreUploadFileDataModel
// resourceUUID Used to uniquely identify the resource
// resourceUri resource uri, maybe file path (local) or file url (API)
type PreUploadFileDataModel struct {
	FileDataModel
	scraperType  model.ScraperType
	resourceUUID uuid.UUID
	resourceUri  string
}

// UploadFileDataModel
// fileUUID Used to uniquely identify the uploaded file
// filePath The actual path of the file locally
type UploadFileDataModel struct {
	FileDataModel
	scraperType model.ScraperType
	fileUUID    uuid.UUID
	filePath    string
}

func (s *PreUploadFileDataModel) calculateUUID() error {
	_uuid := uuidTool.GenerateStrUUID(s.resourceUri)
	s.resourceUUID = _uuid
	return nil
}

func (s *UploadFileDataModel) calculateUUID() error {
	_uuid, err := uuidTool.GenerateFileUUID(s.filePath)
	if err != nil {
		return log.ErrorWrap(err)
	}
	s.fileUUID = _uuid
	return nil
}

func NewScraperPreUploadFileData(scType model.ScraperType, uri string) (*PreUploadFileDataModel, error) {
	m := &PreUploadFileDataModel{
		scraperType: scType,
		resourceUri: uri,
	}
	err := m.calculateUUID()
	return m, err
}

func NewScraperUploadFileData(scType model.ScraperType, filePath string) (*UploadFileDataModel, error) {
	m := &UploadFileDataModel{
		scraperType: scType,
		filePath:    filePath,
	}
	err := m.calculateUUID()
	return m, err
}
