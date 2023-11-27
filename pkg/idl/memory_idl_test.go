package idl

import (
	"path/filepath"
	"testing"

	"github.com/anthony-dong/golang/pkg/utils"
)

func TestNewLocalIDLProvider(t *testing.T) {
	mainIdl := filepath.Join(utils.GetGoProjectDir(), "pkg/idl/test/api.thrift")
	provider := NewDescriptorProvider(NewMemoryIDLProvider(mainIdl))
	descriptorProvider, err := provider.DescriptorProvider()
	if err != nil {
		t.Fatal(err)
	}
	provide := <-descriptorProvider.Provide()

	for _, function := range provide.Functions {
		t.Logf("func: %s\n", function.Name)
	}
}
