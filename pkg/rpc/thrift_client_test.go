package rpc

import (
	"context"
	"path/filepath"
	"testing"

	kitex_client "github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/transport"

	"github.com/anthony-dong/golang/pkg/rpc/kitex_demo/handler"
	"github.com/anthony-dong/golang/pkg/utils"
)

func GetTestIDLPath() string {
	return filepath.Join(utils.GetGoProjectDir(), "pkg/idl/test/api.thrift")
}

func TestThriftClient(t *testing.T) {
	addr := ":10086"
	server, err := handler.NewServer(addr)
	if err != nil {
		t.Fatal(err)
	}
	go server.Run()
	defer server.Stop()
	runTestClient(t, "a.b.c", addr)
}

func runTestClient(t *testing.T, serviceName, addr string) {
	ctx := context.Background()
	client, err := NewThriftClient(NewLocalIDLProvider(map[string]string{
		serviceName: GetTestIDLPath(),
	}))
	if err != nil {
		t.Fatal(err)
	}
	methods := []string{"TestStruct", "TestVoid", "TestOnewayVoid", "TestList", "TestSet", "TestMap", "TestIntMap", "TestString"}
	for _, method := range methods {
		t.Run(method, func(t *testing.T) {
			exampleCode, err := client.GetExampleCode(ctx, serviceName, nil, method)
			if err != nil {
				t.Fatal(err)
			}
			resp, err := client.Do(ctx, &Request{RPCMethod: method, ServiceName: serviceName, Addr: addr, Body: exampleCode, Header: []*KV{
				NewKV("k1", "v1"),
				NewPersistentHeader("k2", "v2"),
				NewTransientHeader("k11", "1"),
				NewTransientHeader("k22", "2"),
			}}, kitex_client.WithTransportProtocol(transport.TTHeaderFramed))
			if err != nil {
				t.Fatal(err)
			}
			t.Logf("method: %s, resp: %s\n", method, string(resp.Body))
		})
	}
}
