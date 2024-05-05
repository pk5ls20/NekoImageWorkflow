package impl

import (
	"fmt"
	"github.com/pk5ls20/NekoImageWorkflow/common/log"
	"github.com/robertkrimen/otto"
	"github.com/sirupsen/logrus"
	"io"
	"os"
	"path/filepath"
	"reflect"
	"sync"
)

const JSParseFunctionName = "pasteJson"

type APIParser interface {
	// Init initializes the APIParser
	Init()
	// RegisterParser registers a new parser from a JS file
	RegisterParser(jsFilePath string) error
	// ParseJson parses a rawJSON string using $parserName , which is the name of JS file (aka registered parser)
	ParseJson(rawJson string, parserName string) ([]string, error)
}

type APIParserImpl struct {
	APIParser
	vm          *otto.Otto
	mutex       sync.Mutex
	once        sync.Once
	parserMap   map[string]string
	initialized bool
}

func (a *APIParserImpl) Init() {
	a.mutex.Lock()
	defer a.mutex.Unlock()
	a.vm = otto.New()
	a.parserMap = make(map[string]string)
	a.initialized = true
}

func (a *APIParserImpl) RegisterParser(jsFilePath string) error {
	a.mutex.Lock()
	defer a.mutex.Unlock()
	if !a.initialized {
		return log.ErrorWrap(fmt.Errorf("APIParser not initialized"))
	}
	fileName := filepath.Base(jsFilePath)
	file, err := os.Open(jsFilePath)
	if err != nil {
		return log.ErrorWrap(err)
	}
	defer func(file *os.File) {
		if _err := file.Close(); _err != nil {
			logrus.Error("Error closing file:", _err)
		}
	}(file)
	content, err := io.ReadAll(file)
	if err != nil {
		return log.ErrorWrap(err)
	}
	a.parserMap[fileName] = string(content)
	return nil
}

func (a *APIParserImpl) ParseJson(rawJson string, parserName string) ([]string, error) {
	a.mutex.Lock()
	defer a.mutex.Unlock()
	if !a.initialized {
		return nil, log.ErrorWrap(fmt.Errorf("APIParser not initialized"))
	}
	jsCode, ok := a.parserMap[parserName]
	if !ok {
		return nil, log.ErrorWrap(fmt.Errorf("parser %s not found", parserName))
	}
	_, err := a.vm.Run(jsCode)
	if err != nil {
		return nil, log.ErrorWrap(err)
	}
	value, err := a.vm.Call(JSParseFunctionName, nil, rawJson)
	if err != nil {
		return nil, log.ErrorWrap(err)
	}
	result, err := value.Export()
	if err != nil {
		return nil, log.ErrorWrap(err)
	}
	if strings, _err := convertSliceToStrings(result); _err == nil {
		return strings, nil
	} else {
		return nil, log.ErrorWrap(fmt.Errorf("conversion to []string failed: %v", _err))
	}
}

func convertSliceToStrings(slice interface{}) ([]string, error) {
	s := reflect.ValueOf(slice)
	if s.Kind() != reflect.Slice {
		return nil, log.ErrorWrap(fmt.Errorf("provided value is not a slice, it is %T", slice))
	}
	result := make([]string, s.Len())
	for i := 0; i < s.Len(); i++ {
		element := s.Index(i).Interface()
		result[i] = fmt.Sprintf("%v", element)
	}
	return result, nil
}
