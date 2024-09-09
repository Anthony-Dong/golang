package dag

import (
	"bytes"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"sync"
)

const beginNode = "__begin__"
const endNode = "__end__"

type Graph struct {
	name  string
	nodes map[string]*Node // all nodes

	begin *Node
	end   *Node

	buildOnce sync.Once
	buildErr  error

	revert bool // 反向申明依赖，相当于只申明Node.From
}

func NewGraph(opts ...func(*Graph)) *Graph {
	output := &Graph{}
	for _, op := range opts {
		op(output)
	}
	return output
}

func WithName(name string) func(graph *Graph) {
	return func(graph *Graph) {
		graph.name = name
	}
}

func WithRevert(revert bool) func(graph *Graph) {
	return func(graph *Graph) {
		graph.revert = revert
	}
}

func (g *Graph) GetNode(name string) *Node {
	return g.nodes[name]
}

func (g *Graph) Graphviz() string {
	builder := bytes.NewBuffer(nil)
	name := g.name
	if name == "" {
		name = "G"
	}
	builder.WriteString("digraph ")
	builder.WriteString(strconv.Quote(name))
	builder.WriteString(" {\n")
	for _, node := range g.nodes {
		for _, to := range node.To {
			builder.WriteString(fmt.Sprintf("  %q -> %q ;\n", node.Name, to.Name))
		}
	}
	builder.WriteString("}\n")
	return builder.String()
}

func (g *Graph) toSliceNode() []*Node {
	result := make([]*Node, 0, len(g.nodes))
	for _, elem := range g.nodes {
		result = append(result, elem)
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].Name < result[j].Name
	})
	return result
}

func (g *Graph) init() {
	if g.nodes == nil {
		g.nodes = map[string]*Node{}
	}
}

func (g *Graph) MustNode(name string, run interface{}) {
	node, err := NewNode(name, run)
	if err != nil {
		panic(err)
	}
	if err := g.AddNode(node); err != nil {
		panic(err)
	}
}

func (g *Graph) addNode(name string, run interface{}) error {
	node, err := NewNode(name, run)
	if err != nil {
		return err
	}
	g.init()
	g.nodes[name] = node
	return nil
}

func (g *Graph) AddNode(node *Node) error {
	g.init()
	if node == nil {
		return fmt.Errorf("cannot add nil node")
	}
	if _, isExist := g.nodes[node.Name]; isExist {
		return fmt.Errorf(`node [%s] is exist`, node.Name)
	}
	g.nodes[node.Name] = node
	return nil
}

func (g *Graph) Build() error {
	g.buildOnce.Do(func() {
		g.buildErr = g._build()
	})
	return g.buildErr
}

func (g *Graph) checkCycles(route *Route, node *Node) error {
	if route.Contains(node.Name) {
		return fmt.Errorf("cycle: %s", strings.Join(append(route.route, node.Name), " -> "))
	}
	route.Push(node.Name)
	for _, elem := range node.To {
		if err := g.checkCycles(route, elem); err != nil {
			return err
		}
		route.Pop()
	}
	return nil
}

func (g *Graph) _build() error {
	var beginNodes []*Node
	var endNodes []*Node

	if g.revert {
		g.fillTo(nil, g.toSliceNode())
	} else {
		g.fillFrom(nil, g.toSliceNode())
	}

	// 入度不为0的Node 取反就是 begin node
	toNodes := make(map[string]bool)
	for _, node := range g.nodes {
		for _, to := range node.To {
			toNodes[to.Name] = true
		}
		// 出度为0的Node为end node
		if len(node.To) == 0 {
			endNodes = append(endNodes, node)
		}
	}

	if len(g.nodes) == len(toNodes) {
		for _, elem := range g.nodes {
			route := Route{}
			if err := g.checkCycles(&route, elem); err != nil {
				return err
			}
		}
		return fmt.Errorf(`cycle error`)
	}

	for _, node := range g.nodes {
		if toNodes[node.Name] {
			continue
		}
		beginNodes = append(beginNodes, node)
	}

	// 美化 beginNodes / endNodes 合成一个
	g.beautify(beginNodes, endNodes)
	return nil
}

func (g *Graph) beautify(beginNodes []*Node, endNodes []*Node) {
	if len(beginNodes) > 1 {
		virtualNode := &Node{Name: beginNode}
		for _, node := range beginNodes {
			virtualNode.To = append(virtualNode.To, node)
			node.From = append(node.From, virtualNode)
		}
		g.nodes[virtualNode.Name] = virtualNode
		g.begin = virtualNode
	} else {
		g.begin = beginNodes[0]
	}

	endNodeIsGenType := func() bool {
		return endNodes[0].GetRunType() == NodeRunTypeGen
	}
	if len(endNodes) > 1 || endNodeIsGenType() {
		virtualNode := &Node{Name: endNode}
		for _, node := range endNodes {
			virtualNode.From = append(virtualNode.From, node)
			node.To = append(node.To, virtualNode)
		}
		g.nodes[virtualNode.Name] = virtualNode
		g.end = virtualNode
	} else {
		g.end = endNodes[0]
	}
}

func (g *Graph) fillFrom(from *Node, nodes []*Node) {
	for _, node := range nodes {
		if from != nil && !ContainsNode(node.From, from.Name) {
			node.From = append(node.From, from)
		}
		g.fillFrom(node, node.To)
	}
}

func (g *Graph) fillTo(to *Node, nodes []*Node) {
	for _, node := range nodes {
		if to != nil && !ContainsNode(node.To, to.Name) {
			node.To = append(node.To, to)
		}
		g.fillTo(node, node.From)
	}
}
