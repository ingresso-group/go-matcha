package matcha

import (
	"encoding/json"
	"fmt"
	"reflect"
)

func ShouldMatchExpectedJSONResponse(actual interface{}, expectedList ...interface{}) string {

	// Check number of arguments
	if len(expectedList) != 2 {
		return fmt.Sprintf("ShouldMatchExpectedJSONResponse expects three arguments: the actual JSON as a byte slice, the expected JSON format as a Struct, and a map to hold captured values")
	}

	actualJSON, ok := actual.([]byte)
	if !ok {
		return fmt.Sprintf("Expected first argument to be a byte slice")
	}
	expectedResponseStruct := expectedList[0]
	var capturedValues CapturedValues
	if expectedList[1] != nil {
		capturedValues, ok = expectedList[1].(CapturedValues)
		if !ok {
			return fmt.Sprintf("Expected third argument to be a map[string]interface or nil")
		}
	}
	var actualResponse interface{}
	err := json.Unmarshal(actualJSON, &actualResponse)
	if err != nil {
		return fmt.Sprintf("Was not possible to unmarshal JSON into a Go struct. JSON data:\n%v", string(actualJSON))
	}

	matcher := Matcher{format: "json", capturedValues: capturedValues}

	result := matcher.shouldMatchExpectedField(actualResponse, reflect.TypeOf(expectedResponseStruct), "Result")
	if result != success {
		result = fmt.Sprintf("%v\nJSON data:\n%v", result, string(actualJSON))
	}
	return result
}
