package graph

import (
	"fmt"
	"runtime/debug"
	"testing"

	"github.com/anthony-dong/golang/pkg/utils"
)

func TestName(t *testing.T) {
	graph := Graph{}

	// g.addV().setProperty("id",1).setProperty("type",1).setProperty("step","stg_orgs")
	graph.AddVertex(1, 1, NewProperty("step", "stg_orgs"))
	graph.AddVertex(2, 1, NewProperty("step", "stg_users"))
	graph.AddVertex(3, 1, NewProperty("step", "stg_user_groups"))
	graph.AddVertex(4, 1, NewProperty("step", "init_users"))
	graph.AddVertex(5, 1, NewProperty("step", "dim_users"))

	// g.addE("follow").from(1, 1).to(2,1)
	graph.AddEdge(NewVertexKey(1, 1), NewVertexKey(2, 1), 1)
	graph.AddEdge(NewVertexKey(2, 1), NewVertexKey(4, 1), 1)
	graph.AddEdge(NewVertexKey(3, 1), NewVertexKey(4, 1), 1)
	graph.AddEdge(NewVertexKey(4, 1), NewVertexKey(5, 1), 1)

	// g.V().has("id", 4).has("type", 1).bothE(1)
	// 查询ID为4且类型为1，边类型为1的出入边

	fmt.Println(graph.Print())
	fmt.Println(utils.ToJson(graph.Out(NewVertexKey(4, 1), 1), true))
}

func TestName3(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Log("recover:", r)
			t.Log("stack:", string(debug.Stack()))
		}
	}()
	t.Log(foo(nil))
}

func foo(req *Request) (resp *Response, err error) {
	defer func() {
		for _, product := range resp.Product {
			fmt.Println(product)
		}
	}()
	return foo2(req)
}

type Response struct {
	Product []string
}

type Request struct {
	Data int
}

func foo2(req *Request) (*Response, error) {
	foo3(req)
	return nil, nil
}

func foo3(req *Request) {
	fmt.Println("req: ", req)
	_ = req.Data
}
