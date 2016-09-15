package matcha

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
		Success bool `json:"success"`
	} `json:"result"`
}

type expectedListOfObjects []ExpectedJSONComplex

type expectedNestedListOfObjects struct {
	Results []ExpectedJSONComplex `json:"results"`
}

type expectedFieldNoTag struct {
	StringField string
}

type expectedURL struct {
	URL string `json:"url" pattern:"https://.*"`
}

type expectedJSONCapture struct {
	NumberField float64 `capture:"captured_number"`
	StringField string  `capture:""`
}

func TestGenericMatching(t *testing.T) {
	Convey("Given an expected field", t, func() {

		var expected expectedJSONString

		Convey("When invalid JSON data", func() {
			fakeJSON := []byte(`{a}`)

			Convey("It should return an error string", func() {
				success := ShouldMatchExpectedJSONResponse(fakeJSON, expected, nil)
				So(success, ShouldStartWith, "Was not possible to unmarshal JSON into a Go struct")
			})

		})

		Convey("When not present in actual JSON", func() {
			fakeJSON := []byte(`{}`)

			Convey("It should return an error string", func() {
				success := ShouldMatchExpectedJSONResponse(fakeJSON, expected, nil)
				So(success, ShouldStartWith, "No field 'string_field' found in response")
			})

		})

		Convey("When other fields present in actual JSON", func() {

			fakeJSON := []byte(`{"string_field": "some string", "another_field": 10}`)

			Convey("It should ignore them and mark as success", func() {
				success := ShouldMatchExpectedJSONResponse(fakeJSON, expected, nil)
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
				success := ShouldMatchExpectedJSONResponse(fakeJSON, expected, nil)
				So(success, ShouldEqual, "")
			})

		})

		Convey("When has different type to actual JSON", func() {

			fakeJSON := []byte(`{"string_field": 5}`)

			Convey("It should return an error string", func() {
				success := ShouldMatchExpectedJSONResponse(fakeJSON, expected, nil)
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
				success := ShouldMatchExpectedJSONResponse(fakeJSON, expected, nil)
				So(success, ShouldEqual, "")
			})

		})

		Convey("When has different type to actual JSON", func() {

			fakeJSON := []byte(`{"number_field": "5"}`)

			Convey("It should return an error string", func() {
				success := ShouldMatchExpectedJSONResponse(fakeJSON, expected, nil)
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
				success := ShouldMatchExpectedJSONResponse(fakeJSON, expected, nil)
				So(success, ShouldEqual, "")
			})

		})

		Convey("When has different type to actual JSON", func() {

			fakeJSON := []byte(`{"boolean_field": "some string"}`)

			Convey("It should return an error string", func() {
				success := ShouldMatchExpectedJSONResponse(fakeJSON, expected, nil)
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
				success := ShouldMatchExpectedJSONResponse(fakeJSON, expected, nil)
				So(success, ShouldEqual, "")
			})

		})

		Convey("When has different type to actual JSON", func() {

			fakeJSON := []byte(`{"array_field": 5}`)

			Convey("It should return an error string", func() {
				success := ShouldMatchExpectedJSONResponse(fakeJSON, expected, nil)
				expectedErrString := "Was expecting an array for field: array_field"
				So(success, ShouldStartWith, expectedErrString)
			})

		})

		Convey("When one element has unexpcted type", func() {

			fakeJSON := []byte(`{"array_field": ["one", 2]}`)

			Convey("It should return an error string", func() {
				success := ShouldMatchExpectedJSONResponse(fakeJSON, expected, nil)
				expectedErrString := "Expected 'array_field array values' to be: 'string' (but was: 'float64')!"
				So(success, ShouldStartWith, expectedErrString)
			})

		})

	})

	Convey("Given an expected array of structs", t, func() {

		var expected expectedListOfObjects

		Convey("When has same structure as actual JSON", func() {

			fakeJSON := []byte(`[ { "result": { "attributes": { "string_field": "fantastic" }, "success": true } }, { "result": { "attributes": { "string_field": "wonderful" }, "success": true } } ]`)

			Convey("It should return success", func() {
				success := ShouldMatchExpectedJSONResponse(fakeJSON, expected, nil)
				So(success, ShouldEqual, "")
			})

		})

		Convey("When has different structure as actual JSON", func() {

			fakeJSON := []byte(`[ { "result": { "attributes": { "string_field": "fantastic" }, "success": true } }, { "result": { "attributes": { "string_field": "fantastic" } } } ]`)

			Convey("It should return the expected error", func() {
				success := ShouldMatchExpectedJSONResponse(fakeJSON, expected, nil)
				expectedErrString := "No field 'success' found in response"
				So(success, ShouldStartWith, expectedErrString)
			})

		})

		Convey("When there are several differences in structure to actual JSON", func() {

			fakeJSON := []byte(`[ { "result": {} }, { "result": { "attributes": { "string_field": "fantastic" } } } ]`)

			Convey("It should return several errors", func() {
				success := ShouldMatchExpectedJSONResponse(fakeJSON, expected, nil)
				expectedErrString := "No field 'attributes' found in response\nNo field 'success' found in response"
				So(success, ShouldStartWith, expectedErrString)
			})

		})
	})

	Convey("Given an expected nested array of structs", t, func() {

		var expected expectedNestedListOfObjects

		Convey("When has same structure as actual JSON", func() {

			fakeJSON := []byte(`{"results": [ { "result": { "attributes": { "string_field": "fantastic" }, "success": true } }, { "result": { "attributes": { "string_field": "wonderful" }, "success": true } } ]}`)

			Convey("It should return success", func() {
				success := ShouldMatchExpectedJSONResponse(fakeJSON, expected, nil)
				So(success, ShouldEqual, "")
			})

		})
	})
}

func TestJSONObjectMatching(t *testing.T) {

	Convey("Given an expected nested struct", t, func() {

		var expected ExpectedJSONComplex

		Convey("When matches actual JSON structure", func() {

			fakeJSON := []byte(`{"result": {"attributes":{ "string_field": "fantastic"}, "success": true } }`)

			Convey("It should return success", func() {
				success := ShouldMatchExpectedJSONResponse(fakeJSON, expected, nil)
				So(success, ShouldEqual, "")
			})

		})

		Convey("When has different type to actual JSON", func() {

			fakeJSON := []byte(`{"result": false}`)

			Convey("It should return an error string", func() {
				success := ShouldMatchExpectedJSONResponse(fakeJSON, expected, nil)
				expectedErrString := "Was expecting an object for field: result"
				So(success, ShouldStartWith, expectedErrString)
			})

		})

		Convey("When received an array instead", func() {

			fakeJSON := []byte(`{"result": [1, 2, 3]}`)

			Convey("It should return an error string", func() {
				success := ShouldMatchExpectedJSONResponse(fakeJSON, expected, nil)
				expectedErrString := "Was expecting an object for field: result"
				So(success, ShouldStartWith, expectedErrString)
			})

		})

	})

}

func TestDefaultFieldName(t *testing.T) {

	Convey("Given an expected string field without 'json' tag", t, func() {

		var expected expectedFieldNoTag

		Convey("When has same name and type as actual JSON", func() {

			fakeJSON := []byte(`{"string_field": "some string"}`)

			Convey("It should return success", func() {
				success := ShouldMatchExpectedJSONResponse(fakeJSON, expected, nil)
				So(success, ShouldEqual, "")
			})

		})

	})

}

func TestCapturingJSONValues(t *testing.T) {

	Convey("Given expected fields to capture", t, func() {

		var expected expectedJSONCapture

		Convey("When has same name and type as actual JSON", func() {

			fakeJSON := []byte(`{"string_field": "I've been captured!", "number_field": 16}`)

			Convey("Values should be captured", func() {
				capturedValues := make(map[string]interface{})
				success := ShouldMatchExpectedJSONResponse(fakeJSON, expected, capturedValues)
				So(success, ShouldEqual, "")
				// Field with custom capture name
				So(capturedValues["captured_number"], ShouldEqual, 16)
				// Field with no capture name, should default to field name
				So(capturedValues["string_field"], ShouldEqual, "I've been captured!")
			})

		})

	})

}

func TestJSONPatternMatching(t *testing.T) {

	Convey("Given expected field with pattern", t, func() {

		var expected expectedURL

		Convey("When actual value matches pattern", func() {

			fakeJSON := []byte(`{"url": "https://www.google.com"}`)

			Convey("It should return success", func() {
				success := ShouldMatchExpectedJSONResponse(fakeJSON, expected, nil)
				So(success, ShouldEqual, "")
			})

		})

		Convey("When actual value doesn't match pattern", func() {

			fakeJSON := []byte(`{"url": "https:www.google.com"}`)

			Convey("It should return an error string", func() {
				success := ShouldMatchExpectedJSONResponse(fakeJSON, expected, nil)
				expectedErrString := "URL: 'https:www.google.com' does not match expected pattern: https://.*"
				So(success, ShouldStartWith, expectedErrString)
			})

		})

	})

}
