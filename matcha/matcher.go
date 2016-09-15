package matcha

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"

	snakecase "github.com/segmentio/go-snakecase"
	. "github.com/smartystreets/goconvey/convey"
)

type Matcher struct {
	format         string // Should be 'json' or 'xml'
	capturedValues map[string]interface{}
}

const (
	success = "" // goconvey uses an empty string to signal success
)

func TypeErrorString(fieldName string, expectedType string, actualType string) string {
	return fmt.Sprintf("Expected '%v' to be: '%v' (but was: '%v')!", fieldName, expectedType, actualType)
}

func (m *Matcher) getFieldName(field reflect.StructField) string {
	dataType := m.format
	newFieldName, ok := field.Tag.Lookup(dataType)
	if !ok {
		// Get field name by looking at StructField name
		newFieldName = snakecase.Snakecase(field.Name)
	}
	return newFieldName
}

func (m *Matcher) shouldMatchPattern(actual interface{}, expectedField reflect.StructField) string {

	// Check if we are expecting to match against a pattern for this field
	pattern, ok := expectedField.Tag.Lookup("pattern")
	if ok {
		// If so, check the expected field type is a string and the actual value is also a string
		if expectedField.Type.Kind() != reflect.String {
			return fmt.Sprintf("'pattern' tag cannot be used on non-string fields: %v", expectedField.Name)
		}
		actualString, isString := actual.(string)
		if !isString {
			return fmt.Sprintf("Expected a string value for field: %v but instead got %v", expectedField.Name, reflect.TypeOf(actual))
		}

		// If ok, then we try to match against the expected pattern
		matched, err := regexp.MatchString(pattern, actualString)
		if err != nil {
			return fmt.Sprintf("Received invalid regular expression: %v", pattern)
		}
		if !matched {
			return fmt.Sprintf("%v: '%v' does not match expected pattern: %v", expectedField.Name, actualString, pattern)
		}
	}

	return success
}

func (m *Matcher) shouldMatchExpectedArray(actual interface{}, expectedType reflect.Type, fieldName string) string {

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
		equal := m.shouldMatchExpectedField(newActualField, expectedArrayElementType, newFieldName)
		if equal != success {
			errorList = append(errorList, equal)
		}
	}

	if errorList != nil {
		return strings.Join(errorList, "\n")
	}
	return success
}

func (m *Matcher) captureValue(expectedField reflect.StructField, value interface{}) {
	// If we're not interested in capturing any values, just return
	if m.capturedValues == nil {
		return
	}

	captureKey, ok := expectedField.Tag.Lookup("capture")
	if ok {
		if captureKey == "" {
			captureKey = m.getFieldName(expectedField)
		}
		m.capturedValues[captureKey] = value
	}
}

func (m *Matcher) shouldMatchExpectedStructField(actual map[string]interface{}, expectedField reflect.StructField) string {

	fieldName := m.getFieldName(expectedField)
	expectedFieldType := expectedField.Type
	actualField, ok := actual[fieldName]
	if !ok {
		return fmt.Sprintf("No field '%v' found in response", fieldName)
	}

	m.captureValue(expectedField, actualField)

	equal := m.shouldMatchPattern(actualField, expectedField)
	if equal != success {
		return equal
	}

	return m.shouldMatchExpectedField(actualField, expectedFieldType, fieldName)
}

func (m *Matcher) shouldMatchExpectedObject(actual interface{}, expectedType reflect.Type, fieldName string) string {

	var errorList []string
	actualMap, ok := actual.(map[string]interface{})
	if !ok {
		return fmt.Sprintf("Was expecting an object for field: %v", fieldName)
	}
	for i := 0; i < expectedType.NumField(); i++ {

		newField := expectedType.Field(i)
		equal := m.shouldMatchExpectedStructField(actualMap, newField)
		if equal != success {
			errorList = append(errorList, equal)
		}
	}

	if errorList != nil {
		return strings.Join(errorList, "\n")
	}
	return success
}

func (m *Matcher) shouldMatchExpectedField(actual interface{}, expectedType reflect.Type, fieldName string) string {

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
		return m.shouldMatchExpectedArray(actual, expectedType, fieldName)
	case reflect.Struct:
		// Type is a JSON object
		return m.shouldMatchExpectedObject(actual, expectedType, fieldName)
	default:
		return fmt.Sprintf("'%v' is of a type I don't know how to handle", expectedType)
	}
	return success
}
