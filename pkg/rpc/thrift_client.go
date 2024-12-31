package rpc

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/cloudwego/kitex/transport"

	"github.com/anthony-dong/golang/pkg/logs"

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

// ThriftClient todo keepalive
type ThriftClient struct {
	idlProvider IDLProvider
	idlCache    map[string]*descriptor.ServiceDescriptor
	Option      ThriftClientOption
}

func NewThriftClient(provider IDLProvider, ops ...func(client *ThriftClientOption)) (*ThriftClient, error) {
	if provider == nil {
		return nil, fmt.Errorf(`idl provider is nil`)
	}
	option := ThriftClientOption{}
	for _, op := range ops {
		op(&option)
	}
	return &ThriftClient{
		idlProvider: provider,
		idlCache:    map[string]*descriptor.ServiceDescriptor{},
		Option:      option,
	}, nil
}

func init() {
	thrift.SetDefaultParseMode(thrift.CombineServices) // with combine service
}

const (
	DefaultConnTimeout = time.Second * 5
	DefaultRcpTimeout  = time.Second * 1800
)

type ThriftClientOption struct {
	NewClient         func(ctx context.Context, serviceName string, provider generic.DescriptorProvider, ops []kitex_client.Option) (v genericclient.Client, err error)
	PreHandlerRequest func(ctx context.Context, req *Request) (context.Context, []callopt.Option, error)
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

func NewTransientHeader(key string, value string) *KV {
	return &KV{Key: metainfo.PrefixTransient + key, Value: value}
}

func NewPersistentHeader(key string, value string) *KV {
	return &KV{Key: metainfo.PrefixPersistent + key, Value: value}
}

func (t *ThriftClient) GetIDLDescriptor(ctx context.Context, serviceName string, idlConfig *IDLConfig) (*descriptor.ServiceDescriptor, error) {
	if idlConfig == nil {
		idlConfig = &IDLConfig{}
	}
	key := serviceName + ":" + utils.ToJson(idlConfig)
	if desc := t.idlCache[key]; desc != nil {
		return desc, nil
	}
	memoryIDL, err := t.idlProvider.MemoryIDL(ctx, serviceName, idlConfig)
	if err != nil {
		return nil, err
	}
	desc, err := idl.ParseThriftIDL(memoryIDL)
	if err != nil {
		return nil, err
	}
	t.idlCache[key] = desc
	return desc, nil
}

func (t *ThriftClient) newClient(ctx context.Context, req *Request, desc *descriptor.ServiceDescriptor, ops []kitex_client.Option) (genericclient.Client, error) {
	if GetValue(req.Tag, "protocol") == "framed" {
		ops = append(ops, kitex_client.WithTransportProtocol(transport.TTHeaderFramed))
		logs.CtxInfo(ctx, "use custom protocol: %s", "TTHeaderFramed")
	}
	provider := idl.NewDescriptorProvider(desc)
	if t.Option.NewClient != nil {
		return t.Option.NewClient(ctx, req.ServiceName, provider, ops)
	}
	return t.newDefaultClient(ctx, req.ServiceName, provider, ops)
}

func (t *ThriftClient) Do(ctx context.Context, req *Request, ops ...kitex_client.Option) (*Response, error) {
	if req.RPCMethod == "" {
		return nil, fmt.Errorf(`req method connot be empty`)
	}
	if req.ServiceName == "" {
		req.ServiceName = "-"
	}
	if ctx == nil {
		ctx = context.Background()
	}
	idlDesc, err := t.GetIDLDescriptor(ctx, req.ServiceName, req.IDLConfig)
	if err != nil {
		return nil, err
	}
	client, err := t.newClient(ctx, req, idlDesc, ops)
	if err != nil {
		return nil, err
	}
	defer client.Close()

	preHandlerRequest := t.preHandlerRequest
	if t.Option.PreHandlerRequest != nil {
		preHandlerRequest = t.Option.PreHandlerRequest
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
			respData, _ := r.(string)
			if formatResp, err := FormatResponse(idlDesc, req.RPCMethod, respData); err == nil {
				respData = formatResp
			}
			return utils.String2Bytes(respData)
		}(response),
		IsError: err != nil,
		Header:  []*KV{}, // thrift rpc 不支持设置response header
	}, nil
}

func (t *ThriftClient) newDefaultClient(ctx context.Context, serviceName string, provider generic.DescriptorProvider, ops []kitex_client.Option) (v genericclient.Client, err error) {
	thriftGeneric, err := generic.JSONThriftGeneric(provider)
	if err != nil {
		return nil, fmt.Errorf("new thrift json client find err: %v", err)
	}
	kClient, err := genericclient.NewClient(serviceName, thriftGeneric, ops...)
	if err != nil {
		return nil, fmt.Errorf("new thrift client find err: %v", err)
	}
	return kClient, nil
}

func (t *ThriftClient) ListMethods(ctx context.Context, serviceName string, idlConfig *IDLConfig) ([]string, error) {
	idlDescriptor, err := t.GetIDLDescriptor(ctx, serviceName, idlConfig)
	if err != nil {
		return nil, err
	}
	return ListMethods(ctx, idlDescriptor)
}

func (t *ThriftClient) GetExampleCode(ctx context.Context, serviceName string, idlConfig *IDLConfig, method string) ([]byte, error) {
	idlDescriptor, err := t.GetIDLDescriptor(ctx, serviceName, idlConfig)
	if err != nil {
		return nil, err
	}
	function := idlDescriptor.Functions[method]
	if function == nil {
		return nil, fmt.Errorf(`not found rpc method: %s`, method)
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
	return utils.ToJsonByte(value, true), nil
}
