package dag

import (
	"fmt"
)

type FunctionInfo struct {
	Name    string
	Inputs  []string // 函数的输入
	Outputs []string // 函数的输出
	Run     interface{}
}

type FunctionBuilder struct {
	functions map[string]*FunctionInfo
}

func NewFunctionBuilder() *FunctionBuilder {
	return &FunctionBuilder{functions: map[string]*FunctionInfo{}}
}

func (f *FunctionBuilder) AddFunction(info *FunctionInfo) error {
	if f.functions == nil {
		f.functions = map[string]*FunctionInfo{}
	}
	if info.Name == "" {
		return fmt.Errorf("function name is nil")
	}
	if f.functions[info.Name] != nil {
		return fmt.Errorf("function %s already exists", info.Name)
	}
	f.functions[info.Name] = info
	return nil
}

func (f *FunctionBuilder) AddFunctions(infos []*FunctionInfo) error {
	for _, info := range infos {
		if err := f.AddFunction(info); err != nil {
			return err
		}
	}
	return nil
}

func contains(inputs []string, item string) bool {
	for _, elem := range inputs {
		if elem == item {
			return true
		}
	}
	return false
}

// NewGraph 根据输入参数 + 已经申明的函数 推导出你需要的Graph
func (f *FunctionBuilder) NewGraph(name string, inputs []string) (*Graph, error) {
	context := buildContext{
		FunctionBuilder: f,
		Graph:           NewGraph(WithName(name), WithRevert(true)),
	}
	for _, input := range inputs {
		search := false
		for _, function := range f.functions {
			if contains(function.Outputs, input) {
				if _, err := context.buildFunction(function); err != nil {
					return nil, err
				}
				search = true
			}
		}
		if !search {
			return nil, fmt.Errorf("graph %q input %q not found", name, input)
		}
	}
	if err := context.Graph.Build(); err != nil {
		return nil, err
	}
	return context.Graph, nil
}

type buildContext struct {
	*FunctionBuilder
	*Graph
}

func (ctx *buildContext) buildFunction(function *FunctionInfo) (*Node, error) {
	if node := ctx.GetNode(function.Name); node != nil {
		return node, nil
	}
	node, err := NewNode(function.Name, function.Run)
	if err != nil {
		return nil, err
	}
	if err := ctx.AddNode(node); err != nil {
		return nil, err
	}
	if err := ctx.buildDeps(node, function.Inputs); err != nil {
		return nil, err
	}
	return node, nil
}

func (ctx *buildContext) buildDeps(node *Node, inputs []string) error {
	for _, input := range inputs {
		isFound := false
		for _, function := range ctx.functions {
			if contains(function.Outputs, input) {
				depNode, err := ctx.buildFunction(function)
				if err != nil {
					return err
				}
				node.From = append(node.From, depNode)
				isFound = true
			}
		}
		if !isFound {
			return fmt.Errorf("function %q input %q not found", node.Name, input)
		}
	}
	return nil
}
