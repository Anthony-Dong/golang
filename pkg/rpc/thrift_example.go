package rpc

import (
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

func GetExampleValue(tType *descriptor.TypeDescriptor, walk map[*descriptor.StructDescriptor]bool, op *Option) interface{} {
	if walk == nil {
		walk = map[*descriptor.StructDescriptor]bool{}
	}
	if basicTypeMap[tType.Type] {
		return op.Instance(tType.Type)
	}
	switch tType.Type {
	case descriptor.LIST:
		result := make([]interface{}, 0, op.GetListSize())
		for x := 0; x < op.GetListSize(); x++ {
			vv := GetExampleValue(tType.Elem, walk, op)
			if vv == nil {
				continue
			}
			result = append(result, vv)
		}
		return result
	case descriptor.MAP:
		result := make(map[string]interface{}, op.GetMapSize())
		for x := 0; x < op.GetMapSize(); x++ {
			if !basicTypeMap[tType.Key.Type] {
				continue
			}
			kv := GetExampleValue(tType.Key, walk, op)
			if kv == nil {
				continue
			}
			vv := GetExampleValue(tType.Elem, walk, op)
			result[utils.ToString(kv)] = vv
		}
		return result
	case descriptor.STRUCT:
		if walk[tType.Struct] {
			return nil
		}
		walk[tType.Struct] = true
		kv := orderedmap.New()
		fields := NewFields(tType.Struct)
		for _, elem := range fields {
			name := elem.TType.FieldName()
			kv.Set(name, GetExampleValue(elem.TType.Type, walk, op))
		}
		delete(walk, tType.Struct)
		return kv
	default:
		panic(`not support type`)
	}
}

type Option struct {
	Generator
}

func (*Option) GetListSize() int {
	return 1
}

func (*Option) GetMapSize() int {
	return 1
}

type Generator interface {
	Instance(t descriptor.Type) interface{}
}

type fixedGenerator struct{}

func NewFixedGenerator() *fixedGenerator {
	return &fixedGenerator{}
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

func (f *fixedGenerator) Instance(t descriptor.Type) interface{} {
	if i, ok := fixedTypeInstanceMap[t]; ok {
		return i
	}
	return nil
}
