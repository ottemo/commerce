package utils

import (
	"math"
	"strconv"
	"strings"
	"time"
)

// checks if value is blank (zero value)
func CheckIsBlank(value interface{}) bool {
	switch typedValue := value.(type) {
	case string:
		return typedValue == ""
	case int, int32, int64:
		return typedValue == 0
	case []string, []interface{}:
		return false
	case map[string]interface{}, map[string]string, map[string]int:
		return false
	case float32, float64:
		return typedValue == 0
	}

	return false
}

// checks presence of non blank values for keys in map
//   - first arg must be map
//   - fallowing arguments are map keys you want to check
func KeysInMapAndNotBlank(mapObject interface{}, keys ...interface{}) bool {
	switch typedMapObject := mapObject.(type) {
	case map[string]interface{}:
		for _, key := range keys {
			if key2, ok := key.(string); ok {
				if mapValue, present := typedMapObject[key2]; present {
					if isBlank := CheckIsBlank(mapValue); isBlank {
						return false
					}
				} else {
					return false
				}
			}
		}
	case map[int]interface{}:
		for _, key := range keys {
			if key, ok := key.(int); ok {
				if mapValue, present := typedMapObject[key]; present {
					if isBlank := CheckIsBlank(mapValue); isBlank {
						return false
					}
				} else {
					return false
				}
			}
		}
	}

	return true
}

// returns map key value or nil if not found, will be returned first found key value
func GetFirstMapValue(mapObject interface{}, keys ...string) interface{} {
	switch typedMapObject := mapObject.(type) {
	case map[string]interface{}:
		for _, key := range keys {
			if value, present := typedMapObject[key]; present {
				return value
			}
		}
	case map[string]string:
		for _, key := range keys {
			if value, present := typedMapObject[key]; present {
				return value
			}
		}
	}

	return nil
}

// searches for item in array/slice, returns true if found
func IsInArray(searchValue interface{}, arrayObject interface{}) bool {
	switch typedObject := arrayObject.(type) {
	case []string:
		searchValue, ok := searchValue.(string)
		if !ok {
			return false
		}

		for _, value := range typedObject {
			if value == searchValue {

				return true
			}
		}

	case []interface{}:
		for _, value := range typedObject {
			if value == searchValue {
				return true
			}
		}
	}

	return false
}

// searches for presence of 1-st arg string option among provided options since 2-nd argument
func IsAmongStr(option string, searchOptions ...string) bool {
	for _, listOption := range searchOptions {
		if option == listOption {
			return true
		}
	}
	return false
}

// searches for a string in []string slice
func IsInListStr(searchItem string, searchList []string) bool {
	for _, listItem := range searchList {
		if listItem == searchItem {
			return true
		}
	}
	return false
}

// checks presence of string keys in map
func StrKeysInMap(mapObject interface{}, keys ...string) bool {
	switch typedMap := mapObject.(type) {
	case map[string]interface{}: //, map[string]string:

		for _, key := range keys {
			if _, present := typedMap[key]; !present {
				return false
			}
		}
	}

	return true
}

// returns trimmed []string array of [separators] delimited values
func Explode(value string, separators string) []string {
	result := make([]string, 0)

	splitResult := strings.Split(value, separators)
	for _, arrayValue := range splitResult {
		arrayValue = strings.TrimSpace(arrayValue)
		if arrayValue != "" {
			result = append(result, arrayValue)
		}
	}

	return result
}

// checks if value is MD5
func IsMD5(value string) bool {
	ok := false
	if len(value) == 32 {
		ok = true
		for i := 0; i < 32; i++ {
			c := value[i]
			if !(c == '1' || c == '2' || c == '3' || c == '4' || c == '5' || c == '6' || c == '7' || c == '8' || c == '9' || c == '0' ||
				c == 'a' || c == 'b' || c == 'c' || c == 'd' || c == 'e' || c == 'f') {
				ok = false
				break
			}
		}
	}
	return ok
}

// converts interface{} to string
func InterfaceToBool(value interface{}) bool {
	switch typedValue := value.(type) {
	case bool:
		return typedValue
	case string:
		if boolValue, err := strconv.ParseBool(typedValue); err == nil {
			return boolValue
		} else {
			return false
		}
	case int:
		return typedValue > 0
	default:
		return false
	}
}

// converts interface{} to map[string]interface{}
func InterfaceToArray(value interface{}) []interface{} {
	result := make([]interface{}, 0)

	switch typedValue := value.(type) {
	case []string:
		result = make([]interface{}, len(typedValue))
		for idx, value := range typedValue {
			result[idx] = value
		}
	case []int:
		result = make([]interface{}, len(typedValue))
		for idx, value := range typedValue {
			result[idx] = value
		}
	case []int64:
		result = make([]interface{}, len(typedValue))
		for idx, value := range typedValue {
			result[idx] = value
		}
	case []float64:
		result = make([]interface{}, len(typedValue))
		for idx, value := range typedValue {
			result[idx] = value
		}
	case []interface{}:
		return typedValue

	case string:
		jsonArray, err := DecodeJsonToArray(typedValue)
		if err == nil {
			return jsonArray
		}

		splitValues := strings.Split(typedValue, ",")
		result = make([]interface{}, len(splitValues))
		for idx, value := range splitValues {
			result[idx] = strings.Trim(value, " \t\n")
		}
	}

	result = append(result, value)
	return result
}

// converts interface{} to map[string]interface{}
func InterfaceToMap(value interface{}) map[string]interface{} {
	switch typedValue := value.(type) {
	case map[string]interface{}:
		return typedValue

	case string:
		value, err := DecodeJsonToStringKeyMap(value)
		if err == nil {
			return value
		}
	}

	return make(map[string]interface{})
}

// converts interface{} to string
func InterfaceToString(value interface{}) string {
	switch value := value.(type) {
	case bool:
		return strconv.FormatBool(value)
	case int:
		return strconv.Itoa(value)
	case int64:
		return strconv.FormatInt(value, 10)
	case float64:
		return strconv.FormatFloat(value, 'f', 6, 64)
	case string:
		return value
	default:
		return ""
	}
}

// converts interface{} to integer
func InterfaceToInt(value interface{}) int {
	switch typedValue := value.(type) {
	case int:
		return typedValue
	case string:
		intValue, _ := strconv.ParseInt(typedValue, 10, 64)
		return int(intValue)
	case float64:
		return int(typedValue)
	default:
		return 0
	}
}

// converts interface{} to float64
func InterfaceToFloat64(value interface{}) float64 {
	switch typedValue := value.(type) {
	case float64:
		return typedValue
	case int64:
		return float64(typedValue)
	case int:
		return float64(typedValue)
	case string:
		floatValue, _ := strconv.ParseFloat(typedValue, 64)
		return floatValue
	default:
		return 0.0
	}
}

// converts interface{} to time.Time
func InterfaceToTime(value interface{}) time.Time {
	switch typedValue := value.(type) {
	case time.Time:
		return typedValue
	case string:
		newValue, err := time.Parse(time.UnixDate, typedValue)
		if err == nil {
			return newValue
		}

		newValue, err = time.Parse(time.RFC3339, typedValue)
		if err == nil {
			return newValue
		}

		newValue, err = time.Parse(time.RFC822Z, typedValue)
		if err == nil {
			return newValue
		}
	}

	return time.Unix(0, 0)
}

// convert string to integer
func StringToInteger(value string) (int, error) {
	return strconv.Atoi(value)
}

// convert string to float64
func StringToFloat(value string) (float64, error) {
	return strconv.ParseFloat(value, 64)
}

// rounds value to given precision (roundOn=0.5 usual cases)
func Round(val float64, roundOn float64, places int) float64 {

	pow := math.Pow(10, float64(places))
	digit := pow * val
	_, div := math.Modf(digit)

	var round float64
	if div >= roundOn {
		round = math.Ceil(digit)
	} else {
		round = math.Floor(digit)
	}

	return round / pow
}

// function you should use to normalize price after calculations
// (in future that function should have config options setup)
func RoundPrice(price float64) float64 {
	return Round(price, 0.5, 2)
}
