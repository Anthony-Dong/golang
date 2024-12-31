package diff

import (
	"testing"

	"github.com/anthony-dong/golang/pkg/utils"
)

func TestDiffJson(t *testing.T) {
	t.Run("case1", func(t *testing.T) {
		jsonString, err := DiffJsonString(`{"k1": "v1", "k2": "v2"}`, `{"k1": "v1", "k2": "v3"}`)
		if err != nil {
			t.Fatal(err)
		}
		t.Log(utils.ToJson(jsonString, true))
	})

	t.Run("case2", func(t *testing.T) {
		jsonString, err := DiffJsonString(`{"k1": "v1", "k2": ["1","2"]}`, `{"k1": "v1", "k2": ["1"]}`)
		if err != nil {
			t.Fatal(err)
		}
		t.Log(utils.ToJson(jsonString, true))
	})

	t.Run("case3", func(t *testing.T) {
		jsonString, err := DiffJsonString(`{"k1": "v1", "k2": [{"k1":"v1"}, {"k2": "v2"}]}`, `{"k1": "v1", "k2": [{"k1":"v2"}, {"k3": "v2"}]}`)
		if err != nil {
			t.Fatal(err)
		}
		t.Log(utils.ToJson(jsonString, true))
	})
}
