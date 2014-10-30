package db

import (
	"github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/utils"
	"strings"
)

// returns database collection or error otherwise
func GetCollection(CollectionName string) (I_DBCollection, error) {
	dbEngine := GetDBEngine()
	if dbEngine == nil {
		return nil, env.ErrorNew("Can't get DBEngine")
	}

	return dbEngine.GetCollection(CollectionName)
}

// returns object that represents GO side value for given valueType
func ConvertTypeFromDbToGo(value interface{}, valueType string) interface{} {
	switch {
	case strings.HasPrefix(valueType, "[]"):
		result := make([]interface{}, 0)
		if value == nil {
			return result
		}

		valueString := utils.InterfaceToString(value)
		if valueString == "" {
			return result
		}

		array := strings.Split(valueString, ", ")
		arrayType := strings.TrimPrefix(valueType, "[]")

		for _, arrayValue := range array {
			result = append(result, ConvertTypeFromDbToGo(arrayValue, arrayType))
		}
		return result

	case strings.HasPrefix(valueType, DB_BASETYPE_BOOLEAN):
		return utils.InterfaceToBool(value)

	case strings.HasPrefix(valueType, DB_BASETYPE_INTEGER):
		return utils.InterfaceToInt(value)

	case strings.HasPrefix(valueType, DB_BASETYPE_DECIMAL),
		strings.HasPrefix(valueType, DB_BASETYPE_FLOAT),
		strings.HasPrefix(valueType, DB_BASETYPE_MONEY):

		return utils.InterfaceToFloat64(value)

	case valueType == DB_BASETYPE_DATETIME:
		return utils.InterfaceToTime(value)

	case valueType == DB_BASETYPE_JSON:
		result, _ := utils.DecodeJsonToStringKeyMap(value)
		return result

	case strings.HasPrefix(valueType, DB_BASETYPE_VARCHAR), valueType == DB_BASETYPE_TEXT, valueType == DB_BASETYPE_ID:
		return utils.InterfaceToString(value)

	}

	return value
}
