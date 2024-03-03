package utils

import (
	"encoding/json"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestFloat642String(t *testing.T) {
	assert.Equal(t, Float642String(1.001, 2), "1.00")
	assert.Equal(t, Float642String(1.006, 2), "1.01")
	assert.Equal(t, Float642String(1.005, 2), "1.00")
	assert.Equal(t, Float642String(1.004, 2), "1.00")
}

func TestString(t *testing.T) {
	pt, err := strconv.Atoi("")
	t.Log(pt, err)
}

func TestJson(t *testing.T) {
	type JsonD struct {
		Time JsonDuration
	}
	{
		data := JsonD{Time: NewJsonDuration(time.Second * 1000)}
		marshal, err := json.Marshal(data)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, string(marshal), `{"Time":"16m40s"}`)
		data.Time = 0
		if err := json.Unmarshal(marshal, &data); err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, data.Time, NewJsonDuration(time.Second*1000))
	}

	{
		data := JsonD{}
		if err := json.Unmarshal([]byte(`{"Time":""}`), &data); err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, data.Time, NewJsonDuration(0))
	}

	{
		data := JsonD{}
		if err := json.Unmarshal([]byte(`{"Time":10000}`), &data); err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, data.Time, NewJsonDuration(10000))
	}

	{
		data := JsonD{}
		if err := json.Unmarshal([]byte(`{"Time":null}`), &data); err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, data.Time, NewJsonDuration(0))
	}
}
