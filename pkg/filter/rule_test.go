package parser

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParse(t *testing.T) {
	input := `a='1' and b='2' and (c=3 and (d='4' or f='5') and (x='6' or y='7' and z in (1,2)))`
	parseRule, err := ParseRule(input)
	if err != nil {
		t.Fatal(err)
	}
	marshal, _ := json.MarshalIndent(parseRule, "", "\t")
	t.Log(string(marshal))
}

func TestRule(t *testing.T) {
	intput := MapInput{
		"a": "1",
		"b": "2",
		"c": "3",
		"d": "4",
		"e": "5.5",
		"f": "xiaoming",
	}
	rule, err := ParseRule(`a = 1 and b = 2 and ( c = 4 or ( d = 4 and e in ( 5.5 , 6.6 ) ) ) and f = "xiaoming"`)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, rule.Result(intput), true)
}

func BenchmarkRule(b *testing.B) {
	b.StopTimer()
	intput := MapInput{
		"a": "1",
		"b": "2",
		"c": "3",
		"d": "4",
		"e": "5.5",
		"f": "xiaoming",
	}
	parseRule, err := ParseRule(`a = 1 and b = 2 and ( c = 4 or ( d = 4 and e in ( 5.5 , 6.6 ) ) ) and f = "xiaoming"`)
	if err != nil {
		b.Fatal(err)
	}
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		if !parseRule.Result(intput) {
			b.Fatal("must ture")
		}
	}
}
