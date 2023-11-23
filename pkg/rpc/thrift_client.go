package rpc

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/anthony-dong/golang/pkg/utils"

	"github.com/bytedance/gopkg/cloud/metainfo"

	kitex_client "github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/client/genericclient"
	"github.com/cloudwego/kitex/pkg/generic"
	"github.com/cloudwego/kitex/pkg/generic/thrift"
	"github.com/cloudwego/kitex/transport"

	"github.com/anthony-dong/golang/pkg/idl"
	"github.com/anthony-dong/golang/pkg/logs"
)

var _ Client = (*thriftClient)(nil)

type thriftClient struct {
}

func NewThriftClient() *thriftClient {
	return &thriftClient{}
}

func init() {
	thrift.SetDefaultParseMode(thrift.CombineServices) // with combine service
}

const (
	defaultConnTimeout = time.Second * 5
	defaultRcpTimeout  = time.Second * 60
)

type ThriftClientOption struct {
	Protocol transport.Protocol
}

func (t *thriftClient) Send(ctx context.Context, req *Request) (*Response, error) {
	if req.Service == "" || req.Method == "" {
		return nil, fmt.Errorf(`service or method connot be empty`)
	}
	logs.CtxDebug(ctx, "thrift client start request: %s", utils.ToJson(req))
	if ctx == nil {
		ctx = context.Background()
	}
	start := time.Now()
	idlProvider, err := req.NewIDLProvider(ctx)
	if err != nil {
		return nil, err
	}
	client, err := t.GetJsonThriftClient(ctx, req.Service, idlProvider)
	if err != nil {
		return nil, err
	}
	defer client.Close()

	instance := req.Instance

	// 获取用户请求地址的集群信息，主要原因是因为mesh这边会根据ip是否匹配集群，来确定流量规则，如果不匹配是不会携带mesh的token从而导致服务鉴定失败，前提是下游开启了服务鉴定！

	reqStart := time.Now()
	response, err := client.GenericCall(ctx, req.Method, utils.Bytes2String(req.Body))

	return &Response{
		TotalSpend: time.Now().Sub(start),
		Spend:      time.Now().Sub(reqStart),
		Body: func(r interface{}) json.RawMessage {
			if err != nil {
				return utils.String2Bytes(utils.ToJson(map[string]string{"error": err.Error()}))
			}
			str, _ := r.(string)
			return utils.String2Bytes(str)
		}(response),
		Extra: ResponseExtra{
			Instance: instance,
			MetaInfo: metainfo.GetAllValues(ctx),
		},
	}, nil
}

func (t *thriftClient) GetJsonThriftClient(ctx context.Context, psm string, provider idl.DescriptorProvider, ops ...func(*ThriftClientOption)) (v genericclient.Client, err error) {
	option := &ThriftClientOption{Protocol: transport.PurePayload}
	descriptorProvider, err := provider.DescriptorProvider()
	if err != nil {
		return nil, err
	}
	thriftGeneric, err := generic.JSONThriftGeneric(descriptorProvider)
	if err != nil {
		return nil, fmt.Errorf("new thrift json client find err: %v", err)
	}
	clientOps := make([]kitex_client.Option, 0)
	if option.Protocol != 0 {
		clientOps = append(clientOps, kitex_client.WithTransportProtocol(option.Protocol))
	}
	kClient, err := genericclient.NewClient(psm, thriftGeneric, clientOps...)
	if err != nil {
		return nil, fmt.Errorf("new thrift client find err: %v", err)
	}
	return kClient, nil
}

func (t *thriftClient) MethodList(ctx context.Context, req *Request) ([]*Method, error) {
	idlProvider, err := req.NewIDLProvider(ctx)
	if err != nil {
		return nil, err
	}
	provider, err := idlProvider.DescriptorProvider()
	if err != nil {
		return nil, err
	}
	defer provider.Close()
	idlAst := <-provider.Provide()
	result := make([]*Method, 0)
	for _, function := range idlAst.Functions {
		result = append(result, &Method{
			RPCMethod: function.Name,
		})
	}
	return result, nil
}

func (t *thriftClient) ExampleCode(ctx context.Context, request *Request) (string, error) {
	_provider, err := request.NewIDLProvider(ctx)
	if err != nil {
		return "", err
	}
	desc, err := _provider.DescriptorProvider()
	if err != nil {
		return "", err
	}
	provider := <-desc.Provide()
	function, err := provider.LookupFunctionByMethod(request.Method)
	if err != nil {
		return "", err
	}
	value := GetExampleValue(function.Request, nil, &Option{
		Generator: NewFixedGenerator(),
	})
	return utils.ToJson(value), nil
}
