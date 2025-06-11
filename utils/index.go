package utils

import "encoding/json"

func JsonToStruct(inputData interface{}, outputData interface{}) error {
	marshalData, err := json.Marshal(inputData)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(marshalData, &outputData); err != nil {
		return err
	}

	return nil
}
