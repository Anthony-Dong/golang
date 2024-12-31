package utils

import (
	"testing"
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
	t.Log(ToJson(structMap, true))

	structMap.Set("age", 2)
	t.Log(ToJson(structMap, true))
}
