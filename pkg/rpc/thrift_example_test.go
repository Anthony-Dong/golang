package rpc

import (
	"path/filepath"
	"testing"

	"github.com/anthony-dong/golang/pkg/idl"
	"github.com/anthony-dong/golang/pkg/utils"
)

func GetTestMemoryIDL() idl.MemoryIDLProvider {
	mainIdl := filepath.Join(utils.GetGoProjectDir(), "pkg/idl/test/api.thrift")
	return idl.NewMemoryIDLProvider(mainIdl)
}

func TestGetThriftExampleValue(t *testing.T) {
	descriptor, err := idl.ParseThriftIDLFromProvider(GetTestMemoryIDL())
	if err != nil {
		t.Fatal(err)
	}
	for _, function := range descriptor.Functions {
		value, err := GetThriftExampleValue(function.Response, nil, nil)
		if err != nil {
			t.Fatal(err)
		}
		t.Logf("method: %s, example: %s\n", function.Name, utils.ToJson(value, true))
	}
}
