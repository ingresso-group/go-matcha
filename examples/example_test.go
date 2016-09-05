package examples

import (
	"testing"

	"github.com/ingresso-group/go-matcha/matcha"
	. "github.com/smartystreets/goconvey/convey"
)

type expectedResponseFormat struct {
	Query struct {
		Count    float64
		Created  string
		Language string `json:"lang"` // Can explicitly define the name of the field we're expecting
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
			response, err := GetWeatherData()

			Convey("It should not return an error", func() {
				So(err, ShouldBeNil)
			})
			Convey("It should have same format", func() {
				So(response, matcha.ShouldMatchExpectedResponse, expected)
			})
		})
	})
}
