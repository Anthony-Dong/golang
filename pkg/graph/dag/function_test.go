package dag

import (
	"context"
	"fmt"
	"testing"
	"time"
)

func TestNewFunctionBuilder(t *testing.T) {
	builder := NewFunctionBuilder()
	err := builder.AddFunctions([]*FunctionInfo{
		{
			Name: "LoaderA",
			Outputs: []string{
				"ctx.a",
			},
			Run: func(ctx context.Context) {
				time.Sleep(time.Millisecond * 20)
				t.Log("LoaderA")
			},
		},
		{
			Name: "LoaderB",
			Outputs: []string{
				"ctx.b",
			},
			Run: func(ctx context.Context) {
				time.Sleep(time.Millisecond * 10)
				t.Log("LoaderB")
			},
		},
		{
			Name: "LoaderC",
			Inputs: []string{
				"ctx.a",
				"ctx.b",
			},
			Outputs: []string{
				"ctx.c",
			},
			Run: func(ctx context.Context) {
				t.Log("LoaderC")
			},
		},
		{
			Name:   "LoaderE",
			Inputs: []string{},
			Outputs: []string{
				"ctx.e",
			},
			Run: func(ctx context.Context) {
				t.Log("LoaderE")
			},
		},

		{
			Name: "AssemblerA",
			Inputs: []string{
				"ctx.c",
			},
			Outputs: []string{
				"resp.a",
			},
			Run: func(ctx context.Context) {
				t.Log("AssemblerA")
			},
		},
		{
			Name: "AssemblerB",
			Inputs: []string{
				"ctx.c",
			},
			Outputs: []string{
				"resp.b",
			},
			Run: func(ctx context.Context) {
				t.Log("AssemblerB")
			},
		},
		{
			Name: "AssemblerE",
			Inputs: []string{
				"ctx.e",
			},
			Outputs: []string{
				"resp.e",
			},
			Run: func(ctx context.Context) {
				t.Log("AssemblerE")
			},
		},
	})
	if err != nil {
		panic(err)
	}
	graph, err := builder.NewGraph("场景一", []string{
		"resp.b",
		"resp.a",
		"resp.e",
	})
	if err != nil {
		t.Fatal(err)
	}
	if err := graph.Build(); err != nil {
		t.Fatal(err)
	}
	fmt.Println(graph.Graphviz())
	if err := NewExecutor().Execute(context.TODO(), graph); err != nil {
		t.Fatal(err)
	}
}
