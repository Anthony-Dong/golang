package dag

import (
	"context"
	"fmt"
)

type BizContext interface {
	Context() context.Context
	Data() interface{}
}

type Node struct {
	Name string

	// one of run func
	Run     func(ctx context.Context) error            // NodeRunTypeNormal
	RunCond func(ctx context.Context) (string, error)  // NodeRunTypeCond
	RunGen  func(ctx context.Context) ([]*Node, error) // NodeRunTypeGen

	To []*Node

	From      []*Node         // use graph.revert func build it
	Downgrade map[string]bool // enable downgrade from node
}

type NodeRunType uint8

const (
	NodeRunTypeEmpty NodeRunType = iota
	NodeRunTypeNormal
	NodeRunTypeCond
	NodeRunTypeGen
)

func (n *Node) GetRunType() NodeRunType {
	if n.Run != nil {
		return NodeRunTypeNormal
	}
	if n.RunCond != nil {
		return NodeRunTypeCond
	}
	if n.RunGen != nil {
		return NodeRunTypeGen
	}
	return NodeRunTypeEmpty
}

func ContainsNode(nodes []*Node, name string) bool {
	for _, node := range nodes {
		if node.Name == name {
			return true
		}
	}
	return false
}

func MustNode(name string, run interface{}) *Node {
	node, err := NewNode(name, run)
	if err != nil {
		panic(err)
	}
	return node
}

func NewNode(name string, run interface{}) (*Node, error) {
	if name == "" {
		return nil, fmt.Errorf(`node name is empty`)
	}
	if name == beginNode || name == endNode {
		return nil, fmt.Errorf(`node name cannot be %s or %s`, beginNode, endNode)
	}
	if run == nil {
		return nil, fmt.Errorf(`node %s run func is nil`, name)
	}
	node := &Node{Name: name}
	switch r := run.(type) {
	case func(ctx context.Context) error:
		node.Run = r
	case func(ctx context.Context) (string, error):
		node.RunCond = r
	case func(ctx context.Context) ([]*Node, error):
		node.RunGen = r
	case func(ctx context.Context):
		node.Run = func(ctx context.Context) error {
			r(ctx)
			return nil
		}
	case func(ctx context.Context) string:
		node.RunCond = func(ctx context.Context) (string, error) {
			return r(ctx), nil
		}
	case func(ctx context.Context) []*Node:
		node.RunGen = func(ctx context.Context) ([]*Node, error) {
			return r(ctx), nil
		}
	default:
		return nil, fmt.Errorf(`unsupport node run type %T`, run)
	}
	return node, nil
}
