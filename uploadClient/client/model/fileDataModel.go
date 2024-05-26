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
	// calculateUUID is a function that calculates the UUID of the ResourceUri / file
	// only call in the constructor
	calculateUUID() error
	// GetFileContent is a function that returns the file data
	GetFileContent() ([]byte, error)
	// PrepareUpload is a function that prepares the data for upload, wait to implement
	// TODO: maybe, here can transform the FileDataModel itself to model which adapt to kitexClient proto
	PrepareUpload() error
	// FinishUpload use to clean up temp files (if exists) after successful upload
	FinishUpload() error
}

type PreUploadFileMetaDataModel struct {
	ScraperType  commonModel.ScraperType
	ScraperID    string
	MsgGroupID   string
	ResourceUUID uuid.UUID
	ResourceUri  string
}

// PreUploadFileDataModel
// ResourceUUID Used to uniquely identify the resource
// ResourceUri resource uri, maybe file path (local) or file url (API)
type PreUploadFileDataModel struct {
	FileDataModel
	PreUploadFileMetaDataModel
}

type UploadFileMetaDataModel struct {
	ScraperType commonModel.ScraperType
	ScraperID   string
	MsgGroupID  string
	FileUUID    uuid.UUID
	FilePath    string
	IsTempFile  bool
}

// UploadFileDataModel
// FileUUID Used to uniquely identify the uploaded file
// FilePath The actual path of the file locally
type UploadFileDataModel struct {
	FileDataModel
	UploadFileMetaDataModel
}

type AnyFileMetaDataModel struct {
	*PreUploadFileMetaDataModel
	*UploadFileMetaDataModel
}

func (s *PreUploadFileDataModel) calculateUUID() error {
	_uuid := commonUUID.GenerateStrUUID(s.ResourceUri)
	s.ResourceUUID = _uuid
	return nil
}

func (s *PreUploadFileDataModel) GetFileContent() ([]byte, error) {
	return nil, nil
}

func (s *PreUploadFileDataModel) FinishUpload() error {
	return nil
}

func (s *UploadFileDataModel) calculateUUID() error {
	_uuid, err := commonUUID.GenerateFileUUID(s.FilePath)
	if err != nil {
		return commonLog.ErrorWrap(err)
	}
	s.FileUUID = _uuid
	return nil
}

func (s *UploadFileDataModel) GetFileContent() ([]byte, error) {
	fileContent, err := os.ReadFile(s.FilePath)
	if err != nil {
		return nil, commonLog.ErrorWrap(err)
	}
	return fileContent, nil
}

func (s *UploadFileDataModel) FinishUpload() error {
	if s.IsTempFile {
		logrus.Debug("Delete temp file: ", s.FilePath)
		if err := os.Remove(s.FilePath); err != nil {
			return commonLog.ErrorWrap(err)
		}
	}
	return nil
}

func NewPreUploadFileData(scType commonModel.ScraperType, scID string, MsgGroupID string, uri string) (*PreUploadFileDataModel, error) {
	m := &PreUploadFileDataModel{
		PreUploadFileMetaDataModel: PreUploadFileMetaDataModel{
			ScraperType: scType,
			ScraperID:   scID,
			MsgGroupID:  MsgGroupID,
			ResourceUri: uri,
		},
	}
	err := m.calculateUUID()
	return m, err
}

func NewUploadFileData(model *PreUploadFileDataModel) *UploadFileDataModel {
	m := &UploadFileDataModel{
		UploadFileMetaDataModel: UploadFileMetaDataModel{
			ScraperType: model.ScraperType,
			ScraperID:   model.ScraperID,
			FileUUID:    model.ResourceUUID,
			FilePath:    model.ResourceUri,
			IsTempFile:  false,
		},
	}
	return m
}

func NewUploadTempFileData(scType commonModel.ScraperType, scID string, MsgGroupID string, fileContent []byte) (*UploadFileDataModel, error) {
	t := tmpStorage.NewTmpFile()
	filePath, fileUUID, err := t.Create(fileContent, ".tmp")
	if err != nil {
		return nil, commonLog.ErrorWrap(err)
	}
	m := &UploadFileDataModel{
		UploadFileMetaDataModel: UploadFileMetaDataModel{
			ScraperType: scType,
			ScraperID:   scID,
			MsgGroupID:  MsgGroupID,
			FileUUID:    fileUUID,
			FilePath:    filePath,
			IsTempFile:  true,
		},
	}
	return m, nil
}

func NewPreUploadFileDataFromMeta(meta *PreUploadFileMetaDataModel) *PreUploadFileDataModel {
	return &PreUploadFileDataModel{
		PreUploadFileMetaDataModel: *meta,
	}
}

func NewUploadFileDataFromMeta(meta *UploadFileMetaDataModel) *UploadFileDataModel {
	return &UploadFileDataModel{
		UploadFileMetaDataModel: *meta,
	}
}
