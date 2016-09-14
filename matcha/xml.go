package matcha

import (
	"fmt"
	"reflect"

	"github.com/clbanning/mxj"
)

//func getJSONFieldNameXML(field reflect.StructField) string {
//        newFieldName, ok := field.Tag.Lookup("json")
//        if !ok {
//                // Get field name by looking at StructField name
//                newFieldName = snakecase.Snakecase(field.Name)
//        }
//        return newFieldName
//}

//func shouldMatchPatternXML(actual interface{}, expectedField reflect.StructField) string {

//        // Check if we are expecting to match against a pattern for this field
//        pattern, ok := expectedField.Tag.Lookup("pattern")
//        if ok {
//                // If so, check the expected field type is a string and the actual value is also a string
//                if expectedField.Type.Kind() != reflect.String {
//                        return fmt.Sprintf("'pattern' tag cannot be used on non-string fields: %v", expectedField.Name)
//                }
//                actualString, isString := actual.(string)
//                if !isString {
//                        return fmt.Sprintf("Expected a string value for field: %v but instead got %v", expectedField.Name, reflect.TypeOf(actual))
//                }

//                // If ok, then we try to match against the expected pattern
//                matched, err := regexp.MatchString(pattern, actualString)
//                if err != nil {
//                        return fmt.Sprintf("Received invalid regular expression: %v", pattern)
//                }
//                if !matched {
//                        return fmt.Sprintf("%v: '%v' does not match expected pattern: %v", expectedField.Name, actualString, pattern)
//                }
//        }

//        return success
//}

//func shouldMatchExpectedArrayXML(actual interface{}, expectedType reflect.Type, fieldName string, capturedValues map[string]interface{}) string {

//        var errorList []string
//        actualSlice, ok := actual.([]interface{})
//        if !ok {
//                return fmt.Sprintf("Was expecting an array for field: %v", fieldName)
//        }
//        // Get the expected type of each element in the array
//        expectedArrayElementType := expectedType.Elem()
//        // Compare each element in slice
//        for _, newActualField := range actualSlice {
//                // Array fields don't have names, so use something intuitive
//                newFieldName := fmt.Sprintf("%v array values", fieldName)
//                equal := shouldMatchExpectedFieldXML(newActualField, expectedArrayElementType, newFieldName, capturedValues)
//                if equal != success {
//                        errorList = append(errorList, equal)
//                }
//        }

//        if errorList != nil {
//                return strings.Join(errorList, "\n")
//        }
//        return success
//}

//func captureValueXML(capturedValues map[string]interface{}, expectedField reflect.StructField, value interface{}) {
//        // If we're not interested in capturing any values, just return
//        if capturedValues == nil {
//                return
//        }

//        captureKey, ok := expectedField.Tag.Lookup("capture")
//        if ok {
//                if captureKey == "" {
//                        captureKey = getJSONFieldName(expectedField)
//                }
//                capturedValues[captureKey] = value
//        }
//}

//func shouldMatchExpectedStructFieldXML(actual map[string]interface{}, expectedField reflect.StructField, capturedValues map[string]interface{}) string {

//        fieldName := getJSONFieldName(expectedField)
//        expectedFieldType := expectedField.Type
//        actualField, ok := actual[fieldName]
//        if !ok {
//                return fmt.Sprintf("No field '%v' found in response JSON", fieldName)
//        }

//        captureValueXML(capturedValues, expectedField, actualField)

//        equal := shouldMatchPatternXML(actualField, expectedField)
//        if equal != success {
//                return equal
//        }

//        return shouldMatchExpectedFieldXML(actualField, expectedFieldType, fieldName, capturedValues)
//}

//func shouldMatchExpectedObjectXML(actual interface{}, expectedType reflect.Type, fieldName string, capturedValues map[string]interface{}) string {

//        var errorList []string
//        fmt.Println("type = ", reflect.TypeOf(actual))
//        actualMap, ok := actual.(map[string]interface{})
//        if !ok {
//                return fmt.Sprintf("Was expecting a regular XML field: %v", fieldName)
//        }
//        for i := 0; i < expectedType.NumField(); i++ {

//                newField := expectedType.Field(i)
//                fmt.Println("new field = ", newField)
//                equal := shouldMatchExpectedStructFieldXML(actualMap, newField, capturedValues)
//                if equal != success {
//                        errorList = append(errorList, equal)
//                }
//        }

//        if errorList != nil {
//                return strings.Join(errorList, "\n")
//        }
//        return success
//}

//func shouldMatchExpectedFieldXML(actual interface{}, expectedType reflect.Type, fieldName string, capturedValues map[string]interface{}) string {

//        expectedKind := expectedType.Kind()
//        fmt.Println("expectedKind = ", expectedKind)
//        fmt.Println("actual = ", actual)
//        actualType := reflect.TypeOf(actual)
//        switch expectedKind {
//        case reflect.String:
//                if equal := ShouldEqual(expectedType, actualType); equal != success {
//                        return TypeErrorString(fieldName, expectedType.String(), actualType.String())
//                }
//        case reflect.Float64:
//                if equal := ShouldEqual(expectedType, actualType); equal != success {
//                        return TypeErrorString(fieldName, expectedType.String(), actualType.String())
//                }
//        case reflect.Bool:
//                if equal := ShouldEqual(expectedType, actualType); equal != success {
//                        return TypeErrorString(fieldName, expectedType.String(), actualType.String())
//                }
//        case reflect.Slice:
//                return shouldMatchExpectedArrayXML(actual, expectedType, fieldName, capturedValues)
//        case reflect.Struct:
//                // Type is a JSON object
//                return shouldMatchExpectedObjectXML(actual, expectedType, fieldName, capturedValues)
//        default:
//                return fmt.Sprintf("'%v' is of a type I don't know how to handle", expectedType)
//        }
//        return success
//}

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
	fmt.Println("actualResponse = ", actualResponse)

	matcher := Matcher{format: "xml", capturedValues: capturedValues}

	result := matcher.shouldMatchExpectedField(actualResponse, reflect.TypeOf(expectedResponseStruct), "Result")
	if result != success {
		result = fmt.Sprintf("%v\nXML data:\n%v", result, string(actualXML))
	}
	return result
}
