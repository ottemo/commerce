package utils

import (
	"strconv"
)

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
		if  listItem == searchItem {
			return true
		}
	}
	return false
}

// checks presence of string keys in map
func StrKeysInMap(mapObject interface{}, keys ...string) bool {
	switch typedMap := mapObject.(type) {
	case map[string]interface{}://, map[string]string:

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
	if len(value) == 32  {
		ok = true
		for i:=0; i<32; i++ {
			c := value[i]
			if !(c=='1' || c=='2' || c=='3' || c=='4' || c=='5' || c=='6' || c=='7' || c=='8' || c=='9' || c=='0' ||
			     c=='a' || c=='b' || c=='c' || c=='d' || c=='e' || c=='f') {
				ok = false
				break
			}
		}
	}
	return ok
}

// TODO: should be somwhere in other place
func GetSiteBackUrl() string {
	return "http://dev.ottemo.com:3000/"
}

// converts interface{} to string
func InterfaceToString(value interface{}) string {
	switch value := value.(type) {
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
		intValue, _ := StrToInt(typedValue)
		return intValue
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
	case string:
		floatValue, _ := strconv.ParseFloat(typedValue, 64)
		return floatValue
	default:
		return 0.0
	}
}

// convert string to integer
func StrToInt(value string) (int, error) {
	return strconv.Atoi(value)
}
