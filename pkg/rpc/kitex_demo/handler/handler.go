package handler

import (
	"context"

	"github.com/bytedance/gopkg/cloud/metainfo"

	"github.com/anthony-dong/golang/pkg/logs"
	api "github.com/anthony-dong/golang/pkg/rpc/kitex_demo/kitex_gen/api"
	"github.com/anthony-dong/golang/pkg/rpc/kitex_demo/kitex_gen/base"
	"github.com/anthony-dong/golang/pkg/utils"
)

var _ api.APIService = (*APIServiceImpl)(nil)

// APIServiceImpl implements the last service interface defined in the IDL.
type APIServiceImpl struct{}

// TestStruct implements the APIServiceImpl interface.
func (s *APIServiceImpl) TestStruct(ctx context.Context, req *api.Request) (resp *api.Response, err error) {
	logs.CtxInfo(ctx, "req: %s", utils.ToJson(req))
	return NewMockResponse(ctx, req), nil
}

// TestVoid implements the APIServiceImpl interface.
func (s *APIServiceImpl) TestVoid(ctx context.Context, req *api.Request) (err error) {
	logs.CtxInfo(ctx, "req: %s", utils.ToJson(req))
	return
}

// TestOnewayVoid implements the APIServiceImpl interface.
// Oneway methods are not guaranteed to receive 100% of the requests sent by the client.
// And the client may not perceive the loss of requests due to network packet loss.
// If possible, do not use oneway methods.
func (s *APIServiceImpl) TestOnewayVoid(ctx context.Context, req *api.Request) (err error) {
	logs.CtxInfo(ctx, "req: %s", utils.ToJson(req))
	return
}

// TestList implements the APIServiceImpl interface.
func (s *APIServiceImpl) TestList(ctx context.Context, req *api.Request) (resp []*api.Response, err error) {
	logs.CtxInfo(ctx, "req: %s", utils.ToJson(req))
	return []*api.Response{NewMockResponse(ctx, req)}, nil
}

// TestSet implements the APIServiceImpl interface.
func (s *APIServiceImpl) TestSet(ctx context.Context, req *api.Request) (resp []*api.Response, err error) {
	logs.CtxInfo(ctx, "req: %s", utils.ToJson(req))
	return []*api.Response{NewMockResponse(ctx, req)}, nil
}

// TestMap implements the APIServiceImpl interface.
func (s *APIServiceImpl) TestMap(ctx context.Context, req *api.Request) (resp map[string]*api.Response, err error) {
	logs.CtxInfo(ctx, "req: %s", utils.ToJson(req))
	return map[string]*api.Response{
		"k1": NewMockResponse(ctx, req),
		"k2": NewMockResponse(ctx, req),
	}, nil
}

// TestIntMap implements the APIServiceImpl interface.
func (s *APIServiceImpl) TestIntMap(ctx context.Context, req *api.Request) (resp map[int64]*api.Response, err error) {
	logs.CtxInfo(ctx, "req: %s", utils.ToJson(req))
	return map[int64]*api.Response{
		1: NewMockResponse(ctx, req),
		2: NewMockResponse(ctx, req),
	}, nil
}

// TestString implements the APIServiceImpl interface.
func (s *APIServiceImpl) TestString(ctx context.Context, req *api.Request) (resp string, err error) {
	logs.CtxInfo(ctx, "req: %s", utils.ToJson(req))
	return "hello kitex", nil
}

func NewMockResponse(ctx context.Context, req interface{}) *api.Response {
	persistentValues := metainfo.GetAllPersistentValues(ctx)
	values := metainfo.GetAllValues(ctx)
	return &api.Response{
		Req: utils.StringPtr(utils.ToJson(req)),
		MetaInfo: utils.StringPtr(utils.ToJson(map[string]interface{}{
			"PersistentValues": persistentValues,
			"Values":           values,
		})),
		BaseResp: &base.BaseResp{},
	}
}
