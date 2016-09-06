package matcha

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	snakecase "github.com/segmentio/go-snakecase"
	. "github.com/smartystreets/goconvey/convey"
)

const (
	success = "" // goconvey uses an empty string to signal success
)

func TypeErrorString(fieldName string, expectedType string, actualType string) string {
	return fmt.Sprintf("Expected '%v' to be: '%v' (but was: '%v')!", fieldName, expectedType, actualType)
}

func getJSONFieldName(field reflect.StructField) string {
	newFieldName, ok := field.Tag.Lookup("json")
	if !ok {
		// Get field name by looking at StructField name
		newFieldName = snakecase.Snakecase(field.Name)
	}
	return newFieldName
}

func shouldMatchExpectedArray(actual interface{}, expectedType reflect.Type, fieldName string, capturedValues map[string]interface{}) string {

	var errorList []string
	actualSlice, ok := actual.([]interface{})
	if !ok {
		return fmt.Sprintf("Was expecting an array for field: %v", fieldName)
	}
	// Get the expected type of each element in the array
	expectedArrayElementType := expectedType.Elem()
	// Compare each element in slice
	for _, newActualField := range actualSlice {
		// Array fields don't have names, so use something intuitive
		newFieldName := fmt.Sprintf("%v array values", fieldName)
		equal := shouldMatchExpectedField(newActualField, expectedArrayElementType, newFieldName, capturedValues)
		if equal != success {
			errorList = append(errorList, equal)
		}
	}

	if errorList != nil {
		return strings.Join(errorList, "\n")
	}
	return success
}

func captureValue(capturedValues map[string]interface{}, key string, value interface{}) {
	if capturedValues == nil {
		return
	}
	capturedValues[key] = value
}

func shouldMatchExpectedStructField(actual map[string]interface{}, expectedField reflect.StructField, capturedValues map[string]interface{}) string {

	fieldName := getJSONFieldName(expectedField)
	expectedFieldType := expectedField.Type
	actualField, ok := actual[fieldName]
	if !ok {
		return fmt.Sprintf("No field '%v' found in response JSON", fieldName)
	}

	captureKey, ok := expectedField.Tag.Lookup("capture")
	if ok {
		if captureKey == "" {
			captureKey = getJSONFieldName(expectedField)
		}
		captureValue(capturedValues, captureKey, actualField)
	}
	return shouldMatchExpectedField(actualField, expectedFieldType, fieldName, capturedValues)
}

func shouldMatchExpectedObject(actual interface{}, expectedType reflect.Type, fieldName string, capturedValues map[string]interface{}) string {

	var errorList []string
	actualMap, ok := actual.(map[string]interface{})
	if !ok {
		return fmt.Sprintf("Was expecting a JSON object for field: %v", fieldName)
	}
	for i := 0; i < expectedType.NumField(); i++ {

		newField := expectedType.Field(i)
		equal := shouldMatchExpectedStructField(actualMap, newField, capturedValues)
		if equal != success {
			errorList = append(errorList, equal)
		}
	}

	if errorList != nil {
		return strings.Join(errorList, "\n")
	}
	return success
}

func shouldMatchExpectedField(actual interface{}, expectedType reflect.Type, fieldName string, capturedValues map[string]interface{}) string {

	expectedKind := expectedType.Kind()
	actualType := reflect.TypeOf(actual)
	switch expectedKind {
	case reflect.String:
		if equal := ShouldEqual(expectedType, actualType); equal != success {
			return TypeErrorString(fieldName, expectedType.String(), actualType.String())
		}
	case reflect.Float64:
		if equal := ShouldEqual(expectedType, actualType); equal != success {
			return TypeErrorString(fieldName, expectedType.String(), actualType.String())
		}
	case reflect.Bool:
		if equal := ShouldEqual(expectedType, actualType); equal != success {
			return TypeErrorString(fieldName, expectedType.String(), actualType.String())
		}
	case reflect.Slice:
		return shouldMatchExpectedArray(actual, expectedType, fieldName, capturedValues)
	case reflect.Struct:
		// Type is a JSON object
		return shouldMatchExpectedObject(actual, expectedType, fieldName, capturedValues)
	default:
		fmt.Println(expectedType, "is of a type I don't know how to handle")
	}
	return success
}

func ShouldMatchExpectedResponse(actual interface{}, expectedList ...interface{}) string {

	// Check number of arguments
	if len(expectedList) != 2 {
		return fmt.Sprintf("ShouldMatchExpectedResponse expects two arguments: the expected JSON format as a Struct, and a map to hold captured values")
	}

	actualJSON, ok := actual.([]byte)
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
	var actualResponse interface{}
	err := json.Unmarshal(actualJSON, &actualResponse)
	if err != nil {
		return fmt.Sprintf("Was not possible to unmarshal JSON into a Go struct. JSON data:\n%v", string(actualJSON))
	}

	result := shouldMatchExpectedField(actualResponse, reflect.TypeOf(expectedResponseStruct), "Result", capturedValues)
	if result != success {
		result = fmt.Sprintf("%v\nJSON data:\n%v", result, string(actualJSON))
	}
	return result
}
