// Package model TODO: update these model to protobuf
package model

import (
	"github.com/google/uuid"
	"github.com/pk5ls20/NekoImageWorkflow/common/log"
	"github.com/pk5ls20/NekoImageWorkflow/common/model"
	uuidTool "github.com/pk5ls20/NekoImageWorkflow/common/uuid"
)

type BaseBridgeDataModel interface {
}

// BridgeDataModel Requires specific transport DataModel (just below) implementation
type BridgeDataModel interface {
	BaseBridgeDataModel
	// calculateUUID is a function that calculates the UUID of the resourceUri / file
	// only call in the constructor
	calculateUUID() error
}

type BridgeTransformDataModel interface {
	BaseBridgeDataModel
	// PrepareUpload is a function that prepares the data for upload, wait to implement
	PrepareUpload() error
	// FinishUpload use to clean up temp files (if exists) after successful upload
	FinishUpload() error
}

// ScraperPreUploadFileDataModel
// resourceUUID Used to uniquely identify the resource
// resourceUri resource uri, maybe file path (local) or file url (API)
type ScraperPreUploadFileDataModel struct {
	BridgeDataModel
	scraperType  model.ScraperType
	resourceUUID uuid.UUID
	resourceUri  string
}

// ScraperPostUploadFileDataModel
// fileUUID Used to uniquely identify the uploaded file
// filePath The actual path of the file locally
type ScraperPostUploadFileDataModel struct {
	BridgeDataModel
	scraperType model.ScraperType
	fileUUID    uuid.UUID
	filePath    string
}

// PreTransformDataModel is the model for pre-transform data (aka file metadata)
type PreTransformDataModel struct {
	BridgeTransformDataModel
	preUploadFileData []*ScraperPreUploadFileDataModel
}

// PostTransformDataModel is the actual upload file data
type PostTransformDataModel struct {
	BridgeTransformDataModel
	postUploadFileData []*ScraperPostUploadFileDataModel
}

func (s *ScraperPreUploadFileDataModel) calculateUUID() error {
	_uuid := uuidTool.GenerateStrUUID(s.resourceUri)
	s.resourceUUID = _uuid
	return nil
}

func (s *ScraperPostUploadFileDataModel) calculateUUID() error {
	_uuid, err := uuidTool.GenerateFileUUID(s.filePath)
	if err != nil {
		return log.ErrorWrap(err)
	}
	s.fileUUID = _uuid
	return nil
}

func NewScraperPreUploadFileData(scType model.ScraperType, uri string) (*ScraperPreUploadFileDataModel, error) {
	m := &ScraperPreUploadFileDataModel{
		scraperType: scType,
		resourceUri: uri,
	}
	err := m.calculateUUID()
	return m, err
}

func NewScraperPostUploadFileData(scType model.ScraperType, filePath string) (*ScraperPostUploadFileDataModel, error) {
	m := &ScraperPostUploadFileDataModel{
		scraperType: scType,
		filePath:    filePath,
	}
	err := m.calculateUUID()
	return m, err
}

func (p *PreTransformDataModel) PrepareUpload() error {
	return nil
}

func (p *PreTransformDataModel) FinishUpload() error {
	return nil
}

func (p *PostTransformDataModel) PrepareUpload() error {
	return nil
}

func (p *PostTransformDataModel) FinishUpload() error {
	return nil
}
