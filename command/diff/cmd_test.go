package diff

import (
	"testing"

	"github.com/anthony-dong/golang/pkg/utils"
)

type JsonData struct {
	Key    string
	Value  string `json:"value"`
	Value2 string `json:"value2"`
}

func TestName(t *testing.T) {
	js := JsonData{Value: "1", Value2: "2"}
	t.Log(utils.ToJson(js))
	//parse, err := gojq.Parse(".")
	//if err != nil {
	//	return err
	//}
	//parse.Run()
}
