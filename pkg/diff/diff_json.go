package diff

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/anthony-dong/golang/pkg/orderedmap"
	"github.com/anthony-dong/golang/pkg/utils"
)

const AddType = "add"
const DelType = "del"
const ChangeType = "change"

type Diff struct {
	Key  string
	Type string

	Origin interface{} `json:",omitempty"`
	Patch  interface{} `json:",omitempty"`
}

func unmarshalJson(data []byte) (interface{}, error) {

	orderedMap := orderedmap.New()
	orderedMap.SetUseNumber(true)
	if err := json.Unmarshal(data, orderedMap); err != nil {
		return unmarshalArray(data)
	}
	return orderedMap, nil
}

func unmarshalArray(data []byte) ([]interface{}, error) {
	messages := make([]interface{}, 0)
	if err := json.Unmarshal(data, &messages); err != nil {
		return nil, err
	}
	return messages, nil
}

func DiffJsonString(json1 string, json2 string) ([]*Diff, error) {
	return DiffJson(utils.String2Bytes(json1), utils.String2Bytes(json2))
}

func DiffJson(json1 []byte, json2 []byte) ([]*Diff, error) {
	cmp1, err := unmarshalJson(json1)
	if err != nil {
		return nil, err
	}
	cmp2, err := unmarshalJson(json2)
	if err != nil {
		return nil, err
	}
	return diffValue("", cmp1, cmp2), nil
}

func newDiff(Key string, v1, v2 interface{}) []*Diff {
	return []*Diff{
		{
			Key:    Key,
			Origin: v1,
			Patch:  v2,
			Type: func() string {
				if v1 == nil {
					return AddType
				}
				if v2 == nil {
					return DelType
				}
				return ChangeType
			}(),
		},
	}
}

func diffValue(prefix string, v1 interface{}, v2 interface{}) []*Diff {
	if v1 == nil && v2 == nil {
		return nil
	}
	switch vv1 := v1.(type) {
	case map[string]interface{}:
		orderedMap, isOK := toOrderedMap(vv1)
		if !isOK {
			panic(`cannot convert value to ordered map`)
		}
		return diffOrderedMap(prefix, orderedMap, v2)
	case *orderedmap.OrderedMap:
		return diffOrderedMap(prefix, vv1, v2)
	case orderedmap.OrderedMap:
		return diffOrderedMap(prefix, &vv1, v2)
	case json.Number:
		vv2, isOK := v2.(json.Number)
		if !isOK {
			return newDiff(prefix, v1, v2)
		}
		if !bytes.Equal([]byte(vv1), []byte(vv2)) {
			return newDiff(prefix, v1, v2)
		}
		return nil
	case string:
		vv2, isOK := v2.(string)
		if !isOK {
			return newDiff(prefix, v1, v2)
		}
		if vv1 != vv2 {
			return newDiff(prefix, v1, v2)
		}
		return nil
	case bool:
		vv2, isOK := v2.(bool)
		if !isOK {
			return newDiff(prefix, v1, v2)
		}
		if vv1 != vv2 {
			return newDiff(prefix, v1, v2)
		}
		return nil
	case []interface{}:
		vv2, isOK := v2.([]interface{})
		if !isOK {
			return newDiff(prefix, v1, v2)
		}
		diffs := make([]*Diff, 0)
		lastIndex := len(vv1)
		for index, _ := range vv1 {
			if index < len(vv2) {
				diffs = append(diffs, diffValue(prefix+"."+formatIndex(index), vv1[index], vv2[index])...)
			} else {
				diffs = append(diffs, diffValue(prefix+"."+formatIndex(index), vv1[index], nil)...)
			}
		}
		if len(vv2) > len(vv1) {
			for index := lastIndex; index < len(vv2); index++ {
				diffs = append(diffs, diffValue(prefix+"."+formatIndex(index), nil, vv2[index])...)
			}
		}
		return diffs
	case nil:
		if v2 != nil {
			return newDiff(prefix, v1, v2)
		}
		return nil
	default:
		panic(fmt.Sprintf(`unsupported type [%T]`, v1))
	}
}

func diffOrderedMap(prefix string, vv1 *orderedmap.OrderedMap, v2 interface{}) []*Diff {
	vv2, isOK := toOrderedMap(v2)
	if !isOK {
		return newDiff(prefix, vv1, v2)
	}
	diffs := make([]*Diff, 0)
	vv1.Foreach(func(key string, value interface{}) {
		diffs = append(diffs, diffValue(prefix+"."+key, value, vv2.GetOr(key))...)
	})
	vv2.Foreach(func(key string, value interface{}) {
		if vv1.Exist(key) {
			return
		}
		diffs = append(diffs, diffValue(prefix+"."+key, nil, value)...)
	})
	return diffs
}

func toOrderedMap(v interface{}) (*orderedmap.OrderedMap, bool) {
	switch vv := v.(type) {
	case *orderedmap.OrderedMap:
		return vv, true
	case orderedmap.OrderedMap:
		return &vv, true
	case map[string]interface{}:
		orderedMap := orderedmap.New()
		for key, value := range vv {
			orderedMap.Set(key, value)
		}
		return orderedMap, true
	}
	return nil, false
}

func formatIndex(index int) string {
	return "[" + strconv.Itoa(index) + "]"
}
