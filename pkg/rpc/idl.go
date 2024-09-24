package rpc

import (
	"context"

	"github.com/anthony-dong/golang/pkg/idl"
)

type localIDLProvider struct {
	mains map[string]string
}

func (s *localIDLProvider) MemoryIDL(ctx context.Context, serviceName string, idlConfig *IDLConfig) (*idl.MemoryIDL, error) {
	return idl.NewMemoryIDLProvider(s.mains[serviceName]).MemoryIDL()
}

var _ IDLProvider = (*localIDLProvider)(nil)

func NewLocalIDLProvider(mains map[string]string) IDLProvider {
	return &localIDLProvider{mains: mains}
}
