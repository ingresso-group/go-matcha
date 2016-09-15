package matcha

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

type expectedXMLString struct {
	StringField string
}

type expectedXMLNumber struct {
	NumberField float64
}

type expectedXMLBool struct {
	BooleanField bool
}

// It is invalid XML to have array elements at the root of the document
type expectedXMLArray struct {
	Result struct {
		ArrayField []string
	}
}

type expectedXMLComplexArray struct {
	Result struct {
		ArrayField []struct {
			NumberField float64
		}
	}
}

type expectedXMLComplex struct {
	Result struct {
		Class []struct {
			ClassCode string
			ClassId   float64
			Subclass  struct {
				SubclassCode string
			}
		}
		EventId string
	}
}

type expectedXMLFieldName struct {
	StringField string `xml:"string_t"`
}

type expectedXMLUrl struct {
	URL string `XML:"url" pattern:"https://.*"`
}

type expectedXMLCapture struct {
	Result struct {
		NumberField float64 `capture:"captured_number"`
		StringField string  `capture:""`
	}
}

func TestXMLGenericMatching(t *testing.T) {
	Convey("Given an expected field", t, func() {

		var expected expectedXMLString

		Convey("When invalid XML data", func() {
			fakeXML := []byte(`<a>`)

			Convey("It should return an error string", func() {
				success := ShouldMatchExpectedXMLResponse(fakeXML, expected, nil)
				So(success, ShouldStartWith, "Was not possible to unmarshal XML into a Go struct")
			})

		})

		Convey("When not present in actual XML", func() {
			fakeXML := []byte(`<hello></hello>`)

			Convey("It should return an error string", func() {
				success := ShouldMatchExpectedXMLResponse(fakeXML, expected, nil)
				So(success, ShouldStartWith, "No field 'string_field' found in response")
			})

		})

		Convey("When other fields present in actual XML", func() {

			fakeXML := []byte(`<string_field>some string</string_field><another_field>10</another_field>`)

			Convey("It should ignore them and mark as success", func() {
				success := ShouldMatchExpectedXMLResponse(fakeXML, expected, nil)
				So(success, ShouldEqual, "")
			})

		})

	})
}

func TestXMLStringMatching(t *testing.T) {

	Convey("Given an expected string field", t, func() {

		var expected expectedXMLString

		Convey("When has same type as actual XML", func() {

			fakeXML := []byte(`<string_field>some string</string_field>`)

			Convey("It should return success", func() {
				success := ShouldMatchExpectedXMLResponse(fakeXML, expected, nil)
				So(success, ShouldEqual, "")
			})

		})

		Convey("When has different type to actual XML", func() {

			fakeXML := []byte(`<string_field>5</string_field>`)

			Convey("It should return an error string", func() {
				success := ShouldMatchExpectedXMLResponse(fakeXML, expected, nil)
				expectedErrString := TypeErrorString("string_field", "string", "float64")
				So(success, ShouldStartWith, expectedErrString)
			})

		})

	})

}

func TestXMLNumberMatching(t *testing.T) {

	Convey("Given an expected float field", t, func() {

		var expected expectedXMLNumber

		Convey("When has same type as actual XML", func() {

			fakeXML := []byte(`<number_field>25.2</number_field>`)

			Convey("It should return success", func() {
				success := ShouldMatchExpectedXMLResponse(fakeXML, expected, nil)
				So(success, ShouldEqual, "")
			})

		})

		Convey("When has different type to actual XML", func() {

			fakeXML := []byte(`<number_field>fifty</number_field>`)

			Convey("It should return an error string", func() {
				success := ShouldMatchExpectedXMLResponse(fakeXML, expected, nil)
				expectedErrString := TypeErrorString("number_field", "float64", "string")
				So(success, ShouldStartWith, expectedErrString)
			})

		})

	})

}

func TestXMLBoolMatching(t *testing.T) {

	Convey("Given an expected boolean field", t, func() {

		var expected expectedXMLBool

		Convey("When has same type as actual XML", func() {

			fakeXML := []byte(`<boolean_field>true</boolean_field>`)

			Convey("It should return success", func() {
				success := ShouldMatchExpectedXMLResponse(fakeXML, expected, nil)
				So(success, ShouldEqual, "")
			})

		})

		Convey("When has different type to actual XML", func() {

			fakeXML := []byte(`<boolean_field>yes</boolean_field>`)

			Convey("It should return an error string", func() {
				success := ShouldMatchExpectedXMLResponse(fakeXML, expected, nil)
				expectedErrString := TypeErrorString("boolean_field", "bool", "string")
				So(success, ShouldStartWith, expectedErrString)
			})

		})

	})

}

func TestXMLArrayMatching(t *testing.T) {

	Convey("Given an expected array field", t, func() {

		var expected expectedXMLArray

		Convey("When has same type as actual XML", func() {

			fakeXML := []byte(`<result><array_field>one</array_field><array_field>two</array_field></result>`)

			Convey("It should return success", func() {
				success := ShouldMatchExpectedXMLResponse(fakeXML, expected, nil)
				So(success, ShouldEqual, "")
			})

		})

		// Not sure if this should be the expected behaviour or not
		Convey("When array field only has one element", func() {

			fakeXML := []byte(`<result><array_field>one</array_field></result>`)

			Convey("It should return an error string", func() {
				success := ShouldMatchExpectedXMLResponse(fakeXML, expected, nil)
				expectedErrString := "Was expecting an array for field: array_field"
				So(success, ShouldStartWith, expectedErrString)
			})

		})

		Convey("When one element has unexpcted type", func() {

			fakeXML := []byte(`<result><array_field>one</array_field><array_field>2</array_field></result>`)

			Convey("It should return an error string", func() {
				success := ShouldMatchExpectedXMLResponse(fakeXML, expected, nil)
				expectedErrString := "Expected 'array_field array values' to be: 'string' (but was: 'float64')!"
				So(success, ShouldStartWith, expectedErrString)
			})

		})

	})

	Convey("Given an expected array of structs", t, func() {

		var expected expectedXMLComplexArray

		Convey("When has same structure as actual XML", func() {

			fakeXML := []byte(`<result><array_field><number_field>1</number_field></array_field><array_field><number_field>2</number_field></array_field></result>`)

			Convey("It should return success", func() {
				success := ShouldMatchExpectedXMLResponse(fakeXML, expected, nil)
				So(success, ShouldEqual, "")
			})

		})

		Convey("When has different structure as actual XML", func() {

			fakeXML := []byte(`<result><array_field>one</array_field><array_field>two</array_field></result>`)

			Convey("It should return the expected error", func() {
				success := ShouldMatchExpectedXMLResponse(fakeXML, expected, nil)
				expectedErrString := "Was expecting an object for field: array_field"
				So(success, ShouldStartWith, expectedErrString)
			})

		})

	})

}

func TestXMLComplexMatching(t *testing.T) {

	Convey("Given an expected XML struct", t, func() {

		var expected expectedXMLComplex

		Convey("When has same structure as actual XML", func() {

			fakeXML := []byte(`
				<?xml version="1.0" encoding="UTF-8"?>
				<result>
				    <city_desc>London</city_desc>
				    <class>
					<class_code>concerts</class_code>
					<class_id>2</class_id>
					<is_main_class>yes</is_main_class>
					<subclass>
					    <is_main_subclass>yes</is_main_subclass>
					    <subclass_code>rock</subclass_code>
					    <subclass_desc>Rock &amp; Pop</subclass_desc>
					</subclass>
				    </class>
				    <class>
					<class_code>package</class_code>
					<class_id>7</class_id>
					<is_main_class>no</is_main_class>
					<subclass>
					    <is_main_subclass>no</is_main_subclass>
					    <subclass_code>misc</subclass_code>
					    <subclass_desc>Misc packages</subclass_desc>
					</subclass>
				    </class>
				    <country_code>uk</country_code>
				    <country_desc>United Kingdom</country_desc>
				    <event_id>9ZO</event_id>
				</result>
			`)
			Convey("It should return success", func() {
				success := ShouldMatchExpectedXMLResponse(fakeXML, expected, nil)
				So(success, ShouldEqual, "")
			})

		})

	})

}

func TestCustomXMLFieldName(t *testing.T) {

	Convey("Given an expected string field without 'xml' tag", t, func() {

		var expected expectedXMLFieldName

		Convey("When has same name and type as actual XML", func() {

			fakeXML := []byte(`<string_t>some string</string_t>`)

			Convey("It should return success", func() {
				success := ShouldMatchExpectedXMLResponse(fakeXML, expected, nil)
				So(success, ShouldEqual, "")
			})

		})

	})

}

func TestCapturingXMLValues(t *testing.T) {

	Convey("Given expected fields to capture", t, func() {

		var expected expectedXMLCapture

		Convey("When has same name and type as actual XML", func() {

			fakeXML := []byte(`<result><string_field>I've been captured!</string_field><number_field>16</number_field></result>`)

			Convey("Values should be captured", func() {
				capturedValues := make(map[string]interface{})
				success := ShouldMatchExpectedXMLResponse(fakeXML, expected, capturedValues)
				So(success, ShouldEqual, "")
				// Field with custom capture name
				So(capturedValues["captured_number"], ShouldEqual, 16)
				// Field with no capture name, should default to field name
				So(capturedValues["string_field"], ShouldEqual, "I've been captured!")
			})

		})

	})

}

func TestXMLPatternMatching(t *testing.T) {

	Convey("Given expected field with pattern", t, func() {

		var expected expectedXMLUrl

		Convey("When actual value matches pattern", func() {

			fakeXML := []byte(`<url>https://www.google.com</url>`)

			Convey("It should return success", func() {
				success := ShouldMatchExpectedXMLResponse(fakeXML, expected, nil)
				So(success, ShouldEqual, "")
			})

		})

		Convey("When actual value doesn't match pattern", func() {

			fakeXML := []byte(`<url>https:www.google.com</url>`)

			Convey("It should return an error string", func() {
				success := ShouldMatchExpectedXMLResponse(fakeXML, expected, nil)
				expectedErrString := "URL: 'https:www.google.com' does not match expected pattern: https://.*"
				So(success, ShouldStartWith, expectedErrString)
			})

		})

	})

}
