package main

import (
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewConverterFromYAML(test *testing.T) {
	raw, err := ioutil.ReadFile("converter.yml")
	assert.NoError(test, err)

	converter, err := NewConverterFromYAML(raw)
	assert.NoError(test, err)

	assert.Equal(test, true, len(converter.Conversions) > 5)
}

func TestConverterConvert(test *testing.T) {

	raw, err := ioutil.ReadFile("converter.yml")
	assert.NoError(test, err)

	converter, err := NewConverterFromYAML(raw)
	assert.NoError(test, err)

	input := Quantity{Magnitude: 1000, Unit: "mm"}
	expectedOutput := Quantity{Magnitude: 39.3700787, Unit: "in"}
	output, err := converter.Convert(input, expectedOutput.Unit)
	assert.NoError(test, err)
	assert.Equal(test, expectedOutput, output)

	inputCached := Quantity{Magnitude: 100, Unit: "mm"}
	expectedOutputCached := Quantity{Magnitude: 3.9370078700000004, Unit: "in"}
	outputCached, err := converter.Convert(inputCached, expectedOutputCached.Unit)
	assert.NoError(test, err)
	assert.Equal(test, expectedOutputCached, outputCached)
}

func TestFailConversionConverterConvertWithMissingConversion(test *testing.T) {
	raw, err := ioutil.ReadFile("converter.yml")
	assert.NoError(test, err)

	converter, err := NewConverterFromYAML(raw)
	assert.NoError(test, err)

	input := Quantity{Magnitude: 1000, Unit: "mm"}
	expectedOutput := Quantity{}
	output, err := converter.Convert(input, "idonotexist")
	assert.Error(test, err)
	assert.Equal(test, expectedOutput, output)
}

func TestFailConversionConverterConvertWithBadConversion(test *testing.T) {
	raw, err := ioutil.ReadFile("converter.yml")
	assert.NoError(test, err)

	converter, err := NewConverterFromYAML(raw)
	assert.NoError(test, err)

	input := Quantity{Magnitude: 1000, Unit: "mm"}
	expectedOutput := Quantity{}
	output, err := converter.Convert(input, "idonotexist")
	assert.Error(test, err)
	assert.Equal(test, expectedOutput, output)
}

func TestConvertWithValidInput(test *testing.T) {
	input := Quantity{Magnitude: 1100, Unit: "g"}
	expectedOutput := Quantity{Magnitude: 1.1, Unit: "kg"}
	conversion := Conversion{From: "g", To: "kg", Formula: "magnitude / 1000"}
	output, err := conversion.Convert(input)
	assert.NoError(test, err)
	assert.Equal(test, expectedOutput, output)
}

func TestConversionTest(test *testing.T) {
	conversion := Conversion{
		From:    "g",
		To:      "kg",
		Formula: "magnitude / 1000",
		TestFixtures: []ConversionTestFixture{
			ConversionTestFixture{
				Input:    1000,
				Expected: 1,
			},
		},
	}
	err := conversion.Test()

	assert.NoError(test, err)
}

func TestFailConversionTestWithBadFormula(test *testing.T) {
	conversion := Conversion{
		From:    "g",
		To:      "kg",
		Formula: "magnitude / 1031",
		TestFixtures: []ConversionTestFixture{
			ConversionTestFixture{
				Input:    1000,
				Expected: 1,
			},
		},
	}
	err := conversion.Test()

	assert.Error(test, err)
}

func TestFailConversionConvertWithUnitMismatch(test *testing.T) {
	input := Quantity{Magnitude: 1100, Unit: "g"}
	expectedOutput := Quantity{}
	conversion := Conversion{From: "µg", To: "kg", Formula: "magnitude * 1000"}
	output, err := conversion.Convert(input)

	assert.Error(test, err)
	assert.Equal(test, expectedOutput, output)
}

func TestFailConversionConvertWithInvalidFormula(test *testing.T) {
	input := Quantity{Magnitude: 1100, Unit: "g"}
	expectedOutput := Quantity{}
	conversion := Conversion{From: "g", To: "kg", Formula: "magnitude /**e2> 1000"}
	output, err := conversion.Convert(input)

	assert.Error(test, err)
	assert.Equal(test, expectedOutput, output)
}

func TestFailConversionTestWithInvalidFormula(test *testing.T) {
	conversion := Conversion{
		From:    "g",
		To:      "kg",
		Formula: "magnitude /**e2> 1000",
		TestFixtures: []ConversionTestFixture{
			ConversionTestFixture{
				Input:    1000,
				Expected: 1,
			},
		},
	}
	err := conversion.Test()

	assert.Error(test, err)
}

func TestFailConversionTestWithMissingFormula(test *testing.T) {
	conversion := Conversion{
		From: "g",
		To:   "kg",
		TestFixtures: []ConversionTestFixture{
			ConversionTestFixture{
				Input:    1000,
				Expected: 1,
			},
		},
	}
	err := conversion.Test()

	assert.Error(test, err)
}

func TestFailConversionConvertWithInvalidMagnitudeReference(test *testing.T) {
	input := Quantity{Magnitude: 1100, Unit: "g"}
	expectedOutput := Quantity{}
	conversion := Conversion{From: "g", To: "kg", Formula: "badmagnitude * 1000"}
	output, err := conversion.Convert(input)

	assert.Error(test, err)
	assert.Equal(test, expectedOutput, output)
}

func TestFailsConverterConvertWithBadConverions(test *testing.T) {
	input := Quantity{Magnitude: 1, Unit: "m"}
	expectedOutput := Quantity{}
	converter := Converter{
		Conversions: []Conversion{
			Conversion{
				From:    "m",
				To:      "km",
				Formula: "magnitude *%¤ 1000",
				TestFixtures: []ConversionTestFixture{
					ConversionTestFixture{
						Input:    1,
						Expected: 1000,
					},
				},
			},
			Conversion{
				From:    "km",
				To:      "m",
				Formula: "magnitude / 1000",
				TestFixtures: []ConversionTestFixture{
					ConversionTestFixture{
						Input:    1000,
						Expected: 1,
					},
				},
			},
		},
	}
	output, err := converter.Convert(input, "km")

	assert.Error(test, err)
	assert.Equal(test, expectedOutput, output)
}

func TestConversionsFromYAML(test *testing.T) {
	input := `
conversions:
  - from: m
    to: km
    formula: magnitude * 1000
    testFixtures:
      - input: 1
        expected: 1000

  - from: km
    to: m
    formula: magnitude / 1000
    testFixtures:
      - input: 1000
        expected: 1
`

	expectedOutput := Converter{
		Conversions: []Conversion{
			Conversion{
				From:    "m",
				To:      "km",
				Formula: "magnitude * 1000",
				TestFixtures: []ConversionTestFixture{
					ConversionTestFixture{
						Input:    1,
						Expected: 1000,
					},
				},
			},
			Conversion{
				From:    "km",
				To:      "m",
				Formula: "magnitude / 1000",
				TestFixtures: []ConversionTestFixture{
					ConversionTestFixture{
						Input:    1000,
						Expected: 1,
					},
				},
			},
		},
	}
	output, err := NewConverterFromYAML([]byte(input))
	output.Conversions[0].FormulaExpression = nil
	output.Conversions[1].FormulaExpression = nil

	assert.NoError(test, err)
	assert.Equal(test, expectedOutput, output)
}

func TestFailConversionsFromYAMLWithBadFormula(test *testing.T) {
	input := `
conversions:
  - from: m
    to: km
    formula: magnitude "#%"* 1000
    testFixtures:
      - input: 1
        expected: 1000

  - from: km
    to: m
    formula: magnitude / 1000
    testFixtures:
      - input: 1000
        expected: 1
`

	expectedOutput := Converter{}
	output, err := NewConverterFromYAML([]byte(input))

	assert.Error(test, err)
	assert.Equal(test, expectedOutput, output)
}

func TestFailConversionsFromYAMLWithBadYaml(test *testing.T) {
	input := `
m:2#¤234
  - unit: km
    formula
    testFixtures:
      - input: 1
        expected: 10
  `
	expectedOutput := Converter{}
	output, err := NewConverterFromYAML([]byte(input))

	assert.Error(test, err)
	assert.Equal(test, expectedOutput, output)
}
