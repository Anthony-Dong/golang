package thrift_codec

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/apache/thrift/lib/go/thrift"
	"github.com/cloudwego/kitex/pkg/generic/descriptor"
	"github.com/iancoleman/orderedmap"

	"github.com/anthony-dong/golang/pkg/utils"
)

func EncodeMessage(oprot thrift.TProtocol, desc *descriptor.TypeDescriptor, data interface{}) error {
	return encodeThriftType(oprot, "", desc, data)
}

func EncodeReply(oprot thrift.TProtocol, function *descriptor.FunctionDescriptor, seq int32, data interface{}) error {
	if err := oprot.WriteMessageBegin(function.Name, thrift.REPLY, seq); err != nil {
		return err
	}
	if err := EncodeMessage(oprot, function.Response, data); err != nil {
		return err
	}
	if err := oprot.WriteMessageEnd(); err != nil {
		return err
	}
	return nil
}

func safeToOrderMap(data interface{}) (*orderedmap.OrderedMap, bool) {
	switch v := data.(type) {
	case map[string]interface{}:
		kv := orderedmap.New()
		for key, value := range v {
			kv.Set(key, value)
		}
		return kv, true
	case *orderedmap.OrderedMap:
		return v, true
	case orderedmap.OrderedMap:
		return &v, true
	default:
		return nil, false
	}
}

func toThriftType(p descriptor.Type) thrift.TType {
	return thrift.TType(p.ToThriftTType())
}

func encodeThriftType(iprot thrift.TProtocol, fieldName string, tType *descriptor.TypeDescriptor, data interface{}) error {
	if data == nil {
		return nil
	}
	switch tType.Type {
	case descriptor.STRUCT:
		kv, isOK := safeToOrderMap(data)
		if !isOK {
			return fmt.Errorf(`invalid type: %T`, data)
		}
		if err := iprot.WriteStructBegin(fieldName); err != nil {
			return err
		}
		for _, k := range kv.Keys() {
			v, _ := kv.Get(k)
			desc := tType.Struct.FieldsByName[k]
			if desc == nil {
				return fmt.Errorf(`not found field name %s`, k)
			}
			fieldType := toThriftType(desc.Type.Type)
			// if is void resp. skip write field
			if fieldType == thrift.VOID {
				continue
			}
			if err := iprot.WriteFieldBegin(desc.Name, fieldType, int16(desc.ID)); err != nil {
				return err
			}
			if err := encodeThriftType(iprot, desc.Name, desc.Type, v); err != nil {
				return err
			}
			if err := iprot.WriteFieldEnd(); err != nil {
				return err
			}
		}
		if err := iprot.WriteFieldStop(); err != nil {
			return err
		}
		if err := iprot.WriteStructEnd(); err != nil {
			return err
		}
	case descriptor.LIST, descriptor.SET:
		datas, isOK := data.([]interface{})
		if !isOK {
			return fmt.Errorf(`invalid type: %T`, data)
		}
		if err := iprot.WriteListBegin(toThriftType(tType.Elem.Type), len(datas)); err != nil {
			return err
		}
		for _, v := range datas {
			if err := encodeThriftType(iprot, fieldName, tType.Elem, v); err != nil {
				return err
			}
		}
		if err := iprot.WriteListEnd(); err != nil {
			return err
		}
	case descriptor.MAP:
		kv, isOK := data.(map[string]interface{})
		if !isOK {
			return fmt.Errorf(`invalid type: %T`, data)
		}
		if err := iprot.WriteMapBegin(toThriftType(tType.Key.Type), toThriftType(tType.Elem.Type), len(kv)); err != nil {
			return err
		}
		for k, v := range kv {
			if err := encodeThriftType(iprot, fieldName, tType.Key, k); err != nil {
				return err
			}
			if err := encodeThriftType(iprot, fieldName, tType.Elem, v); err != nil {
				return err
			}
		}
		if err := iprot.WriteMapEnd(); err != nil {
			return err
		}
	case descriptor.STRING:
		if tType.Name == "binary" {
			switch v := data.(type) {
			case []byte:
				if err := iprot.WriteBinary(v); err != nil {
					return err
				}
			default:
				if err := iprot.WriteBinary(utils.String2Bytes(utils.ToString(data))); err != nil {
					return err
				}
			}
		} else {
			if err := iprot.WriteString(utils.ToString(data)); err != nil {
				return err
			}
		}
	case descriptor.I08:
		i8, err := toInt64(data, 8)
		if err != nil {
			return err
		}
		if err := iprot.WriteByte(int8(i8)); err != nil {
			return err
		}
	case descriptor.I16:
		i16, err := toInt64(data, 16)
		if err != nil {
			return err
		}
		if err := iprot.WriteI16(int16(i16)); err != nil {
			return err
		}
	case descriptor.I32:
		i32, err := toInt64(data, 32)
		if err != nil {
			return err
		}
		if err := iprot.WriteI32(int32(i32)); err != nil {
			return err
		}
	case descriptor.I64:
		i64, err := toInt64(data, 64)
		if err != nil {
			return err
		}
		if err := iprot.WriteI64(i64); err != nil {
			return err
		}
	case descriptor.BOOL:
		value := utils.ToString(data)
		parseBool, err := strconv.ParseBool(value)
		if err != nil {
			return err
		}
		if err := iprot.WriteBool(parseBool); err != nil {
			return err
		}
	case descriptor.VOID:
		return nil
	default:
		return fmt.Errorf(`unsupported type: %s`, tType.Type)
	}
	return nil
}

func toInt64(data interface{}, size int) (int64, error) {
	switch v := data.(type) {
	case json.Number:
		return v.Int64()
	case float64:
		return int64(v), nil
	case int64:
		return v, nil
	}
	return strconv.ParseInt(fmt.Sprintf("%v", data), 10, size)
}
