package utils

import (
	"strconv"
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

// returns map kay value or nil if not found, will be returned first found key value
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

// TODO: should be somewhere in other place
func GetSiteBackUrl() string {
	return "http://dev.ottemo.com:3000/"
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
