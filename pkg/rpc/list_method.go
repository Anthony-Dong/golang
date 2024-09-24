package rpc

import (
	"context"

	"github.com/cloudwego/kitex/pkg/generic/descriptor"
)

func ListMethods(ctx context.Context, desc *descriptor.ServiceDescriptor) ([]string, error) {
	result := make([]string, 0, len(desc.Functions))
	for _, elem := range desc.Functions {
		result = append(result, elem.Name)
	}
	return result, nil
}
