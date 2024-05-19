package tmp

import (
	"github.com/google/uuid"
	commonLog "github.com/pk5ls20/NekoImageWorkflow/common/log"
	commonUtils "github.com/pk5ls20/NekoImageWorkflow/common/utils"
	commonUUID "github.com/pk5ls20/NekoImageWorkflow/common/uuid"
	"github.com/sirupsen/logrus"
	"os"
	"path/filepath"
	"strconv"
	"sync"
)

const tmpDir = "_tmp"

type tmpFile interface {
	// Create creates a temporary file
	Create(fileContent []byte, extension string) (filePath string, fileUUID uuid.UUID, err error)
	Delete(filePath string) error
}

type TmpFile struct {
	tmpFile
	lockChan *commonUtils.IDLock
	once     sync.Once
}

var instance = &TmpFile{}

func NewTmpFile() *TmpFile {
	instance.once.Do(func() {
		if _, err := os.Stat(tmpDir); os.IsNotExist(err) {
			if err := os.Mkdir(tmpDir, 0755); err != nil {
				logrus.Error("Failed to create tmp directory: ", err)
			}
		}
		instance = &TmpFile{
			lockChan: commonUtils.NewIDLock(),
		}
	})
	return instance
}

// modifyLockFile add changeVal to the count in the lock file, return the new count
// if the lockfile don't exist, create a new one with count = changeVal
// if the final count<=0, the lock file will be deleted
func (t *TmpFile) modifyLockFile(lockFilePath string, changeVal int) (int, error) {
	logrus.Debug("Modify lock file: ", lockFilePath)
	var count = 0
	if _, err := os.Stat(lockFilePath); os.IsNotExist(err) {
		logrus.Debugf("Lock file %s not exists, create a new one with count = %d", lockFilePath, changeVal)
		count = 0
	} else {
		data, _err := os.ReadFile(lockFilePath)
		if _err != nil {
			return 0, _err
		}
		count, _err = strconv.Atoi(string(data))
		if _err != nil {
			return 0, _err
		}
	}
	count = count + changeVal
	// final count?
	if count <= 0 {
		if _err := os.Remove(lockFilePath); _err != nil {
			return 0, _err
		}
		return 0, nil
	} else {
		if _err := os.WriteFile(lockFilePath, []byte(strconv.Itoa(count)), 0644); _err != nil {
			return 0, _err
		}
	}
	return count, nil
}

// Create creates a temporary file
// fileContent the content of the file
// extension the extension of the file (e.g. ".jpg")
func (t *TmpFile) Create(fileContent []byte, extension string) (filePath string, fileUUID uuid.UUID, err error) {
	genUUID := commonUUID.GenerateByteSliceUUID(fileContent)
	fileUUIDString := genUUID.String()
	fileName := fileUUIDString + extension
	filePath = filepath.Join(tmpDir, fileName)
	t.lockChan.Lock(filePath)
	defer func(lockChan *commonUtils.IDLock, id string) {
		if _err := lockChan.Unlock(id); _err != nil {
			logrus.Fatal("Failed to unlock: ", _err)
		}
	}(t.lockChan, filePath)
	lockfilePath := filepath.Join(tmpDir, fileName+".lock")
	// if file already exists, just control lockfile
	if _, _err := os.Stat(filePath); _err == nil {
		var newCount int
		if newCount, err = t.modifyLockFile(lockfilePath, 1); err != nil {
			return filePath, genUUID, commonLog.ErrorWrap(err)
		}
		logrus.Debugf("File %s already exists, new count = %d", filePath, newCount)
		return filePath, genUUID, nil
	}
	// if file itself not exists, create the file itself
	if err = os.WriteFile(filePath, fileContent, 0644); err != nil {
		return "", genUUID, commonLog.ErrorWrap(err)
	}
	logrus.Debug("Create temp file: ", filePath)
	return filePath, genUUID, nil
}

// Delete deletes a temporary file
func (t *TmpFile) Delete(filePath string) error {
	t.lockChan.Lock(filePath)
	defer func(lockChan *commonUtils.IDLock, id string) {
		if err := lockChan.Unlock(id); err != nil {
			logrus.Fatal("Failed to unlock: ", err)
		}
	}(t.lockChan, filePath)
	if _, _err := os.Stat(filePath); _err == nil {
		lockFile := filePath + ".lock"
		// if lockfile not exists, delete the file itself
		if _, err := os.Stat(lockFile); os.IsNotExist(err) {
			logrus.Debug("lockfile not exists, so delete file: ", filePath)
			if __err := os.Remove(filePath); __err != nil {
				return commonLog.ErrorWrap(__err)
			}
			return nil
		}
		// if lockfile exists, just control lockfile
		lockCount, err := t.modifyLockFile(lockFile, -1)
		if err != nil {
			return commonLog.ErrorWrap(err)
		}
		logrus.Debugf("File %s have lockfile, new count = %d", filePath, lockCount)
		return nil
	}
	return nil
}
