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

type Data struct {
	Name string
}

func (d *Data) GetName() string {
	return d.Name
}

type GetNamer interface {
	GetName() string
}

type TestWrapper struct {
	Data GetNamer
}

func (w *TestWrapper) GetName() string {
	return w.Data.GetName()
}

func BenchmarkName(b *testing.B) {
	data := &Data{}
	for i := 0; i < b.N; i++ {
		wrapper := TestWrapper{Data: data}
		wrapper.GetName()
	}
}

// BenchmarkName-12    	1000000000	         0.2701 ns/op
// BenchmarkName-12    	1000000000	         0.2664 ns/op

// BenchmarkName-12    	746830246	         1.714 ns/op
