package utils

import (
	"bytes"
	"encoding/json"
)

func UnmarshalJsonUseNumber(input string, data interface{}) error {
	decoder := json.NewDecoder(bytes.NewBufferString(input))
	decoder.UseNumber()
	if err := decoder.Decode(data); err != nil {
		return err
	}
	return nil
}
