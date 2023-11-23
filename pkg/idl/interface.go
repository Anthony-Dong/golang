package idl

import (
	"github.com/cloudwego/kitex/pkg/generic"
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

type DescriptorProvider interface {
	DescriptorProvider() (generic.DescriptorProvider, error)
}
