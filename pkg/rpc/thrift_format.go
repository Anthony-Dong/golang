package rpc

import (
	"encoding/json"
	"fmt"

	"github.com/anthony-dong/golang/pkg/utils"
	"github.com/cloudwego/kitex/pkg/generic/descriptor"
	"github.com/iancoleman/orderedmap"
)

func formatStruct(desc *descriptor.StructDescriptor, data interface{}) (interface{}, error) {
	datas, isOK := data.(map[string]interface{})
	if !isOK {
		return data, nil
	}
	result := orderedmap.New()
	for name, value := range datas {
		fieldDesc := desc.FieldsByName[name]
		if fieldDesc == nil {
			result.Set(name, value)
			continue
		}
		field, err := FormatType(fieldDesc.Type, value)
		if err != nil {
			return nil, err
		}
		result.Set(name, field)
	}
	result.Sort(func(a *orderedmap.Pair, b *orderedmap.Pair) bool {
		aDesc := desc.FieldsByName[a.Key()]
		bDesc := desc.FieldsByName[b.Key()]
		if aDesc == nil {
			return false
		}
		if bDesc == nil {
			return true
		}
		return aDesc.ID < bDesc.ID
	})
	return result, nil
}

func FormatType(desc *descriptor.TypeDescriptor, data interface{}) (interface{}, error) {
	switch desc.Type {
	case descriptor.MAP:
		kv, isOk := data.(map[string]interface{})
		if !isOk {
			return data, nil
		}
		for key, value := range kv {
			vv, err := FormatType(desc.Elem, value)
			if err != nil {
				return nil, err
			}
			kv[key] = vv
		}
		return kv, nil
	case descriptor.LIST, descriptor.SET:
		datas, isOK := data.([]interface{})
		if !isOK {
			return datas, nil
		}
		for index := range datas {
			vv, err := FormatType(desc.Elem, datas[index])
			if err != nil {
				return nil, err
			}
			datas[index] = vv
		}
		return datas, nil
	case descriptor.STRUCT:
		return formatStruct(desc.Struct, data)
	default:
		return data, nil
	}
}

func FormatResponse(idlDesc *descriptor.ServiceDescriptor, method string, input string) (string, error) {
	function := idlDesc.Functions[method]
	if function == nil {
		return "", fmt.Errorf(`not found rpc method: %s`, method)
	}
	// oneway void
	if function.Response == nil {
		return input, nil
	}
	responseType, err := GetUnWrapperResponseType(function)
	if err != nil {
		return "", err
	}
	if !utils.Contains([]descriptor.Type{descriptor.SET, descriptor.LIST, descriptor.MAP, descriptor.STRUCT}, responseType.Type) {
		if responseType.Type == descriptor.VOID {
			return "", nil
		}
		return input, nil
	}
	var payload interface{}
	if err := utils.UnmarshalJsonUseNumber(input, &payload); err != nil {
		return "", err
	}
	result, err := FormatType(responseType, payload)
	if err != nil {
		return "", err
	}
	marshal, err := json.Marshal(result)
	if err != nil {
		return "", err
	}
	return utils.Bytes2String(marshal), nil
}

func GetUnWrapperResponseType(function *descriptor.FunctionDescriptor) (*descriptor.TypeDescriptor, error) {
	if function == nil || function.Response == nil || function.Response.Struct == nil || len(function.Response.Struct.FieldsByID) == 0 {
		return nil, fmt.Errorf(`invalid response type`)
	}
	return function.Response.Struct.FieldsByID[0].Type, nil
}
