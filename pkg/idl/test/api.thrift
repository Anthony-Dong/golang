include "base/base.thrift"
include "common/common.thrift"

struct Request {
     1: optional string                    Field1  // @desc: 字段1的备注信息; @example: "hello, string!"; @name: 字段1
     2: optional i32                       Field2  // @desc: 字段2的备注信息; @example: 1; @name: 字段2
     3: optional bool                      Field3  // @desc: 字段3的备注信息; @example: true; @name: 字段3
     4: optional list<string>              Field4  // @desc: 字段4的备注信息; @example: ["hello, list"]; @name: 字段4
     5: optional list<map<string, string>> Field5  // @desc: 字段5的备注信息，测试list+map组合; @example: [{"k1": "v1"}]; @name: 字段5
     6: optional map<string, i64>          Field6  // @desc: 字段6的备注信息, 测试map; @example: {"k1": 1}; @name: 字段6
     7: optional list<Message>             Field7  // @desc: 字段7的备注信息, 测试list; @name: 字段7
     8: optional map<string, Message>      Field8  // @desc: 字段8的备注信息, 测试map+结构体; @example: {"Field8": ""}; @name: 字段8
     9: optional double                    Field9  // @desc: 字段9的备注信息, 测试浮点类型; @example: 3.141592653; @name: 字段9
    10: optional string                    Field10 // @desc: 字段10的备注信息; @example: abc; @name: 字段10
    11: optional string                    Field11 // @desc: 字段11; @example: abcdefg
    12: optional common.Status             Field12 // @desc: 字段12; @example: 1
    13: common.KVList                      Field13 // @desc: 字段13
    255: optional base.Base Base // base
}

struct Message {
    1: optional string Field1 // @example: hello Field1; @name: 消息字段1
    2: optional i32    Field2 // @example: 10086
}

struct Response {
    1: optional string              Req       // @desc: 请求原路返回; @name: 请求信息
    2: optional string              MetaInfo  // @desc: 请求详情

    255: optional base.BaseResp BaseResp // base
}


service APIService {
    Response RPCAPI1 (1: Request req) (api.get='/api/v1/query'),
    Response RPCAPI2 (1: Request req) (api.get='/api/v1/query'),
}