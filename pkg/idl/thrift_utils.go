package idl

import (
	"fmt"
	"regexp"

	"github.com/cloudwego/thriftgo/parser"

	"github.com/cloudwego/kitex/pkg/generic"
)

var (
	// todo: 目前移除了全部的注解，如果有额外需求，保留一些注解，可能需要把 \w+\.\w+ 改成 字符组的方式.
	_thriftApiMethodRegexp = regexp.MustCompile(`\((\s*\w+\.\w+\s*=\s*(".+"|'.+')\s*,?\s*)+\)`)
	_thriftHashMapRegexp   = regexp.MustCompile(`(hash_map)(\s*<[\w\s._,<>]+>)`)
	//apiMethodRegexp = regexp.MustCompile(`\(\s*api\.(post|get|put|delete|patch)\s*=\s*("\S+"|'\S+')\s*\)`)
)

func ParseThrift(filename string, content string) (*parser.Thrift, error) {
	return parser.ParseString(filename, content)
}

func _fixThriftIDLForKitex(idlContent string) string {
	idlContent = _thriftApiMethodRegexp.ReplaceAllString(idlContent, "")
	return _thriftHashMapRegexp.ReplaceAllString(idlContent, "map$2") // 修复kitex不支持hash_map
}

func fixThriftIDLForKitex(idls map[string]string) map[string]string {
	result := make(map[string]string, len(idls))
	for k, v := range idls {
		result[k] = _fixThriftIDLForKitex(v)
	}
	return result
}

func loadThriftDescriptorProviderV1(main string) (_ generic.DescriptorProvider, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf(`parse(v1) thrift idl find err: %v`, r)
		}
	}()
	provider, err := generic.NewThriftFileProvider(main)
	if err != nil {
		return nil, fmt.Errorf("parse(v1) thrift idl find err: %v", err)
	}
	return provider, nil
}

func loadThriftDescriptorProviderV2(main string, idls map[string]string) (_ generic.DescriptorProvider, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf(`parse(v2) thrift idl find err: %v`, r)
		}
	}()
	provider, err := generic.NewThriftContentWithAbsIncludePathProvider(main, idls)
	if err != nil {
		return nil, fmt.Errorf("parse(v2) thrift idl find err: %v", err)
	}
	return provider, nil
}