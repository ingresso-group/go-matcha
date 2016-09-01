package json

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	. "github.com/smartystreets/goconvey/convey"
)

const (
	success = "" // goconvey uses an empty string to signal success
)

func TypeErrorString(fieldName string, expectedType string, actualType string) string {
	return fmt.Sprintf("Expected '%v' to be: '%v' (but was: '%v')!", fieldName, expectedType, actualType)
}

func shouldMatchExpectedArray(actual interface{}, expected interface{}, fieldName string) string {

	var errorList []string
	expectedType := reflect.TypeOf(expected)
	actualSlice, ok := actual.([]interface{})
	if !ok {
		return fmt.Sprintf("Was expecting an array for field: %v", fieldName)
	}
	// Get the expected type of each element in the array
	expectedArrayElementType := expectedType.Elem()
	newExpectedField := reflect.Zero(expectedArrayElementType).Interface()
	// Compare each element in slice
	for _, newActualField := range actualSlice {
		// Array fields don't have names, so use something intuitive
		newFieldName := fmt.Sprintf("%v array values", fieldName)
		equal := shouldMatchExpectedField(newActualField, newExpectedField, newFieldName)
		if equal != success {
			errorList = append(errorList, equal)
		}
	}

	if errorList != nil {
		return strings.Join(errorList, "\n")
	}
	return success
}

func shouldMatchExpectedObject(actual interface{}, expected interface{}, fieldName string) string {

	var errorList []string
	expectedType := reflect.TypeOf(expected)
	expectedValue := reflect.ValueOf(expected)
	actualMap, ok := actual.(map[string]interface{})
	if !ok {
		return fmt.Sprintf("Was expecting a JSON object for field: %v", fieldName)
	}
	for i := 0; i < expectedType.NumField(); i++ {
		newFieldName := expectedType.Field(i).Tag.Get("json")
		newExpectedField := expectedValue.Field(i).Interface()
		newActualField, ok := actualMap[newFieldName]
		if !ok {
			return fmt.Sprintf("No field '%v' found in response JSON", newFieldName)
		}
		equal := shouldMatchExpectedField(newActualField, newExpectedField, newFieldName)
		if equal != success {
			errorList = append(errorList, equal)
		}
	}

	if errorList != nil {
		return strings.Join(errorList, "\n")
	}
	return success
}

func shouldMatchExpectedField(actual interface{}, expected interface{}, fieldName string) string {

	expectedType := reflect.TypeOf(expected)
	expectedValue := reflect.ValueOf(expected)
	actualType := reflect.TypeOf(actual)

	switch expected.(type) {
	case string:
		if equal := ShouldHaveSameTypeAs(expected, actual); equal != success {
			return TypeErrorString(fieldName, expectedType.String(), actualType.String())
		}
	case float64:
		if equal := ShouldHaveSameTypeAs(expected, actual); equal != success {
			return TypeErrorString(fieldName, expectedType.String(), actualType.String())
		}
	case bool:
		if equal := ShouldHaveSameTypeAs(expected, actual); equal != success {
			return TypeErrorString(fieldName, expectedType.String(), actualType.String())
		}
	case interface{}:

		if expectedType.Kind() == reflect.Slice {
			return shouldMatchExpectedArray(actual, expected, fieldName)
		} else {
			// Type is a JSON object
			return shouldMatchExpectedObject(actual, expected, fieldName)
		}

	default:
		fmt.Println(expectedType, "is of a type I don't know how to handle", expectedValue)
	}

	return success
}

func ShouldMatchExpectedResponse(actual interface{}, expectedList ...interface{}) string {

	actualJSON := actual.([]byte)
	expectedResponseStruct := expectedList[0]
	var actualResponse interface{}
	err := json.Unmarshal(actualJSON, &actualResponse)
	if err != nil {
		return fmt.Sprintf("Was not possible to unmarshal JSON into a Go struct. JSON data:\n%v", string(actualJSON))
	}

	result := shouldMatchExpectedField(actualResponse, expectedResponseStruct, "Result")
	if result != success {
		result = fmt.Sprintf("%v\nJSON data:\n%v", result, string(actualJSON))
	}
	return result
}
