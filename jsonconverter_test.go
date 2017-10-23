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

func TestJSONConverterConvertToPreferredUnitsWithLargeDataSet(test *testing.T) {
	input, err := ioutil.ReadFile("fixtures/inputLarge.json")
	assert.NoError(test, err)
	expectedOutput, err := ioutil.ReadFile("fixtures/outputLarge.json")
	assert.NoError(test, err)
	converterConfig, err := ioutil.ReadFile("converter.yml")
	assert.NoError(test, err)
	converter, err := NewJSONConverterFromYAML(converterConfig)
	assert.NoError(test, err)

	output, errors := converter.ConvertToPreferredUnits(string(input))
	assert.Empty(test, errors)
	assert.JSONEq(test, string(expectedOutput), output)
}

func TestNewJSONConverterFromYAML(test *testing.T) {
	badConfig := "broken yaml¤-:4"
	expectedOutput := JSONConverter{}
	converter, err := NewJSONConverterFromYAML([]byte(badConfig))

	assert.Error(test, err)
	assert.Equal(test, expectedOutput, converter)
}

func BenchmarkNewJSONConverterFromYAMLLargeDataSet(benchmark *testing.B) {
	converterConfig, err := ioutil.ReadFile("converter.yml")
	if err != nil {
		panic(err)
	}

	for index := 0; index < benchmark.N; index++ {
		_, err = NewJSONConverterFromYAML(converterConfig)
		if err != nil {
			panic(err)
		}
	}
}

func BenchmarkJSONConverterConvertToPreferredUnitsLargeDataSet(benchmark *testing.B) {
	input, err := ioutil.ReadFile("fixtures/inputLarge.json")
	if err != nil {
		panic(err)
	}

	converterConfig, err := ioutil.ReadFile("converter.yml")
	if err != nil {
		panic(err)
	}

	for index := 0; index < benchmark.N; index++ {
		converter, err := NewJSONConverterFromYAML(converterConfig)
		if err != nil {
			panic(err)
		}

		converter.ConvertToPreferredUnits(string(input))
	}
}

func BenchmarkJSONConverterConvertToPreferredUnit(benchmark *testing.B) {
	converter := Converter{
		PreferredUnits: []string{"in"},
		Conversions: []Conversion{
			Conversion{From: "cm", To: "in", Formula: "magnitude * 2.54"},
		},
	}

	for index := 0; index < benchmark.N; index++ {
		converter.ConvertToPreferredUnit(Quantity{Magnitude: 1, Unit: "cm"})
	}
}
