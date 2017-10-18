package main // import "github.com/tirithen/unit-conversion"

import (
	"fmt"
	"io/ioutil"
)

func main() {
	configuration, err := ioutil.ReadFile("./converter.yml")
	converter := unitconversion.NewConverterFromYAML(configuration)

	input := unitconversion.Quantity{Magnitude: 10, Unit: "in"}
	output, err := converter.Convert(input, "cm")
	if err != nil {
		panic(err)
	}

	fmt.Print(output)
}
