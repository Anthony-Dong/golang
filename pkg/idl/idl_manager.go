package idl

import "github.com/cloudwego/kitex/pkg/generic/descriptor"

type ThriftIDLManager interface {
	GetThriftIDL(service string) (*descriptor.ServiceDescriptor, error)
}
