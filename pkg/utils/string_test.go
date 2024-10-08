package utils

import (
	"regexp"
	"testing"
	"time"

	"gopkg.in/yaml.v3"

	"github.com/stretchr/testify/assert"
)

func TestSlug(t *testing.T) {
	t.Log(Slug("影師"))
}

func TestToString(t *testing.T) {
	assert.Equal(t, ToString(byte(1)), "1")

	assert.Equal(t, ToString(true), "true")
	assert.Equal(t, ToString(false), "false")

	assert.Equal(t, ToString(float64(1.11111)), "1.11111")
	assert.Equal(t, ToString(float64(1.000)), "1")
	assert.Equal(t, ToString(float64(1.001)), "1.001")
	assert.Equal(t, ToString(float64(-1.001)), "-1.001")

	assert.Equal(t, ToString(float32(1.11111)), "1.11111")
	assert.Equal(t, ToString(float32(1.000)), "1")
	assert.Equal(t, ToString(float32(1.001)), "1.001")
	assert.Equal(t, ToString(float32(-1.001)), "-1.001")

	assert.Equal(t, ToString(uint64(1)), "1")
	assert.Equal(t, ToString(uint32(1)), "1")
	assert.Equal(t, ToString(uint16(1)), "1")
	assert.Equal(t, ToString(uint8(1)), "1")

	assert.Equal(t, ToString(int64(-1)), "-1")
	assert.Equal(t, ToString(int32(-1)), "-1")
	assert.Equal(t, ToString(int16(-1)), "-1")
	assert.Equal(t, ToString(int8(-1)), "-1")

	assert.Equal(t, ToString(-1), "-1")

	now := time.Now()
	assert.Equal(t, ToString(now), now.String())

	type data struct {
		K1 string `json:"k1"`
	}
	assert.Equal(t, ToString(data{K1: "1"}), `{"k1":"1"}`)
}

func TestFormatFloat(t *testing.T) {
	assert.Equal(t, FormatFloat(1.1, 64), "1.1")
	assert.Equal(t, FormatFloat(1, 64), "1")
	assert.Equal(t, FormatFloat(1.0, 64), "1")
}

func TestContainsString(t *testing.T) {
	assert.Equal(t, ContainsString([]string{"1", "2"}, "2"), true)
	assert.Equal(t, ContainsString([]string{"1", "2"}, "3"), false)
}

func TestToPrettyJsonString(t *testing.T) {
	testData := map[string]interface{}{
		"k1": 1,
		"k2": "k2",
	}
	assert.Equal(t, ToJson(testData), `{"k1":1,"k2":"k2"}`)
}

func TestSplitSliceString(t *testing.T) {
	assert.Equal(t, SplitSliceString([]string{"1", "2", "3"}, 2), [][]string{{"1", "2"}, {"3"}})
	assert.Equal(t, SplitSliceString([]string{"1", "2", "3", "4"}, 2), [][]string{{"1", "2"}, {"3", "4"}})
	assert.Equal(t, SplitSliceString([]string{"1", "2"}, 2), [][]string{{"1", "2"}})
	assert.Equal(t, SplitSliceString([]string{"1", "2"}, 1), [][]string{{"1"}, {"2"}})
	assert.Equal(t, SplitSliceString([]string{}, 1), [][]string{})
	assert.Equal(t, SplitSliceString([]string{"1", "2"}, 3), [][]string{{"1", "2"}})
}

func TestLinesToString(t *testing.T) {
	assert.Equal(t, LinesToString([]string{"1", "2", "3"}), `1
2
3`)
	assert.Equal(t, LinesToString([]string{"1"}), `1`)

	assert.Equal(t, LinesToString([]string{""}), ``)
}

func TestNewString(t *testing.T) {
	assert.Equal(t, NewString('a', 0), "")
	assert.Equal(t, NewString('a', 1), "a")
	assert.Equal(t, NewString('a', 2), "aa")
}

func TestGenerateUUID(t *testing.T) {
	t.Log(GenerateUUID())
}

func TestUnsafeBytes(t *testing.T) {
	assert.Equal(t, String2Bytes("123"), []byte("123"))
	assert.Equal(t, Bytes2String([]byte("123")), "123")

	assert.Equal(t, String2Bytes(""), []byte(nil))
	assert.Equal(t, Bytes2String([]byte("")), "")

	data := []byte(nil)
	assert.Equal(t, len(data), 0)
	assert.Equal(t, cap(data), 0)
	data = append(data, 0)
	assert.Equal(t, data, []byte{0})
}

func TestSplitString(t *testing.T) {
	assert.Equal(t, SplitString(`hello,world, a, ,c`, ","), []string{"hello", "world", "a", "c"})
}

func TestTrimLeftSpace(t *testing.T) {
	assert.Equal(t, TrimLeftSpace("\t \n hello world \t"), "hello world \t")
	assert.Equal(t, TrimRightSpace("\t \n hello world \t\n\r "), "\t \n hello world")
}

func checkScaffoldRuleAuth() error {
	return nil
}

type WrapperError struct {
}

func updateScaffoldGenRule() *WrapperError {
	return nil
}

func TestName222(t *testing.T) {
	err := checkScaffoldRuleAuth()
	if err != nil {
		t.Fatal(err)
	}
	if err := checkScaffoldRuleAuth(); err != nil {
		t.Fatal(err)
	}
	if err := updateScaffoldGenRule(); err != nil {
		t.Fatal(err)
	}
}

func TestRegexp(t *testing.T) {
	compile := regexp.MustCompile(`^.+$`)
	t.Log(compile.MatchString("..."))
	t.Log(compile.MatchString("...a"))
}

func TestUniqueString(t *testing.T) {
	assert.Equal(t, UniqueString([]string{"k3", "k1", "k2", "k3", "k4"}), []string{"k3", "k1", "k2", "k4"})
}

func TestYamlUnmarshal(t *testing.T) {
	type Data struct {
		K1 string `yaml:"K1"`
		K2 string `yaml:"K2"`
	}
	type Test struct {
		Data1 string `yaml:"Data1"`
		Data2 Data   `yaml:"Data2"`
	}
	test := Test{
		Data1: "1",
		Data2: Data{
			K1: "1111",
			K2: "2222",
		},
	}
	if err := yaml.Unmarshal([]byte(`
Data1: 222
Data2:
  K2: 3333`), &test); err != nil {
		t.Fatal(err)
	}
	t.Logf("%#v", test)
}
