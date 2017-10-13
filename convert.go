package main

import (
	"errors"
	"fmt"
	"io"
	"sort"

	"github.com/kylelemons/go-gypsy/yaml"
	govaluate "gopkg.in/Knetic/govaluate.v2"
)

// Quantity defines properties that is needed to make a conversion
type Quantity struct {
	Magnitude float64 `json:"magnitude"`
	Unit      string  `json:"unit"`
}

// Conversion defines properties that describes how a value with one unit can be converted into a value in another unit with a formula
type Conversion struct {
	From    string
	To      string
	Formula string
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
func ConversionsFromYAML(reader io.Reader) (conversions []Conversion, err error) {
	rootNode, err := yaml.Parse(reader)
	if err != nil {
		err = errors.New("Conversions is not valid YAML")
		return
	}

	root, ok := rootNode.(yaml.Map)
	if !ok {
		err = errors.New("Conversions root node has to be a map")
		return
	}

	var defaultMeasurements []string
	for key := range root {
		defaultMeasurements = append(defaultMeasurements, key)
	}
	sort.Strings(defaultMeasurements)

	for key := range defaultMeasurements {
		possibleConversions, ok := root[defaultMeasurements[key]].(yaml.List)
		if !ok {
			err = errors.New("Each default unit needs to have a list of available conversions")
			conversions = []Conversion{}
			return
		}

		for index := range possibleConversions {
			possibleConversion, ok := possibleConversions[index].(yaml.Map)
			if !ok {
				err = errors.New("Each possible conversion needs to be a map")
				conversions = []Conversion{}
				return
			}

			unit, ok := possibleConversion.Key("unit").(yaml.Scalar)
			if !ok {
				err = errors.New("Each possible conversion needs to be a map with a unit property")
				conversions = []Conversion{}
				return
			}

			formula, ok := possibleConversion.Key("formula").(yaml.Scalar)
			if !ok {
				err = errors.New("Each possible conversion needs to be a map with a formula property")
				conversions = []Conversion{}
				return
			}

			conversion := Conversion{From: defaultMeasurements[key], To: unit.String(), Formula: formula.String()}

			conversions = append(conversions, conversion)
		}
	}

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
