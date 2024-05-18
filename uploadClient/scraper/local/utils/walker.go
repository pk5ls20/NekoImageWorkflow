package utils

import (
	commonLog "github.com/pk5ls20/NekoImageWorkflow/common/log"
	"io/fs"
	"path/filepath"
)

func WalkDir(rootPath string) (*[]string, error) {
	dirs := make([]string, 0)
	dirs = append(dirs, rootPath)
	walkDirErr := filepath.WalkDir(rootPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return commonLog.ErrorWrap(err)
		}
		if d.IsDir() {
			dirs = append(dirs, path)
		}
		return nil
	})
	if walkDirErr != nil {
		return nil, commonLog.ErrorWrap(walkDirErr)
	}
	return &dirs, nil
}
