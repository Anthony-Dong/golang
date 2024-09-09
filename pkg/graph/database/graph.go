package graph

import (
	"fmt"
	"strings"
)

type Property struct {
	Key   string
	Value interface{}
}

func NewProperty(key string, value interface{}) *Property {
	return &Property{Key: key, Value: value}
}

type Vertex struct {
	ID   int64 // 顶点ID
	Type int32 // 顶点类型

	Properties []*Property // 顶点属性
}

func (v *Vertex) Key() string {
	return newVertexKey(v.ID, v.Type)
}

func newVertexKey(id int64, Type int32) string {
	return fmt.Sprintf(`%d$%d`, id, Type)
}

type Edge struct {
	Type int32 // 边类型

	OutV *Vertex // 出 from
	InV  *Vertex // 入 to

	Properties []*Property // 边属性
}

func (v *Edge) Key() string {
	return newEdgeKey(v.Type, v.OutV.Key(), v.InV.Key())
}

func newEdgeKey(Type int32, outV string, inV string) string {
	return fmt.Sprintf(`%d$%s$%s`, Type, outV, inV)
}

type Graph struct {
	Name     string
	Edges    map[string]*Edge
	Vertexes map[string]*Vertex

	graph       map[string][]string // 正向
	revertGraph map[string][]string // 反向
}

func (g *Graph) AddVertex(id int64, Type int32, properties ...*Property) {
	if g.Vertexes == nil {
		g.Vertexes = map[string]*Vertex{}
	}
	v := &Vertex{
		ID:         id,
		Type:       Type,
		Properties: properties,
	}
	if g.Vertexes[v.Key()] == nil {
		g.Vertexes[v.Key()] = v
	}
}

func (g *Graph) GetVertex(id VertexKey) *Vertex {
	return g.Vertexes[fmt.Sprintf("%d$%d", id.ID, id.Type)]
}

func (g *Graph) AddEdge(outV VertexKey, inV VertexKey, Type int32, properties ...*Property) {
	edge := &Edge{OutV: g.GetVertex(outV), InV: g.GetVertex(inV), Type: Type, Properties: properties}
	if g.Edges == nil {
		g.Edges = map[string]*Edge{}
	}
	if g.Edges[edge.Key()] == nil {
		g.Edges[edge.Key()] = edge
	}
	if g.graph == nil {
		g.graph = map[string][]string{}
	}
	if g.revertGraph == nil {
		g.revertGraph = map[string][]string{}
	}
	// 添加正向索引(临接表) out->in
	g.graph[edge.OutV.Key()] = append(g.graph[edge.OutV.Key()], edge.InV.Key())
	// 添加反向索引(临接表) in->out
	g.revertGraph[edge.InV.Key()] = append(g.revertGraph[edge.InV.Key()], edge.OutV.Key())
}

func (g *Graph) Both(vertex VertexKey, Type int32) []*Edge {
	result := make([]*Edge, 0)
	result = append(result, g.Out(vertex, Type)...)
	result = append(result, g.In(vertex, Type)...)
	return result
}

func (g *Graph) Out(outVertex VertexKey, Type int32) []*Edge {
	outKey := outVertex.Key()
	result := make([]*Edge, 0)
	for _, inKey := range g.graph[outKey] {
		result = append(result, g.Edges[newEdgeKey(Type, outKey, inKey)])
	}
	return result
}

func (g *Graph) In(inVertex VertexKey, Type int32) []*Edge {
	inKey := inVertex.Key()
	result := make([]*Edge, 0)
	for _, outKey := range g.revertGraph[inKey] {
		result = append(result, g.Edges[fmt.Sprintf(`%d$%s$%s`, Type, outKey, inKey)])
	}
	return result
}

func (g *Graph) Print() string {
	buildLabel := func(pros []*Property) string {
		tags := make([]string, 0)
		for _, pro := range pros {
			tags = append(tags, fmt.Sprintf("%s=%s", pro.Key, pro.Value))
		}
		tagsStr := ""
		if len(tags) > 0 {
			tagsStr = fmt.Sprintf("[lable=%q]", strings.Join(tags, " "))
		}
		return tagsStr
	}
	builder := strings.Builder{}
	builder.WriteString(fmt.Sprintf("digraph %q {\n", g.Name))

	for _, elem := range g.Vertexes {
		builder.WriteString(fmt.Sprintf("%d %s\n", elem.ID, buildLabel(elem.Properties)))
	}

	for _, elem := range g.Edges {
		builder.WriteString(fmt.Sprintf("%d -> %d %s\n", elem.OutV.ID, elem.InV.ID, buildLabel(elem.Properties)))
	}
	builder.WriteString("}")
	return builder.String()
}

type VertexKey struct {
	ID   int64
	Type int32
}

func (v *VertexKey) Key() string {
	return newVertexKey(v.ID, v.Type)
}

func NewVertexKey(id int64, Type int32) VertexKey {
	return VertexKey{ID: id, Type: Type}
}
