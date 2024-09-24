package rpc

import (
	"fmt"
	"sort"

	"github.com/cloudwego/kitex/pkg/generic/descriptor"
	"github.com/iancoleman/orderedmap"

	"github.com/anthony-dong/golang/pkg/utils"
)

type Field struct {
	ID    int32
	TType *descriptor.FieldDescriptor
}

func NewFields(sStruct *descriptor.StructDescriptor) []*Field {
	fields := make([]*Field, 0, len(sStruct.FieldsByID))
	for id, value := range sStruct.FieldsByID {
		fields = append(fields, &Field{
			ID:    id,
			TType: value,
		})
	}
	sort.SliceStable(fields, func(i, j int) bool {
		return fields[i].ID < fields[j].ID
	})
	return fields
}

var basicTypeMap = map[descriptor.Type]bool{
	descriptor.BOOL:   true,
	descriptor.BYTE:   true,
	descriptor.DOUBLE: true,
	descriptor.I16:    true,
	descriptor.I32:    true,
	descriptor.I64:    true,
	descriptor.STRING: true,
}

func GetThriftExampleValue(tType *descriptor.TypeDescriptor, walk map[*descriptor.StructDescriptor]bool, op *ThriftExampleOption) (interface{}, error) {
	if walk == nil {
		walk = map[*descriptor.StructDescriptor]bool{}
	}
	if op == nil {
		op = NewThriftExampleOption()
	}
	if basicTypeMap[tType.Type] {
		return op.Generator.Instance(tType.Type), nil
	}
	switch tType.Type {
	case descriptor.LIST, descriptor.SET:
		result := make([]interface{}, 0, op.GetListSize())
		for x := 0; x < op.GetListSize(); x++ {
			vv, err := GetThriftExampleValue(tType.Elem, walk, op)
			if err != nil {
				return nil, err
			}
			if vv == nil {
				continue
			}
			result = append(result, vv)
		}
		return result, nil
	case descriptor.MAP:
		result := make(map[string]interface{}, op.GetMapSize())
		for x := 0; x < op.GetMapSize(); x++ {
			if !basicTypeMap[tType.Key.Type] {
				continue
			}
			kv, err := GetThriftExampleValue(tType.Key, walk, op)
			if err != nil {
				return nil, err
			}
			if kv == nil {
				continue
			}
			vv, err := GetThriftExampleValue(tType.Elem, walk, op)
			if err != nil {
				return nil, err
			}
			result[utils.ToString(kv)] = vv
		}
		return result, nil
	case descriptor.STRUCT:
		if walk[tType.Struct] {
			return map[string]interface{}{}, nil
		}
		walk[tType.Struct] = true
		kv := orderedmap.New()
		fields := NewFields(tType.Struct)
		for _, elem := range fields {
			name := elem.TType.FieldName()
			value, err := GetThriftExampleValue(elem.TType.Type, walk, op)
			if err != nil {
				return nil, err
			}
			kv.Set(name, value)
		}
		delete(walk, tType.Struct)
		return kv, nil
	case descriptor.VOID:
		return nil, nil
	default:
		return nil, fmt.Errorf(`not support thrift type: %v`, tType.Type)
	}
}

type ThriftExampleOption struct {
	Generator ThriftGenerator
}

func NewThriftExampleOption() *ThriftExampleOption {
	return &ThriftExampleOption{Generator: NewFixedGenerator()}
}

func (*ThriftExampleOption) GetListSize() int {
	return 1
}

func (*ThriftExampleOption) GetMapSize() int {
	return 1
}

type ThriftGenerator interface {
	Instance(t descriptor.Type) interface{}
}

type fixedThriftGenerator struct{}

func NewFixedGenerator() *fixedThriftGenerator {
	return &fixedThriftGenerator{}
}

var fixedTypeInstanceMap = map[descriptor.Type]interface{}{
	descriptor.BOOL:   false,
	descriptor.BYTE:   byte(1),
	descriptor.DOUBLE: float64(0),
	descriptor.I16:    int16(0),
	descriptor.I32:    int32(0),
	descriptor.I64:    int64(0),
	descriptor.STRING: "",
}

func (f *fixedThriftGenerator) Instance(t descriptor.Type) interface{} {
	if i, ok := fixedTypeInstanceMap[t]; ok {
		return i
	}
	return nil
}
