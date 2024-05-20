package tmp

import (
	"bytes"
	"github.com/sirupsen/logrus"
	"os"
	"sync"
	"testing"
)

func TestCreateAndDelete(t *testing.T) {
	tf := NewTmpFile()
	content := []byte("hello world-1")
	ext := ".txt"
	filePath, _, _err := tf.Create(content, ext)
	if _err != nil {
		t.Errorf("Failed to create file: %s", _err)
	}
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		t.Errorf("File does not exist after creation")
	}
	readContent, _err := os.ReadFile(filePath)
	if _err != nil {
		t.Errorf("Failed to read file: %s", _err)
	}
	if !bytes.Equal(content, readContent) {
		t.Errorf("File content mismatch: expected %s, got %s", content, readContent)
	}
	if err := tf.Delete(filePath); err != nil {
		t.Errorf("Failed to delete file: %s", err)
	}
	if _, err := os.Stat(filePath); !os.IsNotExist(err) {
		t.Errorf("File exists after deletion")
	}
}

func TestConcurrentAccess(t *testing.T) {
	tf := NewTmpFile()
	content := []byte("test data-2")
	ext := ".dat"
	var wg sync.WaitGroup
	n := 1000
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			filePath, _, err := tf.Create(content, ext)
			if err != nil {
				t.Errorf("Failed to create file: %s", err)
			}
			if _err := tf.Delete(filePath); _err != nil {
				t.Errorf("Failed to delete file: %s", _err)
			}
		}()
	}
	wg.Wait()
}

func TestUniqueIDCreation(t *testing.T) {
	tf := NewTmpFile()
	content1 := []byte("content-3")
	content2 := []byte("another content-3")
	ext := ".log"
	filePath1, fileUUID1, err1 := tf.Create(content1, ext)
	if err1 != nil {
		t.Errorf("Failed to create first file: %s", err1)
	}
	filePath2, fileUUID2, err2 := tf.Create(content2, ext)
	if err2 != nil {
		t.Errorf("Failed to create second file: %s", err2)
	}
	if fileUUID1.String() == fileUUID2.String() {
		t.Errorf("uuid1 %v == %v", fileUUID1, fileUUID2)
	}
	if filePath1 == filePath2 {
		t.Errorf("filepath1 %s == filepath2 %s", filePath1, filePath2)
	}
	if err := tf.Delete(filePath1); err != nil {
		t.Errorf("Failed to delete first file: %s", err)
	}
	if err := tf.Delete(filePath2); err != nil {
		t.Errorf("Failed to delete second file: %s", err)
	}
}

func TestMain(m *testing.M) {
	code := m.Run()
	if err := os.RemoveAll(TmpDir); err != nil {
		logrus.Fatalf("Failed to remove TMP_DIR: %s", err)
	}
	logrus.Debugf("Removed TMP_DIR: %s", TmpDir)
	os.Exit(code)
}
