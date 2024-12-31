package thrift_codec

import (
	"bytes"
	"encoding/json"
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/apache/thrift/lib/go/thrift"
)

var protocolToString = map[Protocol]string{
	UnknownProto:             "Unknown",
	UnframedBinary:           "UnframedBinary",
	UnframedCompact:          "UnframedCompact",
	FramedBinary:             "FramedBinary",
	FramedCompact:            "FramedCompact",
	UnframedHeader:           "UnframedHeader",
	FramedHeader:             "FramedHeader",
	UnframedUnStrictBinary:   "UnframedUnStrictBinary",
	FramedUnStrictBinary:     "FramedUnStrictBinary",
	UnframedBinaryTTHeader:   "UnframedBinaryTTHeader",
	FramedBinaryTTHeader:     "FramedBinaryTTHeader",
	UnframedBinaryMeshHeader: "UnframedBinaryMeshHeader",
	FramedBinaryMeshHeader:   "FramedBinaryMeshHeader",
}

var stringToProtocol = map[string]Protocol{}

func init() {
	for k, v := range protocolToString {
		stringToProtocol[v] = k
	}
}

func (p Protocol) String() string {
	r, isOk := protocolToString[p]
	if isOk {
		return r
	}
	return protocolToString[UnknownProto]
}

func (p Protocol) MarshalJSON() (text []byte, err error) {
	return []byte(strconv.Quote(p.String())), nil
}

func (p *Protocol) UnmarshalJSON(text []byte) (err error) {
	if len(text) == 0 {
		*p = UnknownProto
		return nil
	}
	unquote, err := strconv.Unquote(string(text))
	if err != nil {
		return err
	}
	protocol, ok := stringToProtocol[unquote]
	if !ok {
		*p = UnknownProto
		return nil
	}
	*p = protocol
	return nil
}

var thriftTMessageTypeToString = map[ThriftTMessageType]string{
	InvalidTMessageType: "invalid",
	CALL:                "call",
	REPLY:               "reply",
	EXCEPTION:           "exception",
	ONEWAY:              "oneway",
}

var stringToThriftMessageType = map[string]ThriftTMessageType{}

func init() {
	for k, v := range thriftTMessageTypeToString {
		stringToThriftMessageType[v] = k
	}
}

func (p ThriftTMessageType) String() string {
	v, isExist := thriftTMessageTypeToString[p]
	if isExist {
		return v
	}
	return thriftTMessageTypeToString[InvalidTMessageType]
}

func (p *ThriftTMessageType) UnmarshalJSON(bytes []byte) error {
	if len(bytes) == 0 {
		*p = InvalidTMessageType
		return nil
	}
	unquote, err := strconv.Unquote(string(bytes))
	if err != nil {
		return err
	}
	messageType, isOk := stringToThriftMessageType[unquote]
	if !isOk {
		messageType = InvalidTMessageType
	}
	*p = messageType
	return nil
}

func (p ThriftTMessageType) MarshalJSON() (text []byte, err error) {
	return []byte(strconv.Quote(p.String())), nil
}

var typeNames = map[thrift.TType]string{
	thrift.STOP:   "STOP",
	thrift.VOID:   "VOID",
	thrift.BOOL:   "BOOL",
	thrift.BYTE:   "BYTE",
	thrift.DOUBLE: "DOUBLE",
	thrift.I16:    "I16",
	thrift.I32:    "I32",
	thrift.I64:    "I64",
	thrift.STRING: "STRING",
	thrift.STRUCT: "STRUCT",
	thrift.MAP:    "MAP",
	thrift.SET:    "SET",
	thrift.LIST:   "LIST",
	thrift.UTF8:   "UTF8",
	thrift.UTF16:  "UTF16",
}

var revertTypeNames = map[string]thrift.TType{}

func init() {
	for k, v := range typeNames {
		revertTypeNames[v] = k
	}
}

func (t Field) MarshalJSON() ([]byte, error) {
	return []byte(strconv.Quote(t.String())), nil
}

func (t *Field) UnmarshalJSON(b []byte) error {
	parts := strings.Split(string(b), "_")
	if len(parts) != 2 {
		return fmt.Errorf("invalid field string: %s", string(b))
	}
	id, err := strconv.ParseInt(parts[0], 10, 16)
	if err != nil {
		return fmt.Errorf("invalid fieldId: %s", parts[0])
	}
	t.FieldId = int16(id)
	fieldType, isOk := revertTypeNames[parts[1]]
	if !isOk {
		return fmt.Errorf("invalid fieldType: %s", parts[1])
	}
	t.FieldType = fieldType
	return nil
}

func (t Field) String() string {
	return fmt.Sprintf("%d_%s", t.FieldId, t.FieldType)
}

func (f *FieldOrderMap) String() string {
	marshal, err := json.Marshal(f)
	if err != nil {
		return fmt.Sprintf("%#v", f)
	}
	return string(marshal)
}

func (f FieldOrderMap) MarshalJSON() ([]byte, error) {
	sort.Slice(f.list, func(i, j int) bool {
		return f.list[i].FieldId < f.list[j].FieldId
	})
	result := bytes.Buffer{}
	result.WriteString("{")
	for index, v := range f.list {
		result.WriteByte('"')
		result.WriteString(v.String())
		result.WriteByte('"')
		result.WriteByte(':')
		marshal, err := json.Marshal(f.data[v])
		if err != nil {
			return nil, err
		}
		result.Write(marshal)
		if index == len(f.list)-1 {
			continue
		}
		result.WriteByte(',')
	}
	result.WriteByte('}')
	return result.Bytes(), nil
}

func (f *FieldOrderMap) UnmarshalJSON(b []byte) error {
	var temp map[string]interface{}
	if err := json.Unmarshal(b, &temp); err != nil {
		return err
	}
	f.list = make([]Field, 0, len(temp))
	f.data = make(map[Field]interface{}, len(temp))
	for key, value := range temp {
		var field Field
		if err := field.UnmarshalJSON([]byte(key)); err != nil {
			return err
		}
		f.list = append(f.list, field)
		f.data[field] = value
	}
	sort.Slice(f.list, func(i, j int) bool {
		return f.list[i].FieldId < f.list[j].FieldId
	})
	return nil
}

func ToFieldOrderMap(data interface{}) (*FieldOrderMap, error) {
	switch v := data.(type) {
	case FieldOrderMap:
		return &v, nil
	case *FieldOrderMap:
		return v, nil
	case map[string]interface{}:
		marshal, err := json.Marshal(data)
		if err != nil {
			return nil, err
		}
		result := NewFieldOrderMap(len(v))
		err = json.Unmarshal(marshal, result)
		if err != nil {
			return nil, err
		}
		return result, nil
	case nil:
		return &FieldOrderMap{}, nil
	}
	return nil, fmt.Errorf(`data type [%T] convert stauct failed`, data)
}
