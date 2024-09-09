package idl

import (
	"fmt"
	"regexp"

	"github.com/cloudwego/kitex/pkg/generic/descriptor"

	"github.com/anthony-dong/golang/pkg/idl/internal"

	"github.com/cloudwego/thriftgo/parser"

	"github.com/cloudwego/kitex/pkg/generic"
)

var (
	_thriftApiMethodRegexp = regexp.MustCompile(`\((\s*\w+\.\w+\s*=\s*(".+"|'.+')\s*,?\s*)+\)`)
	_thriftHashMapRegexp   = regexp.MustCompile(`(hash_map)(\s*<[\w\s._,<>]+>)`)
)

func ParseThrift(filename string, content string) (*parser.Thrift, error) {
	return parser.ParseString(filename, content)
}

func _fixThriftIDLForKitex(idlContent string) string {
	idlContent = _thriftApiMethodRegexp.ReplaceAllString(idlContent, "")
	return _thriftHashMapRegexp.ReplaceAllString(idlContent, "map$2") // kitex不支持hash_map
}

func fixThriftIDLForKitex(idls map[string]string) map[string]string {
	result := make(map[string]string, len(idls))
	for k, v := range idls {
		result[k] = _fixThriftIDLForKitex(v)
	}
	return result
}

func loadThriftDescriptorProvider(main string, idls map[string]string) (provider generic.DescriptorProvider, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf(`parse thrift idl find err: %v`, r)
		}
	}()
	provider, err = generic.NewThriftContentWithAbsIncludePathProvider(main, idls)
	if err != nil {
		return nil, fmt.Errorf("parse thrift idl find err: %v", err)
	}
	return provider, nil
}

func ParseThriftType(tree *parser.Thrift, t *parser.Type) (*descriptor.TypeDescriptor, error) {
	return internal.ParseThriftType(tree, t)
}
