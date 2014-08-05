package objectref

import (
	"errors"
	"strings"

	"github.com/ottemo/foundation/app/models"
	"github.com/ottemo/foundation/app/models/visitor"
	"github.com/ottemo/foundation/app/utils"
)

// returns attribute value for current object or nil if no such attribute
func (it *DBObjectRef) Get(attribute string) interface{} {
	var result interface{} = nil

	if it.currData != nil {
		result, _ = it.currData[attribute]
	}

	return result
}

// sets attribute value for current object
func (it *DBObjectRef) Set(attribute string, value interface{}) error {
	if it.currData == nil {
		it.currData = new(map[string]interface{})
	}

	it.currData[attribute] = value

	return nil
}

// fills attributes values based on provided map
func (it *DBObjectRef) FromHashMap(input map[string]interface{}) error {

	if it.currData == nil {
		it.currData = new(map[string]interface{})
	}

	for attribute, value := range input {
		if err := it.Set(attribute, value); err != nil {
			return err
		}
	}

	return nil
}

// returns attribute values that current object holds
func (it *DBObjectRef) ToHashMap() map[string]interface{} {

	result := new(map[string]interface{})

	if it.currData != nil {
		for attribute, value := range it.currData {
			result[attribute] = value
		}
	}

	return result
}

// returns stub information about current object attributes
//   - if you using this helper you should rewrite this function in your class
func (it *DBObjectRef) GetAttributesInfo() []models.T_AttributeInfo {

	result := new([]models.T_AttributeInfo)

	if it.currData != nil {
		for attribute, _ := range it.currData {
			result = append(result,
				models.T_AttributeInfo{
					Model:      "",
					Collection: "",
					Attribute:  attribute,
					Type:       "",
					IsRequired: false,
					IsStatic:   true,
					Label:      attribute,
					Group:      "General",
					Editors:    "not_editable",
					Options:    "",
					Default:    "",
				})
		}
	}

	return result
}
