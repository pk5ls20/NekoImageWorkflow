package impl

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"sync"
	"testing"
)

// Helper function to create a temporary JS file
func createTempFile(content string, id string, t *testing.T) (string, func()) {
	tmpFile, err := os.CreateTemp("", fmt.Sprintf("test_%s_*.js", id))
	if err != nil {
		t.Fatal(err)
	}
	if _, err = tmpFile.Write([]byte(content)); err != nil {
		t.Fatal(err)
	}
	if err = tmpFile.Close(); err != nil {
		t.Fatal(err)
	}
	return tmpFile.Name(), func() {
		err := os.Remove(tmpFile.Name())
		if err != nil {
			return
		}
	}
}

// testing Register
func TestAPIParserImpl_RegisterSingleParser_Success(t *testing.T) {
	parser := APIParserImpl{}
	parser.Init()
	filePath, cleanup := createTempFile("function pasteJson(json) { return JSON.parse(json); }", "0", t)
	defer cleanup()
	err := parser.RegisterParser(filePath)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if _, exists := parser.parserMap[filepath.Base(filePath)]; !exists {
		t.Errorf("Parser was not registered properly")
	}
}

func TestAPIParserImpl_RegisterMultipleParsers_Success(t *testing.T) {
	parser := APIParserImpl{}
	parser.Init()
	filePath1, cleanup1 := createTempFile("function pasteJson(json) { return JSON.parse(json); }", "0", t)
	defer cleanup1()
	filePath2, cleanup2 := createTempFile("function pasteJson(json) { return JSON.parse(json); }", "0", t)
	defer cleanup2()
	if err := parser.RegisterParser(filePath1); err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if err := parser.RegisterParser(filePath2); err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if _, exists := parser.parserMap[filepath.Base(filePath1)]; !exists {
		t.Errorf("First parser was not registered properly")
	}
	if _, exists := parser.parserMap[filepath.Base(filePath2)]; !exists {
		t.Errorf("Second parser was not registered properly")
	}
}

func TestAPIParserImpl_RegisterParser_Fail(t *testing.T) {
	parser := APIParserImpl{}
	parser.Init()
	if err := parser.RegisterParser("non_existent_file.js"); err == nil {
		t.Errorf("Expected an error for nonexistent file, got nil")
	}
}

// testing ParseJson
func TestAPIParserImpl_ParseJson_Success(t *testing.T) {
	parser := APIParserImpl{}
	parser.Init()
	jsContent := `
	function pasteJson(json) {
		var data = JSON.parse(json);
		return data.key;
	}`
	filePath, cleanup := createTempFile(jsContent, "0", t)
	defer cleanup()
	if err := parser.RegisterParser(filePath); err != nil {
		t.Fatal("Failed to register parser:", err)
	}
	result, err := parser.ParseJson(`{"key":[0,1,2,3,4,5]}`, filepath.Base(filePath))
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	expected := []string{"0", "1", "2", "3", "4", "5"}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

func TestAPIParserImpl_ParseJson_JSError(t *testing.T) {
	parser := APIParserImpl{}
	parser.Init()
	jsContent := `function pasteJson(json) { return JSON.parse(json); }`
	filePath, cleanup := createTempFile(jsContent, "0", t)
	defer cleanup()
	if err := parser.RegisterParser(filePath); err != nil {
		t.Fatal("Failed to register parser:", err)
	}
	if _, err := parser.ParseJson(`{invalid_json}`, filepath.Base(filePath)); err == nil {
		t.Errorf("Expected an error due to invalid JSON, got nil")
	}
}

func TestAPIParserImpl_ParseJson_NoInit(t *testing.T) {
	parser := APIParserImpl{}
	jsContent := `function pasteJson(json) { return JSON.parse(json); }`
	filePath, cleanup := createTempFile(jsContent, "0", t)
	defer cleanup()
	if err := parser.RegisterParser(filePath); err == nil {
		t.Fatal("Excepted an error due to not init", err)
	}
	if _, err := parser.ParseJson(`{}`, filepath.Base(filePath)); err == nil {
		t.Fatal("Excepted an error due to not init, got", err)
	}
}

type pasteTask struct {
	id     int
	result []string
}

func TestAPIParserImpl_Concurrency(t *testing.T) {
	parser := &APIParserImpl{}
	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			parser.Init()
		}()
	}
	wg.Wait()
	if !parser.initialized {
		t.Errorf("APIParserImpl was not initialized correctly")
	}
	wg = sync.WaitGroup{}
	results := make(chan pasteTask, 100)
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(id int) {
			task := pasteTask{id: id}
			defer wg.Done()
			filePath, cleanup := createTempFile(
				fmt.Sprintf("function pasteJson(json) { return JSON.parse(json).key%s; }", strconv.Itoa(id)),
				strconv.Itoa(id), t)
			defer cleanup()
			if err := parser.RegisterParser(filePath); err != nil {
				t.Errorf("error registering parser with %s: %v", filePath, err)
			}
			result, err := parser.ParseJson(fmt.Sprintf("{\"key%s\": [%s]}", strconv.Itoa(id), strconv.Itoa(id)),
				filepath.Base(filePath))
			if err != nil {
				t.Errorf("error parsing JSON with %s: %v", filePath, err)
			} else {
				task.result = result
				results <- task
			}
		}(i)
	}
	wg.Wait()
	close(results)
	if len(parser.parserMap) != 100 {
		t.Errorf("expected 10 parsers to be registered, got %d", len(parser.parserMap))
	}
	for task := range results {
		expected := []string{strconv.Itoa(task.id)}
		if !reflect.DeepEqual(task.result, expected) {
			t.Errorf("expected %v, got %v", expected, task.result)
		}
	}
}
