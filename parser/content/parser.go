package content

var cTypeParserMap = map[string]parserType{}

type parserType interface {
	Unmarshal(data []byte, dest interface{}) error
}

func NewParserType(contentType string) parserType {
	return cTypeParserMap[contentType]
}

func RegisterParserType(contentType string, parserType parserType) {
	cTypeParserMap[contentType] = parserType
}
