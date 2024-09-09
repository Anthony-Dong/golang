package curl

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"

	"github.com/anthony-dong/golang/command"
	"github.com/anthony-dong/golang/pkg/idl"
	"github.com/anthony-dong/golang/pkg/logs"
	"github.com/anthony-dong/golang/pkg/rpc"
	"github.com/anthony-dong/golang/pkg/utils"
)

func NewCurlCommand(configProvider func() *command.CurlConfig) (*cobra.Command, error) {
	reqUrl := ""
	reqBody := ""
	reqHeader := make([]string, 0)
	listMethods := false
	showExample := false
	idlInfo := rpc.IDLInfo{}
	timeout := time.Second * 180
	enableModifyReq := false
	cmd := &cobra.Command{
		Use:     "curl",
		Short:   `Send thrift like curl`,
		Example: `curl --url 'thrift://xxx.xxx.xxx/RPCMethod?addr=localhost:8888&env=prod&cluster=default' --header 'h1: v1' --header 'h2: v2' --data '{"k1": "v1"} -v'`,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			var (
				client     rpc.Client
				rpcRequest *rpc.Request
				err        error
			)
			rpcRequest, err = rpc.NewRpcRequest(reqUrl, reqHeader, reqBody)
			if err != nil {
				return err
			}
			rpcRequest.Timeout = utils.NewJsonDuration(timeout)
			rpcRequest.EnableModifyRequest = enableModifyReq
			if !showExample && !listMethods {
				logs.CtxInfo(ctx, "rpc request: %s", rpcRequest.String())
			}
			config := configProvider()
			if config != nil && config.NewClient != nil {
				if client, err = config.NewClient(ctx, rpcRequest, &idlInfo); err != nil {
					return err
				}
			} else {
				if idlInfo.Main == "" {
					return fmt.Errorf(`new local idl find err: not found main idl: %q`, idlInfo.Main)
				}
				client = rpc.NewThriftClient(idl.NewDescriptorProvider(idl.NewMemoryIDLProvider(idlInfo.Main)))
			}

			if listMethods {
				allMethods, err := client.ListMethods(ctx)
				if err != nil {
					return fmt.Errorf(`list methods find err: %v`, err)
				}
				logs.CtxInfo(ctx, "methods:\n%s", utils.ToJson(allMethods, true))
				return nil
			}

			if showExample {
				jsonExample, err := client.GetExampleCode(ctx, &rpc.Method{RPCMethod: rpcRequest.RPCMethod})
				if err != nil {
					return fmt.Errorf(`new request example find err: %v`, err)
				}
				logs.CtxInfo(ctx, "new request example\n%s", jsonExample.Body)
				return nil
			}

			rpcResponse, err := client.Do(ctx, rpcRequest)
			if err != nil {
				return fmt.Errorf(`do rpc request find err: %v`, err)
			}
			flag := "success"
			if rpcResponse.IsError {
				flag = "error"
			}
			logs.CtxInfo(ctx, "rpc response %s: %s", flag, rpcResponse.String())
			return nil
		},
	}
	cmd.Flags().StringVar(&reqUrl, "url", "", "The request url")
	cmd.Flags().StringVar(&idlInfo.Main, "idl", "", "The main IDL local path")
	cmd.Flags().StringVar(&idlInfo.Branch, "branch", "", "The Remote IDL branch/version/commit(if supports it)")
	cmd.Flags().StringSliceVarP(&reqHeader, "header", "H", []string{}, "The request header")
	cmd.Flags().StringVar(&reqBody, "data", "", "The request body")
	cmd.Flags().BoolVar(&listMethods, "methods", false, "List all the methods")
	cmd.Flags().BoolVar(&showExample, "example", false, "Generate request example data")
	cmd.Flags().DurationVar(&timeout, "timeout", timeout, "The request timeout")
	cmd.Flags().BoolVar(&enableModifyReq, "modify", false, "Enable the cli to modify the request")
	return cmd, nil
}
