package handler

import (
	"context"
	"testing"

	"github.com/anthony-dong/golang/pkg/rpc/kitex_demo/kitex_gen/api"
	"github.com/anthony-dong/golang/pkg/rpc/kitex_demo/kitex_gen/api/apiservice"
	"github.com/anthony-dong/golang/pkg/utils"
	"github.com/cloudwego/kitex/client"
)

func TestClient(t *testing.T) {
	server, err := NewServer(":10011")
	if err != nil {
		t.Fatal(err)
	}
	go server.Run()
	defer server.Stop()

	cc := apiservice.MustNewClient("a.b.c", client.WithHostPorts("localhost:10011"))
	err = cc.TestVoid(context.Background(), &api.Request{Field1: utils.StringPtr("1111")})
	if err != nil {
		t.Fatal(err)
	}
}
