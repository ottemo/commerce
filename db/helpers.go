package db

import (
	"github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/utils"
	"strings"
)

// GetCollection returns database collection or error otherwise
func GetCollection(CollectionName string) (InterfaceDBCollection, error) {
	dbEngine := GetDBEngine()
	if dbEngine == nil {
		return nil, env.ErrorNew(ConstErrorModule, ConstErrorLevel, "7f379d36-eb93-4add-b1ee-df2c3a35a590", "Can't get DBEngine")
	}

	return dbEngine.GetCollection(CollectionName)
}

// ConvertTypeFromDbToGo returns object that represents GO side value for given valueType
func ConvertTypeFromDbToGo(value interface{}, valueType string) interface{} {
	switch {
	case strings.HasPrefix(valueType, "[]"):
		var result []interface{}
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

	case strings.HasPrefix(valueType, ConstTypeBoolean):
		return utils.InterfaceToBool(value)

	case strings.HasPrefix(valueType, ConstTypeInteger):
		return utils.InterfaceToInt(value)

	case strings.HasPrefix(valueType, ConstTypeDecimal),
		strings.HasPrefix(valueType, ConstTypeFloat),
		strings.HasPrefix(valueType, ConstTypeMoney):

		return utils.InterfaceToFloat64(value)

	case valueType == ConstTypeDatetime:
		return utils.InterfaceToTime(value)

	case valueType == ConstTypeJSON:
		result, _ := utils.DecodeJSONToStringKeyMap(value)
		return result

	case strings.HasPrefix(valueType, ConstTypeVarchar), valueType == ConstTypeText, valueType == ConstTypeID:
		return utils.InterfaceToString(value)

	}

	return value
}
