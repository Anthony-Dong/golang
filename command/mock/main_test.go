package mock

import (
	"path/filepath"
	"testing"

	"github.com/anthony-dong/golang/pkg/utils"
)

func Test_mockThriftData(t *testing.T) {
	dir := utils.GetGoProjectDir()
	v, err := GetThriftMockData(filepath.Join(dir, "pkg/idl/test/api.thrift"), "Request")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(utils.ToJson(v))
}
