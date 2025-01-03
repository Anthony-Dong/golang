package thrift_codec

import (
	"context"
	"encoding/base64"
	"unicode/utf8"

	"github.com/anthony-dong/golang/pkg/utils"

	"github.com/anthony-dong/golang/pkg/codec/thrift_codec/kitex"

	"github.com/anthony-dong/golang/pkg/logs"

	"github.com/apache/thrift/lib/go/thrift"
)

type ThriftException struct {
	TypeId  int32  `json:"type_id"`
	Message string `json:"message"`
	//Exception thrift.TApplicationException `json:"-"`
}

type ThriftMessage struct {
	Method      string             `json:"method"`
	SeqId       int32              `json:"seq_id"`
	Protocol    Protocol           `json:"protocol"` // set protocol
	MessageType ThriftTMessageType `json:"message_type"`
	Transport   *TransportInfo     `json:"transport,omitempty"`
	Payload     *FieldOrderMap     `json:"payload,omitempty"`
	Exception   *ThriftException   `json:"exception,omitempty"` // MessageType=EXCEPTION 存在异常则是这个字段
	MetaInfo    *kitex.MetaInfo    `json:"meta_info,omitempty"`
}

type TransportInfo struct {
	FromService string `json:"from_service,omitempty"`
	ToService   string `json:"to_service,omitempty"`
	FromAddr    string `json:"from_addr,omitempty"`
	ToAddr      string `json:"to_addr,omitempty"`
}

type ThriftTMessageType thrift.TMessageType

const (
	InvalidTMessageType ThriftTMessageType = 0
	CALL                ThriftTMessageType = 1
	REPLY               ThriftTMessageType = 2
	EXCEPTION           ThriftTMessageType = 3
	ONEWAY              ThriftTMessageType = 4
)

func DecodeMessage(ctx context.Context, iprot thrift.TProtocol) (*ThriftMessage, error) {
	name, messageType, seqId, err := iprot.ReadMessageBegin()
	if err != nil {
		return nil, err
	}
	result := &ThriftMessage{
		Method:      name,
		SeqId:       seqId,
		MessageType: ThriftTMessageType(messageType),
	}
	switch messageType {
	case thrift.EXCEPTION:
		exception := thrift.NewTApplicationException(thrift.UNKNOWN_APPLICATION_EXCEPTION, "Unknown Exception")
		if err := exception.Read(iprot); err != nil {
			return nil, err
		}
		result.Exception = &ThriftException{
			Message: exception.Error(),
			TypeId:  exception.TypeId(),
		}
	case thrift.REPLY, thrift.CALL, thrift.ONEWAY:
		decodeStruct, err := DecodeStruct(ctx, iprot)
		if err != nil {
			return nil, err
		}
		result.Payload = decodeStruct
	case thrift.INVALID_TMESSAGE_TYPE:
		logs.CtxInfo(ctx, "[DecodeRespMessage] not handler message type: %v, method_name: %s", messageType, name)
	}
	if err := iprot.ReadMessageEnd(); err != nil {
		return nil, err
	}
	return result, nil
}

// DecodeStruct
// thriftStruct 指的是msg的类型
// fileName 指的是thriftStruct所在的文件，必须是 idl map[string]*parser.Thrift.
func DecodeStruct(ctx context.Context, iprot thrift.TProtocol) (*FieldOrderMap, error) {
	if _, err := iprot.ReadStructBegin(); err != nil {
		return nil, err
	}
	result := NewFieldOrderMap(16)
	for {
		_, fieldType, fieldId, err := iprot.ReadFieldBegin()
		if err != nil {
			return nil, err
		}
		if fieldType == thrift.STOP {
			break
		}
		fieldValue, err := DecodeField(ctx, fieldType, iprot)
		if err != nil {
			return nil, err
		}
		result.Set(NewField(fieldId, fieldType), fieldValue)
		if err := iprot.ReadFieldEnd(); err != nil {
			return nil, err
		}
	}
	if err := iprot.ReadStructEnd(); err != nil {
		return nil, err
	}
	return result, nil
}

func isValidUTF8(s string) bool {
	return utf8.ValidString(s)
}

func DecodeField(ctx context.Context, fieldType thrift.TType, iprot thrift.TProtocol) (interface{}, error) {
	switch fieldType {
	case thrift.BOOL:
		return iprot.ReadBool()
	case thrift.DOUBLE:
		return iprot.ReadDouble()
	case thrift.I08: // 或者 BYTE
		readByte, err := iprot.ReadByte() // 返回的是 type byte = uint8 类型，需要强制转换一下!
		if err != nil {
			return nil, err
		}
		return readByte, nil
	case thrift.I16:
		return iprot.ReadI16()
	case thrift.I32:
		return iprot.ReadI32()
	case thrift.I64:
		return iprot.ReadI64()
	case thrift.STRING:
		str, err := iprot.ReadString()
		if err != nil {
			return nil, err
		}
		if isValidUTF8(str) {
			return str, nil
		}
		return base64.StdEncoding.EncodeToString([]byte(str)), nil
	//case thrift.BINARY: // not support! 这里可以skip掉
	//	return iprot.ReadBinary()
	case thrift.MAP:
		keyType, valueType, size, err := iprot.ReadMapBegin()
		if err != nil {
			return nil, err
		}
		result := make(map[string]interface{}, size)
		for i := 0; i < size; i++ {
			var (
				err   error
				key   interface{}
				value interface{}
			)
			if key, err = DecodeField(ctx, keyType, iprot); err != nil {
				return nil, err
			}
			if value, err = DecodeField(ctx, valueType, iprot); err != nil {
				return nil, err
			}
			if key == nil { // key为空skip掉!
				continue
			}
			result[utils.ToString(key)] = value
		}
		if err := iprot.ReadMapEnd(); err != nil {
			return nil, err
		}
		return result, nil
	case thrift.SET:
		elemType, size, err := iprot.ReadSetBegin()
		if err != nil {
			return nil, err
		}
		result := make([]interface{}, 0, size)
		for i := 0; i < size; i++ {
			if elem, err := DecodeField(ctx, elemType, iprot); err != nil {
				return nil, err
			} else if elem != nil {
				result = append(result, elem)
			}
		}
		if err := iprot.ReadSetEnd(); err != nil {
			return nil, err
		}
		return result, nil
	case thrift.LIST:
		elemType, size, err := iprot.ReadListBegin()
		if err != nil {
			return nil, err
		}
		result := make([]interface{}, 0, size)
		for i := 0; i < size; i++ {
			if elem, err := DecodeField(ctx, elemType, iprot); err != nil {
				return nil, err
			} else if elem != nil { // 不为空再append!
				result = append(result, elem)
			}
		}
		if err := iprot.ReadListEnd(); err != nil {
			return nil, err
		}
		return result, nil
	case thrift.STRUCT:
		return DecodeStruct(ctx, iprot)
	default:
		return nil, iprot.Skip(fieldType)
	}
}

type FieldOrderMap struct {
	list []Field
	data map[Field]interface{}
}

func (f *FieldOrderMap) Range(foo func(Field, interface{}) bool) {
	for _, elem := range f.list {
		if !foo(elem, f.data[elem]) {
			break
		}
	}
}

func (f *FieldOrderMap) RangeErr(foo func(Field, interface{}) error) error {
	for _, elem := range f.list {
		if err := foo(elem, f.data[elem]); err != nil {
			return err
		}
	}
	return nil
}

func NewFieldOrderMap(size int) *FieldOrderMap {
	return &FieldOrderMap{
		list: make([]Field, 0, size),
		data: make(map[Field]interface{}, size),
	}
}

func (f *FieldOrderMap) Set(field Field, v interface{}) {
	if f.data == nil {
		f.data = map[Field]interface{}{}
	}
	if _, isExist := f.data[field]; isExist {
		// not handler order!! swap...
		f.data[field] = v
		return
	}
	f.list = append(f.list, field)
	f.data[field] = v
}

type Field struct {
	FieldId   int16
	FieldType thrift.TType
}

func NewField(fieldId int16, fieldType thrift.TType) Field {
	return Field{FieldId: fieldId, FieldType: fieldType}
}
