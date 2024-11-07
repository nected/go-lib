package content

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"reflect"
)

func init() {
	registerParserType("application/json", &jsonParser{})
	registerParserType("application/xml", &xmlParser{})
}

type jsonParser struct{}
type xmlParser struct{}

// Unmarshal implements parserType.
func (x *xmlParser) Unmarshal(data []byte, dest interface{}) error {
	return xml.Unmarshal(data, &dest)
}

func (j *jsonParser) Unmarshal(data []byte, dest interface{}) error {
	return json.Unmarshal(data, &dest)
}

func validatePtr(data interface{}) error {
	if reflect.ValueOf(data).Kind() != reflect.Ptr {
		return fmt.Errorf("dest must be a pointer")
	}
	return nil
}
