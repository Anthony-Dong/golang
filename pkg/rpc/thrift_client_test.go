package rpc

import (
	"context"
	"path/filepath"
	"testing"

	"github.com/anthony-dong/golang/pkg/idl"

	"github.com/anthony-dong/golang/pkg/utils"
)

func Test_ThriftClient(t *testing.T) {
	return
	// local test
	mainIdl := filepath.Join(utils.GetGoProjectDir(), "pkg/idl/test/api.thrift")
	client := NewThriftClient(idl.NewLocalIDLProvider(mainIdl))
	ctx := context.Background()
	req := &Request{
		Service: "TestServiceName",
		Method:  "RPCAPI1",
		Body:    []byte(`{"Field1": "success"}`),
		Endpoint: Endpoint{
			Addr: "127.0.0.1:8888",
		},
	}
	example, err := client.ExampleCode(ctx, req)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("example code: %s\n", example)

	req.Body = []byte(example)
	send, err := client.Send(ctx, req)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(utils.ToJson(send, true))
}
