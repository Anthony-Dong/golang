package rpc

import (
	"context"
	"path/filepath"
	"testing"

	"github.com/anthony-dong/golang/pkg/utils"
)

func Test_ThriftClient(t *testing.T) {
	return
	client := NewThriftClient()
	ctx := context.Background()
	req := &Request{
		Service:  "xxx.xxx.xxx",
		Method:   "RPCAPI1",
		IDLType:  IDLTypeLocal,
		MainIDL:  filepath.Join("", "pkg/idl/test/thrift/rpc/rpc_api.thrift"),
		Body:     []byte(`{}`),
		Instance: Instance{
			//Host: "10.37.10.60:8888",
		},
	}
	code, err := client.ExampleCode(ctx, req)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(code)

	send, err := client.Send(ctx, req)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(utils.ToJson(send, true))
}
