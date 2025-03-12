package utils

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsBasicType(t *testing.T) {
	t.Run("bool", func(t *testing.T) {
		tt := reflect.TypeOf(true)
		assert.True(t, isBasicType(tt, tt))
	})
	t.Run("typedef", func(t *testing.T) {
		type A int
		assert.False(t, isBasicType(reflect.TypeOf(A(1)), reflect.TypeOf(int(1))))
	})
	t.Run("string", func(t *testing.T) {
		tt := reflect.TypeOf("")
		assert.True(t, isBasicType(tt, tt))
	})
	t.Run("int", func(t *testing.T) {
		tt := reflect.TypeOf(1)
		assert.True(t, isBasicType(tt, tt))
	})
	t.Run("float64", func(t *testing.T) {
		tt := reflect.TypeOf(1.11)
		assert.True(t, isBasicType(tt, tt))
	})
	t.Run("slice", func(t *testing.T) {
		tt := reflect.TypeOf([]string{})
		assert.True(t, isBasicType(tt, tt))
	})
	t.Run("map", func(t *testing.T) {
		tt := reflect.TypeOf(map[string]string{})
		assert.True(t, isBasicType(tt, tt))
	})
	t.Run("map ptr", func(t *testing.T) {
		tt := reflect.TypeOf(map[string]*string{})
		assert.True(t, isBasicType(tt, tt))
	})
	t.Run("struct", func(t *testing.T) {
		type TestStruct struct {
			Name string
		}
		tt := reflect.TypeOf(TestStruct{})
		assert.False(t, isBasicType(tt, tt))
	})
	t.Run("ptr struct", func(t *testing.T) {
		type TestStruct struct {
			Name string
		}
		tt := reflect.TypeOf(TestStruct{})
		assert.False(t, isBasicType(tt, tt))
	})
	t.Run("elem struct", func(t *testing.T) {
		type TestStruct struct {
			Name string
		}
		tt := reflect.TypeOf([]*TestStruct{})
		assert.False(t, isBasicType(tt, tt))
	})

	t.Run("map struct", func(t *testing.T) {
		type TestStruct struct {
			Name string
		}
		tt := reflect.TypeOf(map[string]*TestStruct{})
		assert.False(t, isBasicType(tt, tt))
	})
	t.Run("cycle struct", func(t *testing.T) {
		type TestStruct struct {
			Name string
			Next *TestStruct
		}
		tt := reflect.TypeOf([]*TestStruct{})
		assert.False(t, isBasicType(tt, tt))
	})
}

func Benchmark_unsafeValueSet(b *testing.B) {
	src := reflect.ValueOf(1.111)
	dst := reflect.New(src.Type()).Elem()
	for i := 0; i < b.N; i++ {
		if err := unsafeValueSet(src, dst); err != nil {
			b.Fatal(err)
		}
	}
}

func Benchmark_unsafeValueSet_reflect(b *testing.B) {
	src := reflect.ValueOf(1.111)
	dst := reflect.New(src.Type()).Elem()
	for i := 0; i < b.N; i++ {
		dst.Set(src)
	}
}
