package main

import (
	"fmt"

	govaluate "gopkg.in/Knetic/govaluate.v2"
	validator "gopkg.in/go-playground/validator.v9"
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
	From         string                  `yaml:"from" validate:"required"`
	To           string                  `yaml:"to" validate:"required"`
	Formula      string                  `yaml:"formula" validate:"required"`
	TestFixtures []ConversionTestFixture `yaml:"testFixtures" validate:"required,dive,required"`
}

// Test runs the conversion with all defined test fixtures to verify that the conversion returnes the values expected
func (conversion *Conversion) Test() (err error) {
	validate := validator.New()
	err = validate.Struct(conversion)
	if err != nil {
		return
	}

	for index := range conversion.TestFixtures {
		fixture := conversion.TestFixtures[index]
		input := Quantity{Magnitude: fixture.Input, Unit: conversion.From}
		expected := Quantity{Magnitude: fixture.Expected, Unit: conversion.To}
		output, conversionError := conversion.Convert(input)
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
func (conversion *Conversion) Convert(input Quantity) (output Quantity, err error) {
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

type Converter struct {
	Conversions []Conversion
	PathCache   map[string][]*Conversion
}

func (converter *Converter) Test() (err error) {
	for index := range converter.Conversions {
		err = converter.Conversions[index].Test()
		if err != nil {
			return
		}
	}

	return
}

func (converter *Converter) filterConversionsByFrom(from string, blacklist []*Conversion) (conversions []*Conversion) {
	for conversionsIndex := range converter.Conversions {
		conversion := &converter.Conversions[conversionsIndex]
		if conversion.From == from {
			blacklisted := false
			for blacklistIndex := range blacklist {
				if conversion == blacklist[blacklistIndex] {
					blacklisted = true
				}
			}

			if !blacklisted {
				conversions = append(conversions, conversion)
			}
		}
	}

	return
}

func (converter *Converter) getPath(from string, to string, previousPath []*Conversion) (path []*Conversion, err error) {
	cacheKey := from + " => " + to
	if cachedPath, ok := converter.PathCache[cacheKey]; ok {
		path = cachedPath
		return
	}

	edge := converter.filterConversionsByFrom(from, previousPath)

	for index := range edge {
		node := edge[index]
		if node.To == to {
			path = append(previousPath, node)
			return
		}
	}

	for index := range edge {
		node := edge[index]
		path, err = converter.getPath(node.To, to, append(previousPath, node))
		if err == nil && len(path) > 0 {
			if converter.PathCache == nil {
				converter.PathCache = make(map[string][]*Conversion)
			}
			converter.PathCache[cacheKey] = path
			return
		}
	}

	err = fmt.Errorf("Unable to find a path")
	path = []*Conversion{}

	return
}

func (converter *Converter) Convert(input Quantity, to string) (output Quantity, err error) {
	path, err := converter.getPath(input.Unit, to, []*Conversion{})
	if err != nil {
		return
	}

	output = input
	for index := range path {
		output, err = path[index].Convert(output)
		if err != nil {
			output = Quantity{}
			return
		}
	}

	return
}

func NewConverterFromYAML(raw []byte) (converter Converter, err error) {
	err = yaml.Unmarshal(raw, &converter.Conversions)
	if err != nil {
		converter = Converter{}
		return
	}

	err = converter.Test()
	if err != nil {
		converter = Converter{}
		return
	}

	return
}
