package tmp

import (
	"github.com/google/uuid"
	"github.com/pk5ls20/NekoImageWorkflow/common/log"
	uuidTool "github.com/pk5ls20/NekoImageWorkflow/common/uuid"
	"github.com/sirupsen/logrus"
	"os"
	"path/filepath"
)

type tmpFile interface {
	// Create creates a temporary file
	Create(content []byte, extension string) (filename string, fileUUID uuid.UUID, err error)
}

type TmpFile struct {
	tmpFile
}

// Create creates a temporary file
// TODO: if the file already exists?
func (t *TmpFile) Create(fileContent []byte, extension string) (filePath string, fileUUID uuid.UUID, err error) {
	genUUID := uuidTool.GenerateByteSliceUUID(fileContent)
	fileName := genUUID.String() + extension
	tmpDir := "_tmp"
	if _, _err := os.Stat(tmpDir); os.IsNotExist(_err) {
		if __err := os.Mkdir(tmpDir, 0755); __err != nil {
			return "", genUUID, log.ErrorWrap(__err)
		}
	}
	filePath = filepath.Join(tmpDir, fileName)
	if err = os.WriteFile(filePath, fileContent, 0644); err != nil {
		return "", genUUID, log.ErrorWrap(err)
	}
	logrus.Debug("Create temp file: ", filePath)
	return filePath, genUUID, nil
}
