package idl

import (
	"path/filepath"
	"testing"

	"github.com/anthony-dong/golang/pkg/utils"
)

func GetTestMemoryIDL() MemoryIDLProvider {
	mainIdl := filepath.Join(utils.GetGoProjectDir(), "pkg/idl/test/api.thrift")
	return NewMemoryIDLProvider(mainIdl)
}

func TestNewLocalIDLProvider(t *testing.T) {
	desc, err := ParseThriftIDLFromProvider(GetTestMemoryIDL())
	if err != nil {
		t.Fatal(err)
	}
	provide := <-NewDescriptorProvider(desc).Provide()
	for _, function := range provide.Functions {
		t.Logf("func: %s\n", function.Name)
	}
}
