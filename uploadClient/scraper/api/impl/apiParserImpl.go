package impl

import (
	"fmt"
	commonLog "github.com/pk5ls20/NekoImageWorkflow/common/log"
	apiModel "github.com/pk5ls20/NekoImageWorkflow/uploadClient/scraper/api/model"
	"github.com/robertkrimen/otto"
	"github.com/sirupsen/logrus"
	"io"
	"os"
	"reflect"
	"sync"
)

// APIParser is the implementation of Parser, designed as safe to global singleton
type APIParser struct {
	apiModel.APIParser
	vm          *otto.Otto
	mutex       sync.Mutex
	parserMap   map[string]string
	initialized bool
}

func (a *APIParser) Init() {
	a.mutex.Lock()
	defer a.mutex.Unlock()
	a.vm = otto.New()
	a.parserMap = make(map[string]string)
	a.initialized = true
}

func (a *APIParser) RegisterParser(jsFilePath string) error {
	a.mutex.Lock()
	defer a.mutex.Unlock()
	if !a.initialized {
		return commonLog.ErrorWrap(fmt.Errorf("apiParser not initialized"))
	}
	if _, exists := a.parserMap[jsFilePath]; exists {
		logrus.Debugf("Parser %s already exists, content is %s", jsFilePath, a.parserMap[jsFilePath])
		return nil
	}
	file, err := os.Open(jsFilePath)
	if err != nil {
		return commonLog.ErrorWrap(err)
	}
	defer func(file *os.File) {
		if _err := file.Close(); _err != nil {
			logrus.Error("Error closing file:", _err)
		}
	}(file)
	content, err := io.ReadAll(file)
	if err != nil {
		return commonLog.ErrorWrap(err)
	}
	a.parserMap[jsFilePath] = string(content)
	return nil
}

func (a *APIParser) ParseJson(rawJson string, jsFilePath string) ([]string, error) {
	a.mutex.Lock()
	defer a.mutex.Unlock()
	if !a.initialized {
		return nil, commonLog.ErrorWrap(fmt.Errorf("apiParser not initialized"))
	}
	jsCode, ok := a.parserMap[jsFilePath]
	if !ok {
		return nil, commonLog.ErrorWrap(fmt.Errorf("parser %s not found", jsFilePath))
	}
	_, err := a.vm.Run(jsCode)
	if err != nil {
		return nil, commonLog.ErrorWrap(err)
	}
	value, err := a.vm.Call(apiModel.JSParseFunctionName, nil, rawJson)
	if err != nil {
		return nil, commonLog.ErrorWrap(err)
	}
	result, err := value.Export()
	if err != nil {
		return nil, commonLog.ErrorWrap(err)
	}
	if strings, _err := convertSliceToStrings(result); _err == nil {
		return strings, nil
	} else {
		return nil, commonLog.ErrorWrap(fmt.Errorf("conversion to []string failed: %v", _err))
	}
}

func convertSliceToStrings(slice interface{}) ([]string, error) {
	s := reflect.ValueOf(slice)
	if s.Kind() != reflect.Slice {
		return nil, commonLog.ErrorWrap(fmt.Errorf("provided value is not a slice, it is %T", slice))
	}
	result := make([]string, s.Len())
	for i := 0; i < s.Len(); i++ {
		element := s.Index(i).Interface()
		result[i] = fmt.Sprintf("%v", element)
	}
	return result, nil
}
