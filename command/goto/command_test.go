package _goto

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseGoLink(t *testing.T) {
	url := ParseGoLink("github.com/Anthony-Dong/golang@v0.0.10/cli/devtool/main.go:100")
	t.Log(url)
	assert.Equal(t, url, "https://github.com/Anthony-Dong/golang/blob/v0.0.10/cli/devtool/main.go#L100")
}

func Test_replaceUrlTag(t *testing.T) {
	assert.Equal(t, replaceUrlTag("xxx@v1.0.1"), "xxx/v1.0.1/xxx")
}
