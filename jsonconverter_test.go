package main

import (
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestJSONConverterConvertToPreferredUnits(test *testing.T) {
	input, err := ioutil.ReadFile("fixtures/input.json")
	assert.NoError(test, err)
	expectedOutput, err := ioutil.ReadFile("fixtures/output.json")
	assert.NoError(test, err)
	converterConfig, err := ioutil.ReadFile("converter.yml")
	assert.NoError(test, err)
	converter, err := NewJSONConverterFromYAML(converterConfig)
	assert.NoError(test, err)

	output, errors := converter.ConvertToPreferredUnits(string(input))
	assert.Empty(test, errors)
	assert.JSONEq(test, string(expectedOutput), output)
}

func TestFailJSONConverterConvertToPreferredUnitsWithBadSyntax(test *testing.T) {
	input := `{ "test": { "magnitude": 12  "unit": "cm" } }`
	expectedOutput := ""
	converterConfig, err := ioutil.ReadFile("converter.yml")
	assert.NoError(test, err)
	converter, err := NewJSONConverterFromYAML(converterConfig)
	assert.NoError(test, err)

	output, errors := converter.ConvertToPreferredUnits(input)
	assert.NotEmpty(test, errors)
	assert.Equal(test, expectedOutput, string(output))
}
