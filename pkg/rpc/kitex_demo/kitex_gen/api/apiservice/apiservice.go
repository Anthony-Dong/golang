// Code generated by Kitex v0.11.3. DO NOT EDIT.

package apiservice

import (
	"context"
	"errors"

	client "github.com/cloudwego/kitex/client"
	kitex "github.com/cloudwego/kitex/pkg/serviceinfo"

	api "github.com/anthony-dong/golang/pkg/rpc/kitex_demo/kitex_gen/api"
)

var errInvalidMessageType = errors.New("invalid message type for service method handler")

var serviceMethods = map[string]kitex.MethodInfo{
	"TestStruct": kitex.NewMethodInfo(
		testStructHandler,
		newAPIServiceTestStructArgs,
		newAPIServiceTestStructResult,
		false,
		kitex.WithStreamingMode(kitex.StreamingNone),
	),
	"TestVoid": kitex.NewMethodInfo(
		testVoidHandler,
		newAPIServiceTestVoidArgs,
		newAPIServiceTestVoidResult,
		false,
		kitex.WithStreamingMode(kitex.StreamingNone),
	),
	"TestOnewayVoid": kitex.NewMethodInfo(
		testOnewayVoidHandler,
		newAPIServiceTestOnewayVoidArgs,
		nil,
		true,
		kitex.WithStreamingMode(kitex.StreamingNone),
	),
	"TestList": kitex.NewMethodInfo(
		testListHandler,
		newAPIServiceTestListArgs,
		newAPIServiceTestListResult,
		false,
		kitex.WithStreamingMode(kitex.StreamingNone),
	),
	"TestSet": kitex.NewMethodInfo(
		testSetHandler,
		newAPIServiceTestSetArgs,
		newAPIServiceTestSetResult,
		false,
		kitex.WithStreamingMode(kitex.StreamingNone),
	),
	"TestMap": kitex.NewMethodInfo(
		testMapHandler,
		newAPIServiceTestMapArgs,
		newAPIServiceTestMapResult,
		false,
		kitex.WithStreamingMode(kitex.StreamingNone),
	),
	"TestIntMap": kitex.NewMethodInfo(
		testIntMapHandler,
		newAPIServiceTestIntMapArgs,
		newAPIServiceTestIntMapResult,
		false,
		kitex.WithStreamingMode(kitex.StreamingNone),
	),
	"TestString": kitex.NewMethodInfo(
		testStringHandler,
		newAPIServiceTestStringArgs,
		newAPIServiceTestStringResult,
		false,
		kitex.WithStreamingMode(kitex.StreamingNone),
	),
}

var (
	aPIServiceServiceInfo                = NewServiceInfo()
	aPIServiceServiceInfoForClient       = NewServiceInfoForClient()
	aPIServiceServiceInfoForStreamClient = NewServiceInfoForStreamClient()
)

// for server
func serviceInfo() *kitex.ServiceInfo {
	return aPIServiceServiceInfo
}

// for stream client
func serviceInfoForStreamClient() *kitex.ServiceInfo {
	return aPIServiceServiceInfoForStreamClient
}

// for client
func serviceInfoForClient() *kitex.ServiceInfo {
	return aPIServiceServiceInfoForClient
}

// NewServiceInfo creates a new ServiceInfo containing all methods
func NewServiceInfo() *kitex.ServiceInfo {
	return newServiceInfo(false, true, true)
}

// NewServiceInfo creates a new ServiceInfo containing non-streaming methods
func NewServiceInfoForClient() *kitex.ServiceInfo {
	return newServiceInfo(false, false, true)
}
func NewServiceInfoForStreamClient() *kitex.ServiceInfo {
	return newServiceInfo(true, true, false)
}

func newServiceInfo(hasStreaming bool, keepStreamingMethods bool, keepNonStreamingMethods bool) *kitex.ServiceInfo {
	serviceName := "APIService"
	handlerType := (*api.APIService)(nil)
	methods := map[string]kitex.MethodInfo{}
	for name, m := range serviceMethods {
		if m.IsStreaming() && !keepStreamingMethods {
			continue
		}
		if !m.IsStreaming() && !keepNonStreamingMethods {
			continue
		}
		methods[name] = m
	}
	extra := map[string]interface{}{
		"PackageName": "api",
	}
	if hasStreaming {
		extra["streaming"] = hasStreaming
	}
	svcInfo := &kitex.ServiceInfo{
		ServiceName:     serviceName,
		HandlerType:     handlerType,
		Methods:         methods,
		PayloadCodec:    kitex.Thrift,
		KiteXGenVersion: "v0.11.3",
		Extra:           extra,
	}
	return svcInfo
}

func testStructHandler(ctx context.Context, handler interface{}, arg, result interface{}) error {
	realArg := arg.(*api.APIServiceTestStructArgs)
	realResult := result.(*api.APIServiceTestStructResult)
	success, err := handler.(api.APIService).TestStruct(ctx, realArg.Req)
	if err != nil {
		return err
	}
	realResult.Success = success
	return nil
}
func newAPIServiceTestStructArgs() interface{} {
	return api.NewAPIServiceTestStructArgs()
}

func newAPIServiceTestStructResult() interface{} {
	return api.NewAPIServiceTestStructResult()
}

func testVoidHandler(ctx context.Context, handler interface{}, arg, result interface{}) error {
	realArg := arg.(*api.APIServiceTestVoidArgs)

	err := handler.(api.APIService).TestVoid(ctx, realArg.Req)
	if err != nil {
		return err
	}

	return nil
}
func newAPIServiceTestVoidArgs() interface{} {
	return api.NewAPIServiceTestVoidArgs()
}

func newAPIServiceTestVoidResult() interface{} {
	return api.NewAPIServiceTestVoidResult()
}

func testOnewayVoidHandler(ctx context.Context, handler interface{}, arg, result interface{}) error {
	realArg := arg.(*api.APIServiceTestOnewayVoidArgs)

	err := handler.(api.APIService).TestOnewayVoid(ctx, realArg.Req)
	if err != nil {
		return err
	}

	return nil
}
func newAPIServiceTestOnewayVoidArgs() interface{} {
	return api.NewAPIServiceTestOnewayVoidArgs()
}

func testListHandler(ctx context.Context, handler interface{}, arg, result interface{}) error {
	realArg := arg.(*api.APIServiceTestListArgs)
	realResult := result.(*api.APIServiceTestListResult)
	success, err := handler.(api.APIService).TestList(ctx, realArg.Req)
	if err != nil {
		return err
	}
	realResult.Success = success
	return nil
}
func newAPIServiceTestListArgs() interface{} {
	return api.NewAPIServiceTestListArgs()
}

func newAPIServiceTestListResult() interface{} {
	return api.NewAPIServiceTestListResult()
}

func testSetHandler(ctx context.Context, handler interface{}, arg, result interface{}) error {
	realArg := arg.(*api.APIServiceTestSetArgs)
	realResult := result.(*api.APIServiceTestSetResult)
	success, err := handler.(api.APIService).TestSet(ctx, realArg.Req)
	if err != nil {
		return err
	}
	realResult.Success = success
	return nil
}
func newAPIServiceTestSetArgs() interface{} {
	return api.NewAPIServiceTestSetArgs()
}

func newAPIServiceTestSetResult() interface{} {
	return api.NewAPIServiceTestSetResult()
}

func testMapHandler(ctx context.Context, handler interface{}, arg, result interface{}) error {
	realArg := arg.(*api.APIServiceTestMapArgs)
	realResult := result.(*api.APIServiceTestMapResult)
	success, err := handler.(api.APIService).TestMap(ctx, realArg.Req)
	if err != nil {
		return err
	}
	realResult.Success = success
	return nil
}
func newAPIServiceTestMapArgs() interface{} {
	return api.NewAPIServiceTestMapArgs()
}

func newAPIServiceTestMapResult() interface{} {
	return api.NewAPIServiceTestMapResult()
}

func testIntMapHandler(ctx context.Context, handler interface{}, arg, result interface{}) error {
	realArg := arg.(*api.APIServiceTestIntMapArgs)
	realResult := result.(*api.APIServiceTestIntMapResult)
	success, err := handler.(api.APIService).TestIntMap(ctx, realArg.Req)
	if err != nil {
		return err
	}
	realResult.Success = success
	return nil
}
func newAPIServiceTestIntMapArgs() interface{} {
	return api.NewAPIServiceTestIntMapArgs()
}

func newAPIServiceTestIntMapResult() interface{} {
	return api.NewAPIServiceTestIntMapResult()
}

func testStringHandler(ctx context.Context, handler interface{}, arg, result interface{}) error {
	realArg := arg.(*api.APIServiceTestStringArgs)
	realResult := result.(*api.APIServiceTestStringResult)
	success, err := handler.(api.APIService).TestString(ctx, realArg.Req)
	if err != nil {
		return err
	}
	realResult.Success = &success
	return nil
}
func newAPIServiceTestStringArgs() interface{} {
	return api.NewAPIServiceTestStringArgs()
}

func newAPIServiceTestStringResult() interface{} {
	return api.NewAPIServiceTestStringResult()
}

type kClient struct {
	c client.Client
}

func newServiceClient(c client.Client) *kClient {
	return &kClient{
		c: c,
	}
}

func (p *kClient) TestStruct(ctx context.Context, req *api.Request) (r *api.Response, err error) {
	var _args api.APIServiceTestStructArgs
	_args.Req = req
	var _result api.APIServiceTestStructResult
	if err = p.c.Call(ctx, "TestStruct", &_args, &_result); err != nil {
		return
	}
	return _result.GetSuccess(), nil
}

func (p *kClient) TestVoid(ctx context.Context, req *api.Request) (err error) {
	var _args api.APIServiceTestVoidArgs
	_args.Req = req
	var _result api.APIServiceTestVoidResult
	if err = p.c.Call(ctx, "TestVoid", &_args, &_result); err != nil {
		return
	}
	return nil
}

func (p *kClient) TestOnewayVoid(ctx context.Context, req *api.Request) (err error) {
	var _args api.APIServiceTestOnewayVoidArgs
	_args.Req = req
	if err = p.c.Call(ctx, "TestOnewayVoid", &_args, nil); err != nil {
		return
	}
	return nil
}

func (p *kClient) TestList(ctx context.Context, req *api.Request) (r []*api.Response, err error) {
	var _args api.APIServiceTestListArgs
	_args.Req = req
	var _result api.APIServiceTestListResult
	if err = p.c.Call(ctx, "TestList", &_args, &_result); err != nil {
		return
	}
	return _result.GetSuccess(), nil
}

func (p *kClient) TestSet(ctx context.Context, req *api.Request) (r []*api.Response, err error) {
	var _args api.APIServiceTestSetArgs
	_args.Req = req
	var _result api.APIServiceTestSetResult
	if err = p.c.Call(ctx, "TestSet", &_args, &_result); err != nil {
		return
	}
	return _result.GetSuccess(), nil
}

func (p *kClient) TestMap(ctx context.Context, req *api.Request) (r map[string]*api.Response, err error) {
	var _args api.APIServiceTestMapArgs
	_args.Req = req
	var _result api.APIServiceTestMapResult
	if err = p.c.Call(ctx, "TestMap", &_args, &_result); err != nil {
		return
	}
	return _result.GetSuccess(), nil
}

func (p *kClient) TestIntMap(ctx context.Context, req *api.Request) (r map[int64]*api.Response, err error) {
	var _args api.APIServiceTestIntMapArgs
	_args.Req = req
	var _result api.APIServiceTestIntMapResult
	if err = p.c.Call(ctx, "TestIntMap", &_args, &_result); err != nil {
		return
	}
	return _result.GetSuccess(), nil
}

func (p *kClient) TestString(ctx context.Context, req *api.Request) (r string, err error) {
	var _args api.APIServiceTestStringArgs
	_args.Req = req
	var _result api.APIServiceTestStringResult
	if err = p.c.Call(ctx, "TestString", &_args, &_result); err != nil {
		return
	}
	return _result.GetSuccess(), nil
}
