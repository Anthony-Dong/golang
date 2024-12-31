package utils

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/anthony-dong/golang/pkg/orderedmap"
)

func NewStructMap(t interface{}, tag string) (*orderedmap.OrderedMap, error) {
	value := reflect.ValueOf(t)
	for value.Kind() == reflect.Ptr {
		value = value.Elem()
	}

	if value.Kind() != reflect.Struct {
		return nil, fmt.Errorf("%s is not a struct type", value.Type())
	}
	if t == nil {
		return orderedmap.New(), nil
	}
	Type := value.Type()
	result := orderedmap.New()

	field := value.NumField()

	getTagValue := func(tagv string, defaultV string) string {
		if tagv == "-" {
			return ""
		}
		if tagv == "" {
			return defaultV
		}
		split := strings.Split(tagv, ",")
		return split[0]
	}

	for i := 0; i < field; i++ {
		fieldValue := value.Field(i)
		fieldType := Type.Field(i)
		fieldName := fieldType.Name
		if tag != "" {
			v, _ := fieldType.Tag.Lookup(tag)
			fieldName = getTagValue(v, fieldName)
			if fieldName == "" {
				continue
			}
		}
		result.Set(fieldName, fieldValue.Interface())
	}
	return result, nil
}
