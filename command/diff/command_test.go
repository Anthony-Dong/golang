package diff

import "testing"

func Test_replaceGJson2Jq(t *testing.T) {
	t.Log(jqPath2GJson(".data.draft.plugins.[0].tasks.[0].base"))
}
