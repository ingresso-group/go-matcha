package examples

import (
	"testing"

	"github.com/ingresso-group/go-matcha/matcha"
	. "github.com/smartystreets/goconvey/convey"
)

type expectedResponseFormat struct {
	Query struct {
		Count    float64 `capture:"count"`
		Created  string  `pattern:"^[0-9]{4}-[0-9]{2}-[0-9]{2}T[0-9]{2}:[0-9]{2}:[0-9]{2}Z"` // pattern matching
		Language string  `json:"lang"`                                                       // Can explicitly define the name of the field we're expecting
		Results  struct {
			Channel struct {
				Item struct {
					Condition struct {
						Code string
						Date string
						Temp string
						Text string
					}
				}
			}
		}
	}
}

func TestGetWeatherData(t *testing.T) {
	Convey("Given an expected response format", t, func() {
		var expected expectedResponseFormat

		Convey("When fetching weather data", func() {
			response := GetWeatherData()

			Convey("It should have same format", func() {
				So(response, matcha.ShouldMatchExpectedJSONResponse, expected, nil)
			})

			Convey("Count should be greater than zero", func() {
				capturedValues := make(map[string]interface{})
				So(response, matcha.ShouldMatchExpectedJSONResponse, expected, capturedValues)
				count := capturedValues["count"]
				So(count, ShouldBeGreaterThan, 0)
			})
		})
	})
}
