package cpp

import (
	"encoding/json"
	"runtime"
	"strings"

	"github.com/anthony-dong/golang/pkg/utils"
)

func ReadToolConfigFromFile(fileName string, tool *Tools) error {
	kv := make(map[string]interface{})
	if err := utils.UnmarshalFromFile(fileName, kv); err != nil {
		return err
	}
	return ReadToolConfigFromKV(kv, runtime.GOOS, tool)
}

func ReadToolConfigFromKV(kv map[string]interface{}, os string, tool *Tools) error {
	osKv := make(map[string]interface{}, len(kv))
	for key, value := range kv {
		keys := strings.SplitN(key, "@", 2)
		if len(keys) != 2 {
			continue
		}
		if keys[1] == os {
			osKv[keys[0]] = value
		}
	}
	for k, v := range osKv {
		kv[k] = v
	}
	marshal, err := json.Marshal(kv)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(marshal, tool); err != nil {
		return err
	}
	return nil
}
