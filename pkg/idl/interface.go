package idl

import (
	"github.com/cloudwego/kitex/pkg/generic/descriptor"
	"github.com/cloudwego/thriftgo/parser"
)

type MemoryIDL struct {
	IDLs    map[string]string
	Main    string
	Include []string
}

type MemoryIDLProvider interface {
	MemoryIDL() (*MemoryIDL, error)
}

type ThriftIDLProvider interface {
	ThriftIDL() (*parser.Thrift, error)
}

func ParseThriftIDL(idl *MemoryIDL) (*descriptor.ServiceDescriptor, error) {
	provider, err := loadThriftDescriptorProvider(idl.Main, fixThriftIDLForKitex(idl.IDLs))
	if err != nil {
		return nil, err
	}
	return <-provider.Provide(), nil
}

func ParseThriftIDLFromProvider(idlProvider MemoryIDLProvider) (*descriptor.ServiceDescriptor, error) {
	idl, err := idlProvider.MemoryIDL()
	if err != nil {
		return nil, err
	}
	return ParseThriftIDL(idl)
}
