package curl

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"

	"github.com/anthony-dong/golang/command"
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
	idlConfig := rpc.IDLConfig{}
	timeout := time.Second * 180
	enableModifyReq := false
	cmd := &cobra.Command{
		Use:     "curl",
		Short:   `Send thrift like curl`,
		Example: `curl --url 'thrift://xxx.xxx.xxx/RPCMethod?addr=localhost:8888&env=prod&cluster=default' --header 'h1: v1' --header 'h2: v2' --data '{"k1": "v1"} -v'`,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			var (
				client *rpc.ThriftClient
				req    *rpc.Request
				err    error
			)
			req, err = rpc.NewRpcRequest(reqUrl, reqHeader, reqBody)
			if err != nil {
				return err
			}
			req.Timeout = utils.NewJsonDuration(timeout)
			req.EnableModifyRequest = enableModifyReq
			req.IDLConfig = &idlConfig
			if !showExample && !listMethods {
				logs.CtxInfo(ctx, "rpc request info: %s", utils.ToString(req.BasicInfo()))
			}
			config := configProvider()
			if config != nil && config.NewThriftClient != nil {
				if client, err = config.NewThriftClient(ctx, req); err != nil {
					return err
				}
			} else {
				if idlConfig.Main == "" {
					return fmt.Errorf(`new local idl find err: not found main idl: %q`, idlConfig.Main)
				}
				if client, err = rpc.NewThriftClient(rpc.NewLocalIDLProvider(map[string]string{req.ServiceName: idlConfig.Main})); err != nil {
					return err
				}
			}
			if listMethods {
				allMethods, err := client.ListMethods(ctx, req.ServiceName, req.IDLConfig)
				if err != nil {
					return fmt.Errorf(`list methods find err: %v`, err)
				}
				logs.CtxInfo(ctx, "methods:\n%s", utils.ToJson(allMethods, true))
				return nil
			}
			if showExample {
				jsonExample, err := client.GetExampleCode(ctx, req.ServiceName, req.IDLConfig, req.ServiceName)
				if err != nil {
					return fmt.Errorf(`new request example find err: %v`, err)
				}
				logs.CtxInfo(ctx, "new request example\n%s", jsonExample)
				return nil
			}

			resp, err := client.Do(ctx, req)
			if err != nil {
				return fmt.Errorf(`do rpc request find err: %v`, err)
			}
			flag := "success"
			if resp.IsError {
				flag = "error"
			}
			logs.CtxInfo(ctx, "rpc response %s:\n%s", flag, utils.PrettyJsonBytes(resp.Body))
			return nil
		},
	}
	cmd.Flags().StringVar(&reqUrl, "url", "", "The request url")
	cmd.Flags().StringVar(&idlConfig.Main, "idl", "", "The main IDL local path")
	cmd.Flags().StringVar(&idlConfig.Branch, "branch", "", "The Remote IDL branch/version/commit(if supports it)")
	cmd.Flags().StringSliceVarP(&reqHeader, "header", "H", []string{}, "The request header")
	cmd.Flags().StringVar(&reqBody, "data", "", "The request body")
	cmd.Flags().BoolVar(&listMethods, "methods", false, "List all the methods")
	cmd.Flags().BoolVar(&showExample, "example", false, "Generate request example data")
	cmd.Flags().DurationVar(&timeout, "timeout", timeout, "The request timeout")
	cmd.Flags().BoolVar(&enableModifyReq, "modify", false, "Enable the cli to modify the request")
	return cmd, nil
}
