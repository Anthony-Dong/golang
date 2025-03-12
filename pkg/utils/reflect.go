package utils

import (
	"fmt"
	"reflect"
	"unsafe"
)

// isBasicType 比较 src / dst 底层全部都是基本类型(内置类型)
func isBasicType(src reflect.Type, dst reflect.Type) bool {
	if src != dst {
		return false
	}
	switch src.Kind() {
	case reflect.String, reflect.Bool:
		return true
	case reflect.Float32, reflect.Float64:
		return true
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return true
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return true
	case reflect.Slice:
		return isBasicType(src.Elem(), dst.Elem())
	case reflect.Map:
		return isBasicType(src.Key(), dst.Key()) && isBasicType(src.Elem(), dst.Elem())
	case reflect.Ptr:
		return isBasicType(src.Elem(), dst.Elem())
	default:
		return false
	}
}

func unsafeValueSet(srcValue reflect.Value, dstValue reflect.Value) error {
	if !dstValue.CanAddr() {
		return fmt.Errorf("dst must be addressable")
	}
	if !srcValue.CanInterface() {
		return fmt.Errorf("src must be interface")
	}
	src := srcValue.Interface()
	dstUnsafePtr := unsafe.Pointer(dstValue.UnsafeAddr())
	switch dstValue.Kind() {
	case reflect.String:
		*(*string)(dstUnsafePtr) = src.(string)
	case reflect.Int64:
		*(*int64)(dstUnsafePtr) = src.(int64)
	case reflect.Int32:
		*(*int32)(dstUnsafePtr) = src.(int32)
	case reflect.Int16:
		*(*int16)(dstUnsafePtr) = src.(int16)
	case reflect.Int8:
		*(*int8)(dstUnsafePtr) = src.(int8)
	case reflect.Int:
		*(*int)(dstUnsafePtr) = src.(int)
	case reflect.Uint64:
		*(*uint64)(dstUnsafePtr) = src.(uint64)
	case reflect.Uint32:
		*(*uint32)(dstUnsafePtr) = src.(uint32)
	case reflect.Uint16:
		*(*uint16)(dstUnsafePtr) = src.(uint16)
	case reflect.Uint8:
		*(*uint8)(dstUnsafePtr) = src.(uint8)
	case reflect.Uint:
		*(*uint)(dstUnsafePtr) = src.(uint)
	case reflect.Bool:
		*(*bool)(dstUnsafePtr) = src.(bool)
	case reflect.Float64:
		*(*float64)(dstUnsafePtr) = src.(float64)
	default:
		return fmt.Errorf("builtin type %s not supported", dstValue.Type())
	}
	return nil
}
