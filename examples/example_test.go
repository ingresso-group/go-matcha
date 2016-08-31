package examples

import (
	"testing"

	"github.com/ingresso-group/go-matcha/json"
	. "github.com/smartystreets/goconvey/convey"
)

type expectedResponseFormat struct {
	Query struct {
		Count   float64 `json:"count"`
		Created string  `json:"created"`
		Lang    string  `json:"lang"`
		Results struct {
			Channel struct {
				Item struct {
					Condition struct {
						Code string `json:"code"`
						Date string `json:"date"`
						Temp string `json:"temp"`
						Text string `json:"text"`
					} `json:"condition"`
				} `json:"item"`
			} `json:"channel"`
		} `json:"results"`
	} `json:"query"`
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
				So(response, json.ShouldMatchExpectedResponse, expected)
			})
		})
	})
}
