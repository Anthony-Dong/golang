package dag

import (
	"regexp"
	"strconv"
	"testing"

	"github.com/anthony-dong/golang/pkg/utils"
)

/**
digraph "测试场景" {
  "ServiceVersionLoader" -> "ServiceIDLAssembler" ;
  "ServiceVersionLoader" -> "ServiceIdlLoader" ;
  "ServiceIdlLoader" -> "ServiceIDLAssembler" ;
}
*/

func parseDependency(dependency string) ([]map[string]string, error) {
	compile := regexp.MustCompile(`(['"]?[\w._-]+['"]?)\s+->\s+(['"]?[\w._-]+['"]?)\s*(\[[\w\s._=-]+])?`)
	allString := compile.FindAllStringSubmatch(dependency, -1)
	result := make([]map[string]string, 0)
	unquote := func(input string) string {
		unquote, err := strconv.Unquote(input)
		if err != nil {
			return input
		}
		return unquote
	}
	for _, elem := range allString {
		if len(elem) != 4 {
			continue
		}
		result = append(result, map[string]string{
			"from": unquote(elem[1]),
			"to":   unquote(elem[2]),
		})
	}
	return result, nil
}

func TestNewNodeDependencyBuilder(t *testing.T) {
	r, err := parseDependency(`digraph "测试场景" {
  "ServiceVersionLoader" -> "ServiceIDLAssembler" ;
  "ServiceVersionLoader" -> "ServiceIdlLoader" ;
  "ServiceIdlLoader" -> "ServiceIDLAssembler" ;
}`)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(utils.ToJson(r, true))
}
