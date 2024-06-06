package utils

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

func ReadKV(str string) (string, string) {
	kvs := strings.SplitN(str, "=", 2)
	if len(kvs) != 2 {
		return "", ""
	}
	return strings.TrimSpace(kvs[0]), strings.TrimSpace(kvs[1])
}

func ReadKVBySep(str string, sep string) (string, string) {
	kvs := strings.SplitN(str, sep, 2)
	if len(kvs) != 2 {
		return "", ""
	}
	return strings.TrimSpace(kvs[0]), strings.TrimSpace(kvs[1])
}

func SplitKV(str string, sep string) (string, string) {
	return ReadKVBySep(str, sep)
}

func UnmarshalFromFile(file string, v interface{}) error {
	content, err := os.ReadFile(file)
	if err != nil {
		return fmt.Errorf(`read file %s find err: %v`, file, err)
	}
	if err := json.Unmarshal(content, v); err == nil {
		return nil
	}
	if err := yaml.Unmarshal(content, v); err != nil {
		return err
	}
	return nil
}
