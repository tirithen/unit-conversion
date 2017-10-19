package main // import "github.com/tirithen/unit-conversion"
import (
	"io/ioutil"
	"net/http"

	"github.com/labstack/echo"
)

var converter JSONConverter

func main() {
	server := echo.New()

	converterConfig, err := ioutil.ReadFile("./converter.yml")
	if err != nil {
		server.Logger.Fatal(err)
	}

	converter, err = NewJSONConverterFromYAML(converterConfig)
	if err != nil {
		server.Logger.Fatal(err)
	}

	server.GET("/", getHandler)
	server.POST("/", postHandler)

	server.Logger.Fatal(server.Start(":8080"))
}

func postHandler(context echo.Context) error {
	contentType := context.Request().Header.Get("Content-Type")

	if contentType == "application/json" {
		body, err := ioutil.ReadAll(context.Request().Body)
		if err != nil {
			context.Logger().Error(err)
			return context.String(http.StatusBadRequest, "Bad Request")
		}

		output, errors := converter.ConvertToPreferredUnits(string(body))
		if len(errors) > 0 {
			context.Logger().Debug(errors)
			return context.String(http.StatusBadRequest, "Bad Request")
		}

		return context.JSONBlob(http.StatusOK, []byte(output))
	}

	return context.String(
		http.StatusUnsupportedMediaType,
		"There are currently no support for Content-Type: "+contentType+" , currently application/json is supported.",
	)
}

func getHandler(context echo.Context) error {
	return context.String(
		http.StatusMethodNotAllowed,
		`
Make a POST request to this URL with Content-Type: application/json and a body that contains objects with the properties magnitude and unit. e.g.:

    $ curl -d '{"model":"Supertablet","size":{"magnitude":10, "unit":"in"}}' -H "Content-Type: application/json" -X POST http://localhost:8080/

Or use your favorite HTTP tool.

Some more examples of data structures

    {
      "demographic": {
        "name": "Lei"
      },
      "measurements": {
        "height": {
          "magnitude": 5,
          "unit": "ft",
          "extraproperty": "will be kept without modification"
        },
        "weight": [
          {
            "magnitude": 154.323584,
            "unit": "lb"
          },
          {
            "magnitude": 71.6,
            "unit": "kg"
          }
        ]
      }
    }

Or just lots of values in an array/object or similar. e.g.:

    [
      {
        "magnitude": 5,
        "unit": "ft"
      },
      {
        "magnitude": 5,
        "unit": "cm"
      },
      {
        "magnitude": 5,
        "unit": "ml"
      },
      {
        "magnitude": 5,
        "unit": "oz"
      }
    ]
    `,
	)
}
