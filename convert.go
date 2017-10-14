package main

import (
	"fmt"

	govaluate "gopkg.in/Knetic/govaluate.v2"
	yaml "gopkg.in/yaml.v2"
)

// Quantity defines properties that is needed to make a conversion
type Quantity struct {
	Magnitude float64 `json:"magnitude"`
	Unit      string  `json:"unit"`
}

// ConversionTestFixture holds a test case that can be used to validate a Conversion.Formula
type ConversionTestFixture struct {
	Input    float64 `yaml:"input"`
	Expected float64 `yaml:"expected"`
}

// Conversion defines properties that describes how a value with one unit can be converted into a value in another unit with a formula
type Conversion struct {
	From         string
	To           string
	Formula      string
	TestFixtures []ConversionTestFixture
}

// TestFormula runs the conversion with all defined test fixtures to verify that the conversion returnes the values expected
func (conversion Conversion) TestFormula() (err error) {
	if len(conversion.TestFixtures) < 1 {
		err = fmt.Errorf("Conversion from %q to %q with formula %q does not have any testFixtures", conversion.From, conversion.To, conversion.Formula)
		return
	}

	for index := range conversion.TestFixtures {
		fixture := conversion.TestFixtures[index]
		input := Quantity{Magnitude: fixture.Input, Unit: conversion.From}
		expected := Quantity{Magnitude: fixture.Expected, Unit: conversion.To}
		output, conversionError := Convert(input, conversion)
		if conversionError != nil {
			err = conversionError
			return
		}

		if output.Magnitude != fixture.Expected {
			err = fmt.Errorf("Conversion test failed, from %q to %q with formula %q and input %f expected %f but got %f", conversion.From, conversion.To, conversion.Formula, input.Magnitude, expected.Magnitude, output.Magnitude)
			return
		}
	}

	return
}

// Convert takes a Quantity and a Conversion and returns a new Quantity with the result
func Convert(input Quantity, conversion Conversion) (output Quantity, err error) {
	if input.Unit != conversion.From {
		err = fmt.Errorf("Conversion from unit mismatch got %q but expected %q", input.Unit, conversion.From)
		return
	}

	expression, err := govaluate.NewEvaluableExpression(conversion.Formula)
	if err != nil {
		return
	}

	parameters := make(map[string]interface{}, 8)
	parameters["magnitude"] = input.Magnitude

	magnitude, err := expression.Evaluate(parameters)
	if err != nil {
		return
	}

	output.Magnitude = magnitude.(float64)
	output.Unit = conversion.To

	return
}

// ConversionsFromYAML reads a YAML stream and returns a slice of conversions
func ConversionsFromYAML(raw string) (conversions []Conversion, err error) {
	type conversionGroup struct {
		Unit         string                  `yaml:"unit"`
		Formula      string                  `yaml:"formula"`
		TestFixtures []ConversionTestFixture `yaml:"testFixtures"`
	}
	conversionGroups := make(map[string]conversionGroup)

	fmt.Println("conversionGroups", conversionGroups)
	raw = "m:\n  - unit: km\n    formula: magnitude / 1000\n    testFixture:\n      - input: 1000\n        expected: 1"
	fmt.Println(raw)

	err = yaml.Unmarshal([]byte(raw), &conversionGroups)
	if err != nil {
		return
	}

	fmt.Println("conversionGroups", conversionGroups)

	/*
		for key := range defaultMeasurements {


				conversion := Conversion{From: defaultMeasurements[key], To: unit.String(), Formula: formula.String()}

				err = conversion.TestFormula()
				if err != nil {
					conversions = []Conversion{}
					return
				}

				conversions = append(conversions, conversion)
			}
		}
	*/
	return
}

// ConvertToFirstOption returns a conversion to the first available conversion
func ConvertToFirstOption(input Quantity, conversions []Conversion) (output Quantity, err error) {
	for index := range conversions {
		if conversions[index].From == input.Unit {
			output, err = Convert(input, conversions[index])
			return
		}
	}

	err = fmt.Errorf("Unable to find a converter that can convert from %q", input.Unit)
	return
}
