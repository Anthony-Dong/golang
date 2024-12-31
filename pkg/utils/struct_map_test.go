package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type TestStruct struct {
	Name string `json:"name,omitempty"`
	Age  int    `json:"age,omitempty"`

	Num  int `json:"-"`
	Data string
}

func TestNewStructMap(t *testing.T) {
	tt := TestStruct{
		Name: "1",
		Age:  1,
		Num:  3,
		Data: "4",
	}
	structMap, err := NewStructMap(tt, "json")
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, ToJson(structMap), `{"name":"1","age":1,"Data":"4"}`)
	structMap.Set("age", 2)
	assert.Equal(t, ToJson(structMap), `{"name":"1","age":2,"Data":"4"}`)
}
