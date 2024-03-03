package rpc

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/cloudwego/kitex/pkg/generic/descriptor"
	"github.com/iancoleman/orderedmap"

	"github.com/bytedance/gopkg/cloud/metainfo"
	kitex_client "github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/client/callopt"
	"github.com/cloudwego/kitex/client/genericclient"
	"github.com/cloudwego/kitex/pkg/generic"
	"github.com/cloudwego/kitex/pkg/generic/thrift"

	"github.com/anthony-dong/golang/pkg/idl"
	"github.com/anthony-dong/golang/pkg/utils"
)

var _ Client = (*ThriftClient)(nil)

type ThriftClient struct {
	IDLProvider       idl.DescriptorProvider
	Option            ThriftClientOption
	NewClient         func(ctx context.Context, serviceName string, provider idl.DescriptorProvider, option ThriftClientOption) (v genericclient.Client, err error)
	PreHandlerRequest func(ctx context.Context, req *Request) (context.Context, []callopt.Option, error)
}

func NewThriftClient(provider idl.DescriptorProvider, ops ...func(client *ThriftClientOption)) *ThriftClient {
	if provider == nil {
		panic(`new thrift client find err: idl provider is nil`)
	}
	option := ThriftClientOption{Options: []kitex_client.Option{}}
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
	DefaultConnTimeout = time.Second * 5
	DefaultRcpTimeout  = time.Second * 1800
)

type ThriftClientOption struct {
	Options []kitex_client.Option
}

func (t *ThriftClient) preHandlerRequest(ctx context.Context, req *Request) (context.Context, []callopt.Option, error) {
	options := make([]callopt.Option, 0)
	options = append(options, callopt.WithConnectTimeout(DefaultConnTimeout))
	if req.Timeout > 0 {
		options = append(options, callopt.WithRPCTimeout(req.Timeout.Duration()))
	} else {
		options = append(options, callopt.WithRPCTimeout(DefaultRcpTimeout))
	}
	if req.Addr != "" {
		options = append(options, callopt.WithHostPort(req.Addr))
	}
	for _, kv := range req.Tag {
		options = append(options, callopt.WithTag(kv.Key, kv.Value))
	}
	header := make(map[string]string, len(req.Header))
	for _, kv := range req.Header {
		if strings.HasPrefix(kv.Key, metainfo.PrefixPersistent) {
			header[kv.Key] = kv.Value
			continue
		}
		if strings.HasPrefix(kv.Key, metainfo.PrefixTransient) {
			header[kv.Key] = kv.Value
			continue
		}
		header[metainfo.PrefixPersistent+kv.Key] = kv.Value
	}
	ctx = metainfo.SetMetaInfoFromMap(ctx, header)
	return ctx, options, nil
}

func (t *ThriftClient) Do(ctx context.Context, req *Request) (*Response, error) {
	if req.RPCMethod == "" {
		return nil, fmt.Errorf(`req method connot be empty`)
	}
	if req.Service == "" {
		req.Service = "-"
	}
	if ctx == nil {
		ctx = context.Background()
	}

	newClientFunc := t.newClient
	if t.NewClient != nil {
		newClientFunc = t.NewClient
	}
	client, err := newClientFunc(ctx, req.Service, t.IDLProvider, t.Option)
	if err != nil {
		return nil, err
	}
	defer client.Close()

	preHandlerRequest := t.preHandlerRequest
	if t.PreHandlerRequest != nil {
		preHandlerRequest = t.PreHandlerRequest
	}
	rpcCtx, callOps, err := preHandlerRequest(ctx, req)
	if err != nil {
		return nil, err
	}
	reqStart := time.Now()
	response, err := client.GenericCall(rpcCtx, req.RPCMethod, utils.Bytes2String(req.Body), callOps...)
	return &Response{
		Spend: time.Now().Sub(reqStart),
		Body: func(r interface{}) json.RawMessage {
			if err != nil {
				return utils.ToJsonByte(map[string]interface{}{
					"error": err.Error(),
				})
			}
			str, _ := r.(string)
			return utils.String2Bytes(str)
		}(response),
		IsError: err != nil,
		Header:  []*KV{}, // thrift rpc 不支持设置response header
	}, nil
}

func (t *ThriftClient) newClient(ctx context.Context, serviceName string, provider idl.DescriptorProvider, option ThriftClientOption) (v genericclient.Client, err error) {
	descriptorProvider, err := provider.DescriptorProvider()
	if err != nil {
		return nil, err
	}
	thriftGeneric, err := generic.JSONThriftGeneric(descriptorProvider)
	if err != nil {
		return nil, fmt.Errorf("new thrift json client find err: %v", err)
	}
	kClient, err := genericclient.NewClient(serviceName, thriftGeneric, option.Options...)
	if err != nil {
		return nil, fmt.Errorf("new thrift client find err: %v", err)
	}
	return kClient, nil
}

func (t *ThriftClient) GetIDLDescriptor() (*descriptor.ServiceDescriptor, error) {
	provider, err := t.IDLProvider.DescriptorProvider()
	if err != nil {
		return nil, err
	}
	return <-provider.Provide(), nil
}

func (t *ThriftClient) ListMethods(ctx context.Context) ([]*Method, error) {
	provider, err := t.GetIDLDescriptor()
	if err != nil {
		return nil, err
	}
	result := make([]*Method, 0, len(provider.Functions))
	for _, elem := range provider.Functions {
		result = append(result, &Method{
			RPCMethod: elem.Name,
		})
	}
	return result, nil
}

func (t *ThriftClient) GetExampleCode(ctx context.Context, method *Method) (*ExampleCode, error) {
	provider, err := t.GetIDLDescriptor()
	if err != nil {
		return nil, err
	}
	function := provider.Functions[method.RPCMethod]
	if function == nil {
		return nil, fmt.Errorf(`not found rpc method: %s`, method.RPCMethod)
	}
	value, err := GetThriftExampleValue(function.Request, nil, NewThriftExampleOption())
	if err != nil {
		return nil, err
	}
	orderedMap, isOk := value.(*orderedmap.OrderedMap)
	if isOk {
		req, isExist := orderedMap.Get("req")
		if isExist {
			value = req
		}
	}
	return &ExampleCode{Body: utils.ToJsonByte(value, true)}, nil
}
