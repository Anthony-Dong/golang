package parser

// https://github.com/mna/pigeon
//go:generate pigeon -o grammar.peg.go ./grammar.peg

import (
	"fmt"
)

type Conditions struct {
	Group   *Conditions `json:"group,omitempty"`
	Value   *Condition  `json:"value,omitempty"`
	Logical string      `json:"logical,omitempty"`
	Next    *Conditions `json:"next,omitempty"`
}

type Condition struct {
	Key      string   `json:"key,omitempty"`
	Operator string   `json:"operator,omitempty"`
	Value    string   `json:"value,omitempty"`
	Values   []string `json:"values,omitempty"`
}

type MapInput map[string]string

func (m MapInput) Get(key string) string {
	return m[key]
}

type Input interface {
	Get(key string) string
}

func (c *Conditions) Result(intput Input) bool {
	var result bool
	if c.Value != nil {
		result = c.Value.Result(intput)
	} else if c.Group != nil {
		result = c.Group.Result(intput)
	} else {
		result = true // 什么也没定义
	}
	if c.Next == nil {
		return result
	}
	if c.Logical == "and" {
		return result && c.Next.Result(intput)
	}
	return result || c.Next.Result(intput)
}

func (c *Condition) Result(input Input) bool {
	value := input.Get(c.Key)
	switch c.Operator {
	case "in":
		return contains(c.Values, value)
	case "=":
		return value == c.Value
	case "!=":
		return value != c.Value
	case ">":
		return value > c.Value
	case "<":
		return value < c.Value
	case ">=":
		return value >= c.Value
	case "<=":
		return value <= c.Value
	default:
		panic(fmt.Sprintf(`invalid operator (%s)`, c.Operator))
	}
}

func ParseRule(rule string) (*Conditions, error) {
	parse, err := Parse("", []byte(rule))
	if err != nil {
		return nil, err
	}
	return parse.(*Conditions), nil
}

func contains[T comparable](arr []T, data T) bool {
	for _, elem := range arr {
		if elem == data {
			return true
		}
	}
	return false
}
