package idl

import (
	"sync"

	"github.com/cloudwego/kitex/pkg/generic"
	"github.com/cloudwego/kitex/pkg/generic/descriptor"
)

var _ generic.DescriptorProvider = (*defaultDescriptorProvider)(nil)

func NewDescriptorProvider(desc *descriptor.ServiceDescriptor) *defaultDescriptorProvider {
	result := &defaultDescriptorProvider{
		desc: make(chan *descriptor.ServiceDescriptor, 1),
	}
	result.desc <- desc
	return result
}

type defaultDescriptorProvider struct {
	desc  chan *descriptor.ServiceDescriptor
	close sync.Once
}

func (c *defaultDescriptorProvider) Close() error {
	c.close.Do(func() {
		close(c.desc)
	})
	return nil
}

func (c *defaultDescriptorProvider) Provide() <-chan *descriptor.ServiceDescriptor {
	return c.desc
}
