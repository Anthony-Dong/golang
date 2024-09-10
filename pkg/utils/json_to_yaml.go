package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/iancoleman/orderedmap"
	"github.com/tidwall/gjson"
	"gopkg.in/yaml.v3"
)

func JsonToYaml(input []byte) (output []byte, err error) {
	yamlNode := &yaml.Node{}
	if err := toYAMLNode(yamlNode, gjson.ParseBytes(input)); err != nil {
		return nil, err
	}
	var buf bytes.Buffer
	encoder := yaml.NewEncoder(&buf)
	encoder.SetIndent(2)
	if err := encoder.Encode(yamlNode); err != nil {
		return nil, fmt.Errorf("error encoding YAML: %v", err)
	}
	return buf.Bytes(), nil
}

func toYAMLNode(node *yaml.Node, data gjson.Result) error {
	switch data.Type {
	case gjson.JSON:
		var err error
		if data.IsObject() {
			node.Kind = yaml.MappingNode
			data.ForEach(func(key, value gjson.Result) bool {
				keyNode := &yaml.Node{}
				if err = toYAMLNode(keyNode, key); err != nil {
					return false
				}
				keyNode.Tag = ""
				valNode := &yaml.Node{}
				if err = toYAMLNode(valNode, value); err != nil {
					return false
				}
				node.Content = append(node.Content, keyNode, valNode)
				return true
			})
			if err != nil {
				return err
			}
		} else if data.IsArray() {
			node.Kind = yaml.SequenceNode
			// todo style
			data.ForEach(func(_, value gjson.Result) bool {
				valNode := &yaml.Node{}
				if err = toYAMLNode(valNode, value); err != nil {
					return false
				}
				node.Content = append(node.Content, valNode)
				return true
			})
			if err != nil {
				return err
			}
		} else {
			return fmt.Errorf("invalid JSON node type: %v", data.Type)
		}
	case gjson.Number:
		node.Kind = yaml.ScalarNode
		node.Value = data.Raw
	case gjson.String:
		node.Kind = yaml.ScalarNode
		node.Tag = "!!str"
		node.Value = data.Str
	case gjson.False, gjson.True:
		node.Kind = yaml.ScalarNode
		node.Tag = "!!bool"
		node.Value = data.String()
	case gjson.Null:
		node.Kind = yaml.ScalarNode
		node.Tag = "!!null"
		node.Value = "null"
	default:
		return fmt.Errorf("invalid JSON node type: %v", data.Type)
	}
	return nil
}

func YamlToJson(input []byte) ([]byte, error) {
	yamlNode := &yaml.Node{}
	if err := yaml.Unmarshal(input, yamlNode); err != nil {
		return nil, err
	}
	node, err := toJsonNode(yamlNode)
	if err != nil {
		return nil, err
	}
	output := bytes.NewBuffer(nil)
	encoder := json.NewEncoder(output)
	encoder.SetEscapeHTML(false)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(node); err != nil {
		return nil, err
	}
	return output.Bytes(), nil
}

func toJsonNode(node *yaml.Node) (interface{}, error) {
	toString := func(input interface{}) string {
		switch v := input.(type) {
		case string:
			return v
		case json.Number:
			return string(v)
		case bool:
			if v {
				return "true"
			}
			return "false"
		case nil:
			return "null"
		default:
			return fmt.Sprintf(`%v`, input)
		}
	}
	switch node.Kind {
	case yaml.MappingNode:
		result := orderedmap.New()
		result.SetEscapeHTML(false)
		for index, _ := range node.Content {
			if index%2 == 1 {
				key, err := toJsonNode(node.Content[index-1])
				if err != nil {
					return nil, err
				}
				value, err := toJsonNode(node.Content[index])
				if err != nil {
					return nil, err
				}
				result.Set(toString(key), value)
			}
		}
		return result, nil
	case yaml.SequenceNode:
		result := make([]interface{}, 0, len(node.Content))
		for index, _ := range node.Content {
			value, err := toJsonNode(node.Content[index])
			if err != nil {
				return nil, err
			}
			result = append(result, value)
		}
		return result, nil
	case yaml.ScalarNode:
		switch node.Tag {
		case "!!str":
			return node.Value, nil
		case "!!float", "!!int":
			return json.Number(node.Value), nil
		case "!!bool":
			return strconv.ParseBool(node.Value)
		case "!!null":
			return nil, nil
		}
		return node.Value, nil
	case yaml.DocumentNode:
		if len(node.Content) == 1 {
			return toJsonNode(node.Content[0])
		}
		clone := *node
		clone.Kind = yaml.SequenceNode
		return toJsonNode(&clone)
	default:
		return nil, fmt.Errorf("invalid JSON node type: %v", node.Kind)
	}
}
