package json

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

type expectedJSONString struct {
	StringField string `json:"string_field"`
}

type expectedJSONNumber struct {
	NumberField float64 `json:"number_field"`
}

type expectedJSONBool struct {
	BooleanField bool `json:"boolean_field"`
}

type expectedJSONArray struct {
	ArrayField []string `json:"array_field"`
}

type ExpectedJSONComplex struct {
	Result struct {
		Attributes struct {
			StringField string `json:"string_field"`
		} `json:"attributes"`
	} `json:"result"`
}

type expectedListOfObjects []ExpectedJSONComplex

type expectedNestedListOfObjects struct {
	Results []ExpectedJSONComplex `json:"results"`
}

func TestGenericMatching(t *testing.T) {
	Convey("Given an expected field", t, func() {

		var expected expectedJSONString

		Convey("When invalid JSON data", func() {
			fakeJSON := []byte(`{a}`)

			Convey("It should return an error string", func() {
				success := ShouldMatchExpectedResponse(fakeJSON, expected)
				So(success, ShouldStartWith, "Was not possible to unmarshal JSON into a Go struct")
			})

		})

		Convey("When not present in actual JSON", func() {
			fakeJSON := []byte(`{}`)

			Convey("It should return an error string", func() {
				success := ShouldMatchExpectedResponse(fakeJSON, expected)
				So(success, ShouldStartWith, "No field 'string_field' found in response JSON")
			})

		})

		Convey("When other fields present in actual JSON", func() {

			fakeJSON := []byte(`{"string_field": "some string", "another_field": 10}`)

			Convey("It should ignore them and mark as success", func() {
				success := ShouldMatchExpectedResponse(fakeJSON, expected)
				So(success, ShouldEqual, "")
			})

		})

	})
}

func TestJSONStringMatching(t *testing.T) {

	Convey("Given an expected string field", t, func() {

		var expected expectedJSONString

		Convey("When has same type as actual JSON", func() {

			fakeJSON := []byte(`{"string_field": "some string"}`)

			Convey("It should return success", func() {
				success := ShouldMatchExpectedResponse(fakeJSON, expected)
				So(success, ShouldEqual, "")
			})

		})

		Convey("When has different type to actual JSON", func() {

			fakeJSON := []byte(`{"string_field": 5}`)

			Convey("It should return an error string", func() {
				success := ShouldMatchExpectedResponse(fakeJSON, expected)
				expectedErrString := TypeErrorString("string_field", "string", "float64")
				So(success, ShouldStartWith, expectedErrString)
			})

		})

	})

}

func TestJSONNumberMatching(t *testing.T) {

	Convey("Given an expected float field", t, func() {

		var expected expectedJSONNumber

		Convey("When has same type as actual JSON", func() {

			fakeJSON := []byte(`{"number_field": 25.2}`)

			Convey("It should return success", func() {
				success := ShouldMatchExpectedResponse(fakeJSON, expected)
				So(success, ShouldEqual, "")
			})

		})

		Convey("When has different type to actual JSON", func() {

			fakeJSON := []byte(`{"number_field": "5"}`)

			Convey("It should return an error string", func() {
				success := ShouldMatchExpectedResponse(fakeJSON, expected)
				expectedErrString := TypeErrorString("number_field", "float64", "string")
				So(success, ShouldStartWith, expectedErrString)
			})

		})

	})

}

func TestJSONBoolMatching(t *testing.T) {

	Convey("Given an expected boolean field", t, func() {

		var expected expectedJSONBool

		Convey("When has same type as actual JSON", func() {

			fakeJSON := []byte(`{"boolean_field": true}`)

			Convey("It should return success", func() {
				success := ShouldMatchExpectedResponse(fakeJSON, expected)
				So(success, ShouldEqual, "")
			})

		})

		Convey("When has different type to actual JSON", func() {

			fakeJSON := []byte(`{"boolean_field": "some string"}`)

			Convey("It should return an error string", func() {
				success := ShouldMatchExpectedResponse(fakeJSON, expected)
				expectedErrString := TypeErrorString("boolean_field", "bool", "string")
				So(success, ShouldStartWith, expectedErrString)
			})

		})

	})

}

// At the moment, array comparisons are very basic
func TestJSONArrayMatching(t *testing.T) {

	Convey("Given an expected array field", t, func() {

		var expected expectedJSONArray

		Convey("When has same type as actual JSON", func() {

			fakeJSON := []byte(`{"array_field": ["one", "two"]}`)

			Convey("It should return success", func() {
				success := ShouldMatchExpectedResponse(fakeJSON, expected)
				So(success, ShouldEqual, "")
			})

		})

		Convey("When has different type to actual JSON", func() {

			fakeJSON := []byte(`{"array_field": 5}`)

			Convey("It should return an error string", func() {
				success := ShouldMatchExpectedResponse(fakeJSON, expected)
				expectedErrString := "Was expecting an array for field: array_field"
				So(success, ShouldStartWith, expectedErrString)
			})

		})

		Convey("When one element has unexpcted type", func() {

			fakeJSON := []byte(`{"array_field": ["one", 2]}`)

			Convey("It should return an error string", func() {
				success := ShouldMatchExpectedResponse(fakeJSON, expected)
				expectedErrString := "Expected 'array_field array values' to be: 'string' (but was: 'float64')!"
				So(success, ShouldStartWith, expectedErrString)
			})

		})

	})

	Convey("Given an expected array of structs", t, func() {

		var expected expectedListOfObjects

		Convey("When has same structure as actual JSON", func() {

			fakeJSON := []byte(`[ { "result": { "attributes": { "string_field": "fantastic" } } }, { "result": { "attributes": { "string_field": "wonderful" } } } ]`)

			Convey("It should return success", func() {
				success := ShouldMatchExpectedResponse(fakeJSON, expected)
				So(success, ShouldEqual, "")
			})

		})

		Convey("When has different structure as actual JSON", func() {

			fakeJSON := []byte(`[ { "result": { "attributes": { "string_field": "fantastic" } } }, { "result": { "attributes": 0 } } ]`)

			Convey("It should return success", func() {
				success := ShouldMatchExpectedResponse(fakeJSON, expected)
				expectedErrString := "Was expecting a JSON object for field: attributes"
				So(success, ShouldStartWith, expectedErrString)
			})

		})
	})

	Convey("Given an expected nested array of structs", t, func() {

		var expected expectedNestedListOfObjects

		Convey("When has same structure as actual JSON", func() {

			fakeJSON := []byte(`{"results": [ { "result": { "attributes": { "string_field": "fantastic" } } }, { "result": { "attributes": { "string_field": "wonderful" } } } ]}`)

			Convey("It should return success", func() {
				success := ShouldMatchExpectedResponse(fakeJSON, expected)
				So(success, ShouldEqual, "")
			})

		})
	})
}

func TestJSONObjectMatching(t *testing.T) {

	Convey("Given an expected nested struct", t, func() {

		var expected ExpectedJSONComplex

		Convey("When matches actual JSON structure", func() {

			fakeJSON := []byte(`{"result": {"attributes":{ "string_field": "fantastic"} } }`)

			Convey("It should return success", func() {
				success := ShouldMatchExpectedResponse(fakeJSON, expected)
				So(success, ShouldEqual, "")
			})

		})

		Convey("When has different type to actual JSON", func() {

			fakeJSON := []byte(`{"result": false}`)

			Convey("It should return an error string", func() {
				success := ShouldMatchExpectedResponse(fakeJSON, expected)
				expectedErrString := "Was expecting a JSON object for field: result"
				So(success, ShouldStartWith, expectedErrString)
			})

		})

		Convey("When received an array instead", func() {

			fakeJSON := []byte(`{"result": [1, 2, 3]}`)

			Convey("It should return an error string", func() {
				success := ShouldMatchExpectedResponse(fakeJSON, expected)
				expectedErrString := "Was expecting a JSON object for field: result"
				So(success, ShouldStartWith, expectedErrString)
			})

		})

	})

}

// TODO - add tests for multiple errors
