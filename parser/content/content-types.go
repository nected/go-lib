package content

import (
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/clbanning/mxj"
)

func init() {
	registerParserType("application/json", &jsonParser{})
	registerParserType("application/xml", &xmlParser{})
}

type jsonParser struct{}
type xmlParser struct{}

// Unmarshal implements parserType.
func (x *xmlParser) Unmarshal(data []byte, dest interface{}) error {
	resp, err := mxj.NewMapXml(data)
	if err != nil {
		return err
	}
	dest = resp
	return nil
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
