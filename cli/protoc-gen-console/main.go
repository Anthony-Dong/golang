package main

import (
	"io"
	"os"
	"strings"

	"github.com/anthony-dong/golang/pkg/utils"

	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/pluginpb"
)

func protoMessageToJson(m proto.Message) ([]byte, error) {
	return protojson.MarshalOptions{Multiline: true, AllowPartial: true}.Marshal(m)
}

func main() {
	req := pluginpb.CodeGeneratorRequest{}
	stdIn, err := io.ReadAll(os.Stdin)
	if err != nil {
		panic(err)
	}
	if err := proto.Unmarshal(stdIn, &req); err != nil {
		panic(err)
	}
	params := parseParams(req.Parameter)
	output := params["output"]
	if output == "" {
		output = "protoc.json"
	}
	disableSourceCode := params["disable_source_code"]
	if disableSourceCode == "1" {
		for _, file := range req.ProtoFile {
			file.SourceCodeInfo = nil
		}
	}
	json, err := protoMessageToJson(&req)
	if err != nil {
		panic(err)
	}

	resp := pluginpb.CodeGeneratorResponse{
		File: []*pluginpb.CodeGeneratorResponse_File{
			{
				Name:    utils.StringPtr(output),
				Content: utils.StringPtr(utils.Bytes2String(json)),
			},
		},
	}
	if stdout, err := proto.Marshal(&resp); err != nil {
		panic(err)
	} else {
		if _, err := os.Stdout.Write(stdout); err != nil {
			panic(err)
		}
	}
}

// parseParams
// output=output/out.json,disable_source_code=1,params1=1,params2=2
func parseParams(paramsPtr *string) map[string]string {
	if paramsPtr == nil || *paramsPtr == "" {
		return map[string]string{}
	}
	result := make(map[string]string)
	split := strings.Split(*paramsPtr, ",")
	for _, elem := range split {
		kv := strings.SplitN(elem, "=", 2)
		if len(kv) == 2 {
			result[kv[0]] = kv[1]
		} else {
			result[kv[0]] = ""
		}
	}
	return result
}
