package tcp

import (
	"context"
	"io"
	"net/http"
	"strings"

	"github.com/spf13/cobra"

	"github.com/anthony-dong/golang/pkg/logs"
	"github.com/anthony-dong/golang/pkg/utils"
)

func NewHTTPEchoServiceCommand() (*cobra.Command, error) {
	addr := ""
	command := &cobra.Command{
		Use:   "http_echo_service",
		Short: "create a http echo service",
		RunE: func(cmd *cobra.Command, args []string) error {
			logs.Info("listening on %s", addr)
			return newHTTPEchoService(cmd.Context(), addr)
		},
	}
	command.Flags().StringVar(&addr, "addr", ":8080", "监听地址")
	return command, nil
}

func newHTTPEchoService(ctx context.Context, addr string) error {
	return http.ListenAndServe(addr, http.HandlerFunc(func(writer http.ResponseWriter, req *http.Request) {
		for k, vs := range req.Header {
			if strings.HasPrefix(k, "X-") {
				for _, v := range vs {
					writer.Header().Add(k, v)
				}
			}
		}
		if req.Header.Get("content-encoding") != "" {
			writer.Header().Set("content-encoding", req.Header.Get("content-encoding"))
		}
		writer.Header().Set("content-type", req.Header.Get("content-type"))
		if err := copyRequest(writer, req); err != nil {
			logs.CtxError(ctx, "err: %v", err)
		}
		logs.Builder().Info().KV("method", req.Method).KV("host", req.URL.Host).KV("path", req.URL.Path).KV("content-size", req.ContentLength).Emit(ctx)
		return
	}))
}

func copyRequest(writer http.ResponseWriter, req *http.Request) error {
	if req.ContentLength < 1024*1024 {
		all, err := io.ReadAll(req.Body)
		if err != nil {
			return err
		}
		if _, err := writer.Write(all); err != nil {
			return err
		}
		return nil
	}
	if _, err := writer.Write(utils.ToJsonByte(map[string]interface{}{
		"ContentLength": req.ContentLength,
	})); err != nil {
		return err
	}
	return nil
}
