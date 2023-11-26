package idl

import (
	"path/filepath"
	"testing"

	"github.com/anthony-dong/golang/pkg/utils"
)

func TestNewLocalIDLProvider(t *testing.T) {
	mainIdl := filepath.Join(utils.GetGoProjectDir(), "pkg/idl/test/api.thrift")
	provider := NewLocalIDLProvider(mainIdl)
	idl, err := provider.MemoryIDL()
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("main: %s\n", idl.Main)
	for file := range idl.IDLs {
		t.Logf("file: %s\n", file)
	}

	thriftIDL, err := provider.ThriftIDL()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(thriftIDL.Filename)

	descriptorProvider, err := provider.DescriptorProvider()
	if err != nil {
		t.Fatal(err)
	}
	provide := <-descriptorProvider.Provide()

	for _, function := range provide.Functions {
		t.Logf("func: %s\n", function.Name)
	}
}
