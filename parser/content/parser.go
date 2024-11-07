package content

import "strings"

var cTypeParserMap = map[string]parserType{}

type parserType interface {
	Unmarshal(data []byte, dest interface{}) error
}

func newParserType(contentType string) parserType {
	contentTypeList := strings.Split(contentType, ";")
	if contentTypeList[0] == "" {
		contentType = "application/json"
	} else {
		contentType = contentTypeList[0]
	}
	return cTypeParserMap[contentType]
}

func registerParserType(contentType string, parserType parserType) {
	cTypeParserMap[contentType] = parserType
}

func Unmarshal(contentType string, data []byte, dest interface{}) error {
	parser := newParserType(contentType)
	if parser == nil {
		return nil
	}

	// Unmarshal the data and if error return it as string
	if err := parser.Unmarshal(data, dest); err != nil {
		dest = string(data)
	}
	return nil
}
