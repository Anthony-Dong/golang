package dag

import (
	"errors"
	"fmt"

	"strings"
)

type NodeDependencyBuilder struct {
	cur   *Node
	errs  []error
	graph *Graph
}

func NewNodeDependencyBuilder(g *Graph) *NodeDependencyBuilder {
	return &NodeDependencyBuilder{
		graph: g,
	}
}

func (b *NodeDependencyBuilder) isErr() bool {
	return len(b.errs) > 0
}

func (b *NodeDependencyBuilder) Node(name string) *NodeDependencyBuilder {
	if b.errs != nil {
		return b
	}
	b.cur = b.graph.GetNode(name)
	if b.cur == nil {
		b.errs = append(b.errs, fmt.Errorf(`not found node: %s`, name))
	}
	return b
}

// From  D.From(B, C)  D runs after B and C
func (b *NodeDependencyBuilder) From(nodes ...string) *NodeDependencyBuilder {
	if b.isErr() {
		return b
	}
	for _, node := range nodes {
		nn := b.graph.GetNode(node)
		if nn == nil {
			b.errs = append(b.errs, fmt.Errorf(`not found node: %s`, node))
			continue
		}
		nn.To = append(nn.To, b.cur)
	}
	return b
}

// To A.To(B, C)  A run before B and C
func (b *NodeDependencyBuilder) To(nodes ...string) *NodeDependencyBuilder {
	if b.isErr() {
		return b
	}
	for _, node := range nodes {
		nn := b.graph.GetNode(node)
		if nn == nil {
			b.errs = append(b.errs, fmt.Errorf(`not found node: %s`, node))
			continue
		}
		b.cur.To = append(b.cur.To, nn)
	}
	return b
}

func (b *NodeDependencyBuilder) Downgrade(nodes ...string) *NodeDependencyBuilder {
	if b.isErr() {
		return b
	}
	if b.cur.Downgrade == nil {
		b.cur.Downgrade = map[string]bool{}
	}
	for _, node := range nodes {
		nn := b.graph.GetNode(node)
		if nn == nil {
			b.errs = append(b.errs, fmt.Errorf(`not found node: %s`, node))
			continue
		}
		b.cur.Downgrade[node] = true
	}
	return b
}

func (b *NodeDependencyBuilder) Build() error {
	if err := b.graph.Build(); err != nil {
		b.errs = append(b.errs, err)
	}
	if !b.isErr() {
		return nil
	}
	errs := make([]string, 0, len(b.errs))
	for _, err := range b.errs {
		errs = append(errs, err.Error())
	}
	return errors.New(strings.Join(errs, "\n"))
}
