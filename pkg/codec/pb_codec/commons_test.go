package pb_codec

import (
	"encoding/json"
	"testing"
)

func TestFieldOrderMap_MarshalJSON(t *testing.T) {
	orderMap := NewFieldOrderMap(10)

	orderMap.Set(NewField(1, 1), "1")
	orderMap.Set(NewField(2, 2), 2)

	marshal, _ := json.Marshal(orderMap)

	t.Log(string(marshal))
}
