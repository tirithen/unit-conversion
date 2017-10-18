[![Build Status](https://travis-ci.org/tirithen/unit-conversion.svg?branch=master)](https://travis-ci.org/tirithen/unit-conversion)

[![Coverage Status](https://coveralls.io/repos/github/tirithen/unit-conversion/badge.svg?branch=master)](https://coveralls.io/github/tirithen/unit-conversion?branch=master)

# unit-conversion

An effort to solve the common problem of unit conversions. There are lots of language specific libraries that allows to convert in between some specific units but I had a hard time finding one service that would let me configure for specific units that is less common (e.g. µg/l to ng/l or similar).

A configurable HTTP service and a Go package that allows for conversion between units. The project is meant to be run as a docker container and communicated with through HTTP.

## Philosophy

* Ease of use, starting the service and sending a HTTP POST request with any JSON data should be enough to use the service.
* Sane defaults, the service should always start without any specific configurations, where needed good default values should be set.
* This service main goal is to provide a Docker packaged HTTP service, the Go package is nice and should have good interfaces but is second in priority.

## Docker image

The docker image is available here https://hub.docker.com/r/tirithen/unit-conversion/.

## Try it out (HTTP version)

**Note:** HTTP service has yet to be fully completed

Send any JSON structure in a POST request, any object that has both magnitude and unit will have those values converted and updated in the response.

    $ docker run -p 8080:8080 --name unit-conversion -d tirithen/unit-conversion
    $ curl -d '{"model":"Supertablet","size":{"magnitude":10, "unit":"in"}}' -H "Content-Type: application/json" -X POST http://localhost:8080/

    Returns {"model":"Supertablet","size":{"magnitude":25.4, "unit":"cm"}}

You can send one or several of these objects per call and all of them will get converted. To control which units that will be converted to you need to set that in a configuration YAML file (have a look under the heading Configuration).

Any unit that lacks a configuration will just be ignored.

## Go example (Go version)

Basic usage is the following:

    package main

    import (
      "fmt"
      "io/ioutil"

      unitconversion "github.com/tirithen/unit-conversion"
    )

    func main() {
      configuration, err := ioutil.ReadFile("./converter.yml")
      converter := unitconversion.NewConverterFromYAML(configuration)

      input := unitconversion.Quantity{Magnitude: 10, Unit: "in"}
      output, err := converter.Convert(input, "cm")
      if err != nil {
        panic(err)
      }

      fmt.Print(output.Magnitude + " " + output.Unit) // Prints: 25.4 in
    }

## Configuration

The service can be configured with a YAML file according to the following:

    preferredUnits:
      - m
      - µg/l

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

      - from: m
        to: in
        formula: magnitude * 39.37007874
        testFixtures:
          - input: 1
            expected: 39.37007874

      - from: in
        to: m
        formula: magnitude * 0.0254
        testFixtures:
          - input: 1
            expected: 0.0254

      - from: µg/l
        to: ng/ml
        formula: magnitude
        testFixtures:
          - input: 1
            expected: 1

      - from: ng/ml
        to: µg/l
        formula: magnitude
        testFixtures:
          - input: 1
            expected: 1

### preferredUnits:

A list of the units that the service should convert to unless any other unit is specified, this is used entirely for the HTTP service

### conversions:

A list of the conversions that the service can handle.

Each conversion needs *from*, *to*, *formula* and in *testFixtures* at least one test fixture with a set of *input* and *expected* to verify that the formula calculates as intended.

Each time the service starts (or in Go when NewConverterFromYAML is called) all conversions are tested with their testFixtures to ensure that their formulas are correct.

### Conversions are chained automatically

If you want to convert in between cm and in and there are no no direct conversion defined but there is a conversion from cm to m and from m to in the service will automatically find that path and convert the amount of times that is needed to reach the final unit.

This way it's enough to add one conversion in each direction against one of the units in preferredUnits to "hook" that unit into the chain.
