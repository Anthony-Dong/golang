// Code generated by Kitex v0.11.3. DO NOT EDIT.
package apiservice

import (
	server "github.com/cloudwego/kitex/server"

	api "github.com/anthony-dong/golang/pkg/rpc/kitex_demo/kitex_gen/api"
)

// NewServer creates a server.Server with the given handler and options.
func NewServer(handler api.APIService, opts ...server.Option) server.Server {
	var options []server.Option

	options = append(options, opts...)
	options = append(options, server.WithCompatibleMiddlewareForUnary())

	svr := server.NewServer(options...)
	if err := svr.RegisterService(serviceInfo(), handler); err != nil {
		panic(err)
	}
	return svr
}

func RegisterService(svr server.Server, handler api.APIService, opts ...server.RegisterOption) error {
	return svr.RegisterService(serviceInfo(), handler, opts...)
}
