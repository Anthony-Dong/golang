package utils

import (
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
)

func Test_JsonToYaml_Obj(t *testing.T) {
	jsonData := `
{
    "str_data": "hello\nworld",
	"escape": "echo ${FIRST_KEY} && echo ${ALIAS_ENV}",
    "int_data": 132131231231313122312312312321111,
    "arr_data": [
        "Go",
        "Python",
        "JavaScript",
        132131231231313122312312312321111,
        null,
        {
            "k1": "v1",
            "k2": [
                1,
                2,
                3,
                -1,
                -1.1
            ]
        }
    ],
    "bool_data": true,
    "null_data": null,
	"11": 11,
	"1.11": 1.11,
	"1.111": "1.111",
	"1.1111": "11",
	"null": null,
	"false": false,
	"true": "true"
}
`
	yamlData, err := JsonToYaml([]byte(jsonData))
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(yamlData))
	output, err := YamlToJson(yamlData)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(output))
	assert.Equal(t, rmSpace(jsonData), rmSpace(string(output)))
}

func rmSpace(input string) string {
	return regexp.MustCompile(`\s+`).ReplaceAllString(input, "")
}

func Test_JsonToYaml_Array(t *testing.T) {
	jsonData := `["Go", "Python", "JavaScript"]`
	yamlData, err := JsonToYaml([]byte(jsonData))
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(yamlData))
	output, err := YamlToJson(yamlData)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(output))
	assert.Equal(t, rmSpace(jsonData), rmSpace(string(output)))
}

func TestYamlNode(t *testing.T) {
	yamlNode := &yaml.Node{
		Kind: yaml.MappingNode,
		Content: []*yaml.Node{
			{Kind: yaml.ScalarNode, Value: "name"},
			{Kind: yaml.ScalarNode, Value: "John"},
			{Kind: yaml.ScalarNode, Value: "age"},
			{Kind: yaml.ScalarNode, Value: "30"},
			{Kind: yaml.ScalarNode, Value: "languages"},
			{Kind: yaml.SequenceNode, Content: []*yaml.Node{
				{Kind: yaml.ScalarNode, Value: "Go"},
				{Kind: yaml.ScalarNode, Value: "Python"},
				{Kind: yaml.ScalarNode, Value: "JavaScript"},
			}, Style: yaml.FlowStyle},
		},
	}
	output, err := yaml.Marshal(yamlNode)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(output))
}
func Test_JsonToYaml_EmptyObject(t *testing.T) {
	jsonData := `{}`

	yamlData, err := JsonToYaml([]byte(jsonData))
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(yamlData))

	output, err := YamlToJson(yamlData)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(output))

	assert.Equal(t, rmSpace(jsonData), rmSpace(string(output)))
}

func Test_JsonToYaml_EmptyArray(t *testing.T) {
	jsonData := `[]`

	yamlData, err := JsonToYaml([]byte(jsonData))
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(yamlData))

	output, err := YamlToJson(yamlData)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(output))

	assert.Equal(t, rmSpace(jsonData), rmSpace(string(output)))
}

func Test_JsonToYaml_NestedObject(t *testing.T) {
	jsonData := `
{
	"level1": {
		"level2": {
			"level3": {
				"key": "value"
			}
		}
	}
}`

	yamlData, err := JsonToYaml([]byte(jsonData))
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(yamlData))

	output, err := YamlToJson(yamlData)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(output))

	assert.Equal(t, rmSpace(jsonData), rmSpace(string(output)))
}

func Test_JsonToYaml_SpecialCharacters(t *testing.T) {
	jsonData := `
{
	"special": "!@#$%^&*()_+-=[]{}|;':\",./<>?"
}`

	yamlData, err := JsonToYaml([]byte(jsonData))
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(yamlData))

	output, err := YamlToJson(yamlData)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(output))

	assert.Equal(t, rmSpace(jsonData), rmSpace(string(output)))
}

func Test_JsonToYaml_BooleanValues(t *testing.T) {
	jsonData := `
{
	"trueValue": true,
	"falseValue": false
}`

	yamlData, err := JsonToYaml([]byte(jsonData))
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(yamlData))

	output, err := YamlToJson(yamlData)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(output))

	assert.Equal(t, rmSpace(jsonData), rmSpace(string(output)))
}
