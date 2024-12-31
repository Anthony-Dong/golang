package proxy

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/spf13/cobra"

	"github.com/anthony-dong/golang/command"
	"github.com/anthony-dong/golang/pkg/logs"
	"github.com/anthony-dong/golang/pkg/rpc"
	"github.com/anthony-dong/golang/pkg/utils"
)

func NewThriftProxyCommand(cfg func() *command.CurlConfig) (*cobra.Command, error) {
	cmd := &cobra.Command{
		Use:   "thrift_proxy",
		Short: "Thrift proxy tool",
		RunE: func(cmd *cobra.Command, args []string) error {
			ips, err := utils.GetAllIP(true)
			if err != nil {
				return err
			}
			for _, ip := range ips {
				logs.CtxInfo(cmd.Context(), "listener: http://%s:8080", ip.String())
			}
			return http.ListenAndServe(":8080", http.HandlerFunc(handlerThriftRequest(cfg)))
		},
	}
	return cmd, nil
}

func handlerThriftRequest(configProvider func() *command.CurlConfig) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := context.Background()
		respError := func(err error) {
			w.WriteHeader(http.StatusBadRequest)
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write(utils.ToJsonByte(map[string]interface{}{
				"error": err.Error(),
			}))
		}
		respSuccess := func(body []byte) {
			w.WriteHeader(http.StatusOK)
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write(body)
		}
		query := r.URL.Query()
		serviceName := query.Get("service_name")
		rpcMethod := query.Get("rpc_method")
		addr := query.Get("addr")
		mainIdl := query.Get("idl")
		enableModifyRequest := query.Get("enable_modify_request") == "1"
		reqBody, err := io.ReadAll(r.Body)
		if err != nil {
			respError(err)
			return
		}
		req := &rpc.Request{
			Protocol:            rpc.ProtocolThrift,
			ServiceName:         serviceName,
			RPCMethod:           rpcMethod,
			Body:                reqBody,
			Addr:                addr,
			Tag:                 rpc.GetTagFromQuery(query, []string{"service_name", "rpc_method", "addr", "enable_modify_request", "idl"}),
			EnableModifyRequest: enableModifyRequest,
		}

		var client *rpc.ThriftClient
		config := configProvider()
		if config != nil && config.NewThriftClient != nil {
			if client, err = config.NewThriftClient(ctx, req); err != nil {
				respError(err)
				return
			}
		} else {
			if mainIdl == "" {
				respError(fmt.Errorf(`new local idl find err: not found main idl: %q`, mainIdl))
			}
			if client, err = rpc.NewThriftClient(rpc.NewLocalIDLProvider(map[string]string{req.ServiceName: mainIdl})); err != nil {
				respError(err)
				return
			}
		}
		response, err := client.Do(ctx, req)
		if err != nil {
			respError(err)
			return
		}
		respSuccess(response.Body)
	}
}
