package main

import (
	"github.com/anthony-dong/golang/pkg/rpc/kitex_demo/handler"
)

func main() {
	server, err := handler.NewServer(":8888")
	if err != nil {
		panic(err)
	}
	if err := server.Run(); err != nil {
		panic(err)
	}
}
