package thrift_codec

import (
	"bytes"
	"context"
	"fmt"

	"github.com/cloudwego/kitex/pkg/generic/descriptor"

	"github.com/anthony-dong/golang/pkg/orderedmap"
)

func convertToOrderMap(data interface{}) (*orderedmap.OrderedMap, error) {
	switch v := data.(type) {
	case orderedmap.OrderedMap:
		return &v, nil
	case *orderedmap.OrderedMap:
		return v, nil
	case map[string]interface{}:
		result := orderedmap.New()
		for k, elem := range v {
			result.Set(k, elem)
		}
		return result, nil
	}
	return nil, fmt.Errorf("convert ordered map find err: unknown type %T", data)
}

func DecodeThriftData(payload []byte) (*FieldOrderMap, error) {
	protocol := NewTProtocol(bytes.NewBuffer(payload), UnframedBinary)
	return DecodeStruct(context.Background(), protocol)
}

func DecodeThriftDataToJson(method string, payload *FieldOrderMap, thriftIdl *descriptor.ServiceDescriptor, isReq bool) (interface{}, error) {
	function, err := lookupThriftFunc(thriftIdl, method)
	if err != nil {
		return nil, err
	}
	if isReq {
		return decodePayload(payload, function.Request.Struct.FieldsByName["req"].Type)
	}
	return decodePayload(payload, function.Response.Struct.FieldsByName[""].Type)
}

func DecodeThriftMessage(msg *ThriftMessage, thriftIdl *descriptor.ServiceDescriptor) (interface{}, error) {
	function, err := lookupThriftFunc(thriftIdl, msg.Method)
	if err != nil {
		return nil, err
	}
	switch msg.MessageType {
	case CALL:
		data, err := decodePayload(msg.Payload, function.Request)
		if err != nil {
			return nil, err
		}
		orderMap, err := convertToOrderMap(data)
		if err != nil {
			return nil, err
		}
		if v, exists := orderMap.Get("req"); exists {
			return v, nil
		}
		return orderMap, nil
	case REPLY:
		data, err := decodePayload(msg.Payload, function.Response)
		if err != nil {
			return nil, err
		}
		orderMap, err := convertToOrderMap(data)
		if err != nil {
			return nil, err
		}
		if v, exists := orderMap.Get(""); exists {
			return v, nil
		}
		return orderMap, nil
	case EXCEPTION:
		return msg.Exception, nil
	case ONEWAY:
		return msg.Payload, nil
	}
	return nil, fmt.Errorf("unknown message type: %s", msg.MessageType)
}

func decodePayload(data interface{}, desc *descriptor.TypeDescriptor) (interface{}, error) {
	switch desc.Type {
	case descriptor.STRUCT:
		orderMap, err := ToFieldOrderMap(data)
		if err != nil {
			return nil, err
		}
		result := orderedmap.New()
		err = orderMap.RangeErr(func(field Field, elem interface{}) error {
			fieldDescriptor := desc.Struct.FieldsByID[int32(field.FieldId)]
			if fieldDescriptor == nil {
				result.Set(field.String(), elem)
				return nil
			}
			value, err := decodePayload(elem, fieldDescriptor.Type)
			if err != nil {
				return err
			}
			result.Set(fieldDescriptor.Name, value)
			return nil
		})
		if err != nil {
			return nil, err
		}
		return result, nil
	case descriptor.LIST, descriptor.SET:
		list, _ := data.([]interface{})
		for index, elem := range list {
			decoded, err := decodePayload(elem, desc.Elem)
			if err != nil {
				return nil, err
			}
			list[index] = decoded
		}
		return list, nil
	case descriptor.MAP:
		kv, _ := data.(map[string]interface{})
		for k, v := range kv {
			decodeValue, err := decodePayload(v, desc.Elem)
			if err != nil {
				return nil, err
			}
			kv[k] = decodeValue
		}
		return kv, nil
	}
	return data, nil
}

func lookupThriftFunc(idl *descriptor.ServiceDescriptor, method string) (*descriptor.FunctionDescriptor, error) {
	for _, function := range idl.Functions {
		if function.Name == method {
			return function, nil
		}
	}
	return nil, fmt.Errorf(`not found method "%s"`, method)
}
