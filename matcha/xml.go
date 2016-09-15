package matcha

import (
	"fmt"
	"reflect"

	"github.com/clbanning/mxj"
)

func ShouldMatchExpectedXMLResponse(actual interface{}, expectedList ...interface{}) string {

	// Check number of arguments
	if len(expectedList) != 2 {
		return fmt.Sprintf("ShouldMatchExpectedXMLResponse expects two arguments: the expected XML format as a Struct, and a map to hold captured values")
	}

	actualXML, ok := actual.([]byte)
	if !ok {
		return fmt.Sprintf("Expected first argument to be a byte slice")
	}
	expectedResponseStruct := expectedList[0]
	var capturedValues map[string]interface{}
	if expectedList[1] != nil {
		capturedValues, ok = expectedList[1].(map[string]interface{})
		if !ok {
			return fmt.Sprintf("Expected third argument to be a map[string]interface or nil")
		}
	}
	//var actualResponse interface{}
	var actualResponse map[string]interface{}
	actualResponse, err := mxj.NewMapXml(actualXML, true)
	if err != nil {
		return fmt.Sprintf("Was not possible to unmarshal XML into a Go struct. XML data:\n%v", string(actualXML))
	}

	matcher := Matcher{format: "xml", capturedValues: capturedValues}

	result := matcher.shouldMatchExpectedField(actualResponse, reflect.TypeOf(expectedResponseStruct), "Result")
	if result != success {
		result = fmt.Sprintf("%v\nXML data:\n%v", result, string(actualXML))
	}
	return result
}
