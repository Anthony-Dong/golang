package trimer

import (
	"encoding/json"
	"io"
	"os"
	"reflect"

	"github.com/iancoleman/orderedmap"

	"github.com/spf13/cobra"

	"github.com/anthony-dong/golang/pkg/utils"
)

func NewJsonTrimerCommand() (*cobra.Command, error) {
	cmd := &cobra.Command{
		Use:   "json",
		Short: "Trim or filter JSON data",
	}
	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		if !utils.CheckStdInFromPiped() {
			return cmd.Help()
		}
		all, err := io.ReadAll(os.Stdin)
		if err != nil {
			return err
		}
		out, err := TrimJson(all)
		if err != nil {
			return err
		}
		os.Stdout.Write(out)
		return nil
	}
	return cmd, nil
}

func TrimJson(data []byte) ([]byte, error) {
	orderedMap := orderedmap.OrderedMap{}
	//orderedMap.SetUseNumber(true)
	orderedMap.SetEscapeHTML(true)
	if err := json.Unmarshal(data, &orderedMap); err != nil {
		return nil, err
	}
	result := handleData(orderedMap)
	return json.MarshalIndent(result, "", "  ")
}

func handleData(data interface{}) interface{} {
	switch v := data.(type) {
	case *orderedmap.OrderedMap:
		return handleOrderMap(v)
	case orderedmap.OrderedMap:
		return handleOrderMap(&v)
	case []interface{}:
		result := make([]interface{}, 0, len(v))
		for _, elem := range v {
			data = handleData(elem)
			if data == nil {
				continue
			}
			result = append(result, data)
		}
		if len(result) == 0 {
			return nil
		}
		return result
	default:
		if data == nil {
			return nil
		}
		if reflect.ValueOf(data).IsZero() {
			return nil
		}
		return data
	}
}

func handleOrderMap(v *orderedmap.OrderedMap) interface{} {
	keys := v.Keys()
	newMap := orderedmap.New()
	for _, key := range keys {
		if key == "" {
			continue
		}
		value, _ := v.Get(key)
		newValue := handleData(value)
		if newValue == nil {
			continue
		}
		newMap.Set(key, newValue)
	}
	if len(newMap.Keys()) == 0 {
		return nil
	}
	return newMap
}
