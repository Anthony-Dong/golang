package utils

import (
	"reflect"
	"runtime"
)

func FunctionName(Func interface{}) string {
	if Func == nil {
		return ""
	}
	value := reflect.ValueOf(Func)
	if value.Kind() != reflect.Func {
		return ""
	}
	return runtime.FuncForPC(value.Pointer()).Name()
}
