package idl

import "github.com/cloudwego/kitex/pkg/generic"

var _ DescriptorProvider = (*defaultDescriptorProvider)(nil)

type defaultDescriptorProvider struct {
	provider MemoryIDLProvider
}

func NewDescriptorProvider(provider MemoryIDLProvider) *defaultDescriptorProvider {
	return &defaultDescriptorProvider{provider: provider}
}

func (m *defaultDescriptorProvider) DescriptorProvider() (generic.DescriptorProvider, error) {
	idl, err := m.provider.MemoryIDL()
	if err != nil {
		return nil, err
	}
	//if provider, err := loadThriftDescriptorProvider(idl.Main, idl.IDLs); err == nil {
	//	return provider, nil
	//}
	return loadThriftDescriptorProvider(idl.Main, fixThriftIDLForKitex(idl.IDLs))
}
