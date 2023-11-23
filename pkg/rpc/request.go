package rpc

import (
	"context"
	"fmt"

	"github.com/anthony-dong/golang/pkg/idl"
)

func (r *Request) NewIDLProvider(ctx context.Context) (idl.DescriptorProvider, error) {
	switch r.IDLType {
	case IDLTypeLocal:
		return idl.NewLocalIDLProvider(r.MainIDL), nil
	default:
		return nil, fmt.Errorf(`invalid idl type: %s`, r.IDLType)
	}
}
