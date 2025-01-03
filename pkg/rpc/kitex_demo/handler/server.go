package handler

import (
	"net"

	"github.com/cloudwego/kitex/server"

	api "github.com/anthony-dong/golang/pkg/rpc/kitex_demo/kitex_gen/api/apiservice"
)

func NewServer(address string) (server.Server, error) {
	addr, err := net.ResolveTCPAddr("tcp", address)
	if err != nil {
		return nil, err
	}
	svr := api.NewServer(new(APIServiceImpl), server.WithServiceAddr(addr))
	return svr, nil
}
