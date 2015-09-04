package token

import (
	"github.com/ottemo/foundation/app/models"
	"github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/utils"
	"strings"
)

// Get will return the requested attribute when provided a string representation of the attribute
func (it *DefaultVisitorCard) Get(attribute string) interface{} {
	switch strings.ToLower(attribute) {
	case "_id", "id":
		return it.id
	case "visitor_id", "visitorID":
		return it.visitorID
	case "holder", "holder_name":
		return it.Holder
	case "payment", "payment_code":
		return it.Payment
	case "type":
		return it.Type
	case "number":
		return it.Number
	case "expiration_date", "expiration":
		return it.ExpirationDate
	case "token":
		return it.Token
	}

	return nil
}

// Set will set a Visitor Token attribute and requiring a name and a value
func (it *DefaultVisitorCard) Set(attribute string, value interface{}) error {
	switch strings.ToLower(attribute) {
	case "_id", "id":
		it.id = utils.InterfaceToString(value)

	case "visitor_id", "visitorID":
		it.visitorID = utils.InterfaceToString(value)

	case "fname", "first_name":
		it.Holder = utils.InterfaceToString(value)

	case "payment", "payment_code":
		it.Payment = utils.InterfaceToString(value)

	case "type":
		it.Type = utils.InterfaceToString(value)

	case "number":
		it.Number = utils.InterfaceToString(value)

	case "expiration_date", "expiration":
		it.ExpirationDate = utils.InterfaceToString(value)

	case "token":
		it.Token = utils.InterfaceToString(value)
	}
	return nil
}

// FromHashMap will take a map[string]interface and apply the attribute values to the Visitor Token
func (it *DefaultVisitorCard) FromHashMap(input map[string]interface{}) error {

	for attribute, value := range input {
		if err := it.Set(attribute, value); err != nil {
			env.ErrorDispatch(err)
		}
	}

	return nil
}

// ToHashMap will return a set of Visitor Token attributes in a map[string]interface
func (it *DefaultVisitorCard) ToHashMap() map[string]interface{} {

	result := make(map[string]interface{})

	result["_id"] = it.id

	result["visitor_id"] = it.visitorID

	result["holder"] = it.Holder
	result["payment"] = it.Payment
	result["type"] = it.Type
	result["number"] = it.Number
	result["expirationDate"] = it.ExpirationDate
	result["token"] = it.Token

	return result
}

// GetAttributesInfo will return a set of Vistor Token attributes in []models.StructAttributeInfo
func (it *DefaultVisitorCard) GetAttributesInfo() []models.StructAttributeInfo {
	info := []models.StructAttributeInfo{}

	return info
}
