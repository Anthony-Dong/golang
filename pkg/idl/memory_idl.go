package idl

import (
	"io/ioutil"
	"path/filepath"

	"github.com/anthony-dong/golang/pkg/utils"

	"github.com/cloudwego/thriftgo/parser"
)

var _ MemoryIDLProvider = (*defaultMemoryIDLProvider)(nil)
var _ ThriftIDLProvider = (*defaultMemoryIDLProvider)(nil)

type defaultMemoryIDLProvider struct {
	Main string

	idls map[string]string
	ast  *parser.Thrift
}

func NewMemoryIDLProvider(main string) *defaultMemoryIDLProvider {
	return &defaultMemoryIDLProvider{
		Main: main,
	}
}

func (t *defaultMemoryIDLProvider) addLocalIDL(filename string, content string) {
	if t.idls == nil {
		t.idls = map[string]string{}
	}
	t.idls[filename] = content
}

func (t *defaultMemoryIDLProvider) lookup() error {
	if err := t.init(); err != nil {
		return err
	}
	if ast, err := t.parse(t.Main, nil); err != nil {
		return err
	} else {
		t.ast = ast
	}
	return nil
}

func (t *defaultMemoryIDLProvider) MemoryIDL() (*MemoryIDL, error) {
	if err := t.lookup(); err != nil {
		return nil, err
	}
	return &MemoryIDL{
		Main: t.Main,
		IDLs: t.idls,
	}, nil
}

func (t *defaultMemoryIDLProvider) ThriftIDL() (*parser.Thrift, error) {
	if err := t.lookup(); err != nil {
		return nil, err
	}
	return t.ast, nil
}

func (t *defaultMemoryIDLProvider) parse(filename string, walk map[string]*parser.Thrift) (*parser.Thrift, error) {
	if walk == nil {
		walk = map[string]*parser.Thrift{}
	}
	if thrift, isExist := walk[filename]; isExist {
		return thrift, nil
	}
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	ast, err := parser.ParseString(t.Main, utils.Bytes2String(content))
	if err != nil {
		return nil, err
	}
	walk[filename] = ast
	t.addLocalIDL(filename, utils.Bytes2String(content))
	for _, elem := range ast.Includes {
		if elem.GetUsed() {
			continue
		}
		includeAst, err := t.parse(thriftAbsPath(filename, elem.GetPath()), walk)
		if err != nil {
			return nil, err
		}
		elem.Used = utils.BoolPtr(true)
		elem.Reference = includeAst
	}
	return ast, nil
}

func thriftAbsPath(path, includePath string) string {
	if filepath.IsAbs(includePath) {
		return includePath
	}
	return filepath.Join(filepath.Dir(path), includePath)
}

func (t *defaultMemoryIDLProvider) init() error {
	t.idls = nil
	t.ast = nil
	if filepath.IsAbs(t.Main) {
		return nil
	}
	abs, err := filepath.Abs(t.Main)
	if err != nil {
		return err
	}
	t.Main = abs
	return nil
}
