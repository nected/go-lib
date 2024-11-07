package content

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"reflect"
)

func init() {
	RegisterParserType("application/json", &jsonParser{})
	RegisterParserType("application/xml", &xmlParser{})
}

type jsonParser struct{}
type xmlParser struct{}

// Unmarshal implements parserType.
func (x *xmlParser) Unmarshal(data []byte, dest interface{}) error {
	if err := validatePtr(dest); err != nil {
		return err
	}
	return xml.Unmarshal(data, dest)
}

func (j *jsonParser) Unmarshal(data []byte, dest interface{}) error {
	if err := validatePtr(dest); err != nil {
		return err
	}
	return json.Unmarshal(data, dest)
}

func validatePtr(data interface{}) error {
	if reflect.ValueOf(data).Kind() != reflect.Ptr {
		return fmt.Errorf("dest must be a pointer")
	}
	return nil
}
