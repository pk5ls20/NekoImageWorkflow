// Package model TODO: update these model to protobuf
package model

import (
	"github.com/google/uuid"
	commonLog "github.com/pk5ls20/NekoImageWorkflow/common/log"
	commonModel "github.com/pk5ls20/NekoImageWorkflow/common/model"
	commonUUID "github.com/pk5ls20/NekoImageWorkflow/common/uuid"
	tmpStorage "github.com/pk5ls20/NekoImageWorkflow/uploadClient/storage/tmp"
	"github.com/sirupsen/logrus"
	"os"
)

type BaseFileDataModel interface {
}

// FileDataModel Requires specific transport DataModel (just below) implementation
type FileDataModel interface {
	BaseFileDataModel
	// calculateUUID is a function that calculates the UUID of the resourceUri / file
	// only call in the constructor
	calculateUUID() error
	GetScraperID() int
	// PrepareUpload is a function that prepares the data for upload, wait to implement
	// TODO: maybe, here can transform the FileDataModel itself to model which adapt to kitexClient proto
	PrepareUpload() error
	// FinishUpload use to clean up temp files (if exists) after successful upload
	FinishUpload() error
}

// PreUploadFileDataModel
// resourceUUID Used to uniquely identify the resource
// resourceUri resource uri, maybe file path (local) or file url (API)
type PreUploadFileDataModel struct {
	FileDataModel
	scraperType  commonModel.ScraperType
	scraperID    int
	resourceUUID uuid.UUID
	resourceUri  string
}

// UploadFileDataModel
// fileUUID Used to uniquely identify the uploaded file
// filePath The actual path of the file locally
type UploadFileDataModel struct {
	FileDataModel
	scraperType commonModel.ScraperType
	scraperID   int
	fileUUID    uuid.UUID
	filePath    string
	isTempFile  bool
}

func (s *PreUploadFileDataModel) calculateUUID() error {
	_uuid := commonUUID.GenerateStrUUID(s.resourceUri)
	s.resourceUUID = _uuid
	return nil
}

func (s *PreUploadFileDataModel) GetScraperID() int {
	return s.scraperID
}

func (s *UploadFileDataModel) calculateUUID() error {
	_uuid, err := commonUUID.GenerateFileUUID(s.filePath)
	if err != nil {
		return commonLog.ErrorWrap(err)
	}
	s.fileUUID = _uuid
	return nil
}

func (s *UploadFileDataModel) GetScraperID() int {
	return s.scraperID
}

func (s *UploadFileDataModel) FinishUpload() error {
	if s.isTempFile {
		logrus.Debug("Delete temp file: ", s.filePath)
		if err := os.Remove(s.filePath); err != nil {
			return commonLog.ErrorWrap(err)
		}
	}
	return nil
}

func NewScraperPreUploadFileData(scType commonModel.ScraperType, scID int, uri string) (*PreUploadFileDataModel, error) {
	m := &PreUploadFileDataModel{
		scraperType: scType,
		resourceUri: uri,
		scraperID:   scID,
	}
	err := m.calculateUUID()
	return m, err
}

func NewScraperUploadFileData(scType commonModel.ScraperType, scID int, filePath string) (*UploadFileDataModel, error) {
	m := &UploadFileDataModel{
		scraperType: scType,
		filePath:    filePath,
		scraperID:   scID,
	}
	err := m.calculateUUID()
	return m, err
}

func NewScraperUploadTempFileData(scType commonModel.ScraperType, scID int, fileContent []byte) (*UploadFileDataModel, error) {
	t := tmpStorage.NewTmpFile()
	filePath, fileUUID, err := t.Create(fileContent, ".tmp")
	if err != nil {
		return nil, commonLog.ErrorWrap(err)
	}
	m := &UploadFileDataModel{
		scraperType: scType,
		filePath:    filePath,
		isTempFile:  true,
		fileUUID:    fileUUID,
		scraperID:   scID,
	}
	return m, nil
}
