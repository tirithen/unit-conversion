package main

import (
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConvertWithValidInput(test *testing.T) {
	input := Quantity{Magnitude: 1100, Unit: "g"}
	expectedOutput := Quantity{Magnitude: 1.1, Unit: "kg"}
	conversion := Conversion{From: "g", To: "kg", Formula: "magnitude / 1000"}
	output, err := Convert(input, conversion)

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

func TestFailConvertWithUnitMismatch(test *testing.T) {
	input := Quantity{Magnitude: 1100, Unit: "g"}
	expectedOutput := Quantity{}
	conversion := Conversion{From: "µg", To: "kg", Formula: "magnitude * 1000"}
	output, err := Convert(input, conversion)

	assert.Error(test, err)
	assert.Equal(test, expectedOutput, output)
}

func TestFailConvertWithInvalidFormula(test *testing.T) {
	input := Quantity{Magnitude: 1100, Unit: "g"}
	expectedOutput := Quantity{}
	conversion := Conversion{From: "g", To: "kg", Formula: "magnitude /**e2> 1000"}
	output, err := Convert(input, conversion)

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

func TestFailConvertWithInvalidMagnitudeReference(test *testing.T) {
	input := Quantity{Magnitude: 1100, Unit: "g"}
	expectedOutput := Quantity{}
	conversion := Conversion{From: "g", To: "kg", Formula: "badmagnitude * 1000"}
	output, err := Convert(input, conversion)

	assert.Error(test, err)
	assert.Equal(test, expectedOutput, output)
}

func TestConversionsFromYAML(test *testing.T) {
	// TODO: sort the input to prevent random errors due to map sorting
	input := `
m:
  - unit: km
    formula: magnitude / 1000
    testFixtures:
      - input: 1000
        expected: 1

  - unit: dm
    bad: hello
    formula: magnitude * 10
    testFixtures:
      - input: 1
        expected: 10

l:
  - unit: dl
    formula: magnitude * 10
    testFixtures:
      - input: 1
        expected: 10

  - unit: cl
    formula: magnitude * 100
    testFixtures:
      - input: 1
        expected: 100
`

	expectedOutput := []Conversion{
		Conversion{
			From:    "l",
			To:      "dl",
			Formula: "magnitude * 10",
			TestFixtures: []ConversionTestFixture{
				ConversionTestFixture{
					Input:    1,
					Expected: 10,
				},
			},
		},
		Conversion{
			From:    "l",
			To:      "cl",
			Formula: "magnitude * 100",
			TestFixtures: []ConversionTestFixture{
				ConversionTestFixture{
					Input:    1,
					Expected: 100,
				},
			},
		},
		Conversion{
			From:    "m",
			To:      "km",
			Formula: "magnitude / 1000",
			TestFixtures: []ConversionTestFixture{
				ConversionTestFixture{
					Input:    1000,
					Expected: 1,
				},
			},
		},
		Conversion{
			From:    "m",
			To:      "dm",
			Formula: "magnitude * 10",
			TestFixtures: []ConversionTestFixture{
				ConversionTestFixture{
					Input:    1,
					Expected: 10,
				},
			},
		},
	}
	output, err := ConversionsFromYAML(input)

	assert.NoError(test, err)
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
	output, err := ConversionsFromYAML(input)

	assert.Error(test, err)
	assert.Empty(test, output)
}

func TestFailConversionsFromYAMLWithRootAsList(test *testing.T) {
	input := "- m\n- l"
	output, err := ConversionsFromYAML(input)

	assert.Error(test, err)
	assert.Empty(test, output)
}

func TestFailConversionsFromYAMLWithMissingPossibleConversions(test *testing.T) {
	input := `
m:
  - unit: km
    formula: magnitude / 1000
    testFixtures:
      - input: 1
        expected: 10

  - unit: dm
    formula: magnitude * 10
    testFixtures:
      - input: 1
        expected: 10

l:
  bad: 12
  `
	output, err := ConversionsFromYAML(input)

	assert.Error(test, err)
	assert.Empty(test, output)
}

func TestFailConversionsFromYAMLWithBadConversionFormat(test *testing.T) {
	input := "m:\n  - unit"
	output, err := ConversionsFromYAML(input)

	assert.Error(test, err)
	assert.Empty(test, output)
}

func TestFailConversionsFromYAMLWithMissingPossibleConversionUnit(test *testing.T) {
	input := `
m:
  - formula: magnitude / 1000
    testFixtures:
      - input: 1
        expected: 10

  - unit: dm
    formula: magnitude * 10
    testFixtures:
      - input: 1
        expected: 10
  `
	output, err := ConversionsFromYAML(input)

	assert.Error(test, err)
	assert.Empty(test, output)
}

func TestFailConversionsFromYAMLWithMissingPossibleConversionFormula(test *testing.T) {
	input := `
m:
  - unit: km
    formula: magnitude / 1000
    testFixtures:
      - input: 1
        expected: 10

l:
  - unit: dl
    testFixtures:
      - input: 1
        expected: 10
  `
	output, err := ConversionsFromYAML(input)

	assert.Error(test, err)
	assert.Empty(test, output)
}

func TestFailConversionsFromYAMLWithMissingTestFixtures(test *testing.T) {
	input := `
m:
  - unit: km
    formula: magnitude / 1000

l:
  - unit: dl
    formula: magnitude * 10
    testFixtures:
      - input: 1
        expected: 10
  `
	output, err := ConversionsFromYAML(input)

	assert.Error(test, err)
	assert.Empty(test, output)
}

func TestFailConversionsFromYAMLWithEmptyTestFixtures(test *testing.T) {
	input := `
m:
  - unit: km
    formula: magnitude / 1000
    testFixtures:
      - input: 1
        expected: 1000

l:
  - unit: dl
    formula: magnitude * 10
    testFixtures:
      - input: 1
        expected: 10
  `
	output, err := ConversionsFromYAML(input)

	assert.Error(test, err)
	assert.Empty(test, output)
}

func TestConvertToFirstOption(test *testing.T) {
	conversions := []Conversion{
		Conversion{From: "l", To: "dl", Formula: "magnitude * 10"},
		Conversion{From: "kg", To: "g", Formula: "magnitude * 1000"},
		Conversion{From: "l", To: "cl", Formula: "magnitude * 100"},
	}
	input := Quantity{Magnitude: 3.54, Unit: "kg"}
	expectedOutput := Quantity{Magnitude: 3540, Unit: "g"}
	output, err := ConvertToFirstOption(input, conversions)

	assert.NoError(test, err)
	assert.Equal(test, expectedOutput, output)
}

func TestFailConvertToFirstOptionMissingConversion(test *testing.T) {
	conversions := []Conversion{
		Conversion{From: "l", To: "dl", Formula: "magnitude * 10"},
		Conversion{From: "kg", To: "g", Formula: "magnitude * 1000"},
		Conversion{From: "l", To: "cl", Formula: "magnitude * 100"},
	}
	input := Quantity{Magnitude: 4.35, Unit: "km"}
	expectedOutput := Quantity{}
	output, err := ConvertToFirstOption(input, conversions)

	assert.Error(test, err)
	assert.Equal(test, expectedOutput, output)
}

func TestConversionsDefinitionYAML(test *testing.T) {
	raw, err := ioutil.ReadFile("conversions.yml")
	assert.NoError(test, err)

	conversions, err := ConversionsFromYAML(string(raw))
	assert.NoError(test, err)

	assert.True(test, len(conversions) > 10)
}
