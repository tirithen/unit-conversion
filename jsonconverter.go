package main

import (
	"encoding/json"
	"strconv"

	yaml "gopkg.in/yaml.v2"
)

type mapNode map[string]json.RawMessage
type arrayNode []json.RawMessage

type JSONConverter struct {
	Converter
}

func (converter *JSONConverter) walkJSON(path string, rawNode json.RawMessage, input string) (output string, errors []error) {
	output = input

	if rawNode[0] == 123 { // 123 is `{` => object
		var node mapNode
		json.Unmarshal(rawNode, &node)

		quantity := Quantity{}
		hasMagnitude := false
		hasUnit := false

		for property, value := range node {
			newPath := ""
			if path == "" {
				newPath = property
			} else {
				newPath = path + "." + property
			}

			subErrors := []error{}
			output, subErrors = converter.walkJSON(newPath, value, output)
			errors = append(errors, subErrors...)

			if property == "magnitude" {
				var valueInterface interface{}
				json.Unmarshal(value, &valueInterface)
				switch valueTyped := valueInterface.(type) {
				case float64:
					quantity.Magnitude = valueTyped
					hasMagnitude = true
				}
			} else if property == "unit" {
				var valueInterface interface{}
				json.Unmarshal(value, &valueInterface)
				switch valueTyped := valueInterface.(type) {
				case string:
					quantity.Unit = valueTyped
					hasUnit = true
				}
			}
		}

		if hasMagnitude && hasUnit {
		}
	} else if rawNode[0] == 91 { // 91 is `[` => array
		var node arrayNode
		json.Unmarshal(rawNode, &node)
		for index, value := range node {
			newPath := ""
			if path == "" {
				newPath = strconv.Itoa(index)
			} else {
				newPath = path + "." + strconv.Itoa(index)
			}

			subErrors := []error{}
			output, subErrors = converter.walkJSON(newPath, value, output)
			errors = append(errors, subErrors...)
		}
	}

	return output, errors
}

// ConvertToPreferredUnits will search through JSON and convert any magnitude/unit pair that it can find
func (converter *JSONConverter) ConvertToPreferredUnits(input string) (output string, errors []error) {
	var node json.RawMessage
	err := json.Unmarshal([]byte(input), &node)
	if err != nil {
		errors = append(errors, err)
		return
	}

	output, errors = converter.walkJSON("", node, input)

	return
}

// NewJSONConverterFromYAML is used to parse and verify YAML data into a Converter
func NewJSONConverterFromYAML(raw []byte) (converter JSONConverter, err error) {
	err = yaml.Unmarshal(raw, &converter)
	if err != nil {
		converter = JSONConverter{}
		return
	}

	err = converter.Test()
	if err != nil {
		converter = JSONConverter{}
		return
	}

	return
}
