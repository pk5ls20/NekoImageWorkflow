package model

const JSParseFunctionName = "pasteJson"

// APIParser is the interface for parsing JSON data to []string
type APIParser interface {
	// Init initializes the apiParser
	Init()
	// RegisterParser registers a new parser from a JS file
	RegisterParser(jsFilePath string) error
	// ParseJson parses a rawJSON string using $parserName , which is the name of JS file (aka registered parser)
	ParseJson(rawJson string, parserName string) ([]string, error)
}
