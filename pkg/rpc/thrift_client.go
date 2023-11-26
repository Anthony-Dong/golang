package rpc

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/iancoleman/orderedmap"

	"github.com/cloudwego/kitex/client/callopt"

	"github.com/anthony-dong/golang/pkg/utils"

	"github.com/bytedance/gopkg/cloud/metainfo"

	kitex_client "github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/client/genericclient"
	"github.com/cloudwego/kitex/pkg/generic"
	"github.com/cloudwego/kitex/pkg/generic/thrift"
	"github.com/cloudwego/kitex/transport"

	"github.com/anthony-dong/golang/pkg/idl"
)

var _ Client = (*ThriftClient)(nil)

type ThriftClient struct {
	IDLProvider idl.DescriptorProvider
	Option      ThriftClientOption
}

func NewThriftClient(provider idl.DescriptorProvider, ops ...func(client *ThriftClientOption)) *ThriftClient {
	if provider == nil {
		panic(`new thrift client find err: idl provider is nil`)
	}
	option := ThriftClientOption{Protocol: transport.PurePayload}
	for _, op := range ops {
		op(&option)
	}
	return &ThriftClient{
		IDLProvider: provider,
		Option:      option,
	}
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

func (t *ThriftClient) Send(ctx context.Context, req *Request) (*Response, error) {
	if req.Method == "" {
		return nil, fmt.Errorf(`req method connot be empty`)
	}
	if req.Service == "" {
		req.Service = "-"
	}
	if ctx == nil {
		ctx = context.Background()
	}
	start := time.Now()
	client, err := t.NewThriftClient(ctx, req.Service, t.IDLProvider)
	if err != nil {
		return nil, err
	}
	defer client.Close()

	endpoint := req.Endpoint
	options := make([]callopt.Option, 0)
	options = append(options, callopt.WithConnectTimeout(defaultConnTimeout))
	options = append(options, callopt.WithRPCTimeout(defaultRcpTimeout))
	if endpoint.Addr != "" {
		options = append(options, callopt.WithHostPort(endpoint.Addr))
	}
	for k, v := range endpoint.Tag {
		options = append(options, callopt.WithTag(k, v))
	}
	reqStart := time.Now()
	response, err := client.GenericCall(ctx, req.Method, utils.Bytes2String(req.Body), options...)
	if err != nil {
		return nil, err
	}
	return &Response{
		TotalSpend: time.Now().Sub(start),
		Spend:      time.Now().Sub(reqStart),
		Body: func(r interface{}) json.RawMessage {
			str, _ := r.(string)
			return utils.String2Bytes(str)
		}(response),
		Extra: ResponseExtra{
			Endpoint: endpoint,
			MetaInfo: metainfo.GetAllValues(ctx),
		},
	}, nil
}

func (t *ThriftClient) NewThriftClient(ctx context.Context, psm string, provider idl.DescriptorProvider) (v genericclient.Client, err error) {
	descriptorProvider, err := provider.DescriptorProvider()
	if err != nil {
		return nil, err
	}
	thriftGeneric, err := generic.JSONThriftGeneric(descriptorProvider)
	if err != nil {
		return nil, fmt.Errorf("new thrift json client find err: %v", err)
	}
	clientOps := make([]kitex_client.Option, 0, 1)
	clientOps = append(clientOps, kitex_client.WithTransportProtocol(t.Option.Protocol))
	kClient, err := genericclient.NewClient(psm, thriftGeneric, clientOps...)
	if err != nil {
		return nil, fmt.Errorf("new thrift client find err: %v", err)
	}
	return kClient, nil
}

func (t *ThriftClient) MethodList(ctx context.Context, req *Request) ([]*Method, error) {
	provider, err := t.IDLProvider.DescriptorProvider()
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

func (t *ThriftClient) ExampleCode(ctx context.Context, request *Request) (string, error) {
	desc, err := t.IDLProvider.DescriptorProvider()
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
	// 如果请求参数仅有一个
	orderedMap, isOk := value.(*orderedmap.OrderedMap)
	if isOk {
		req, isExist := orderedMap.Get("req")
		if isExist {
			return utils.ToJson(req), nil
		}
	}
	return utils.ToJson(value), nil
}
