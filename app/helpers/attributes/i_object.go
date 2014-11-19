package attributes

import (
	"github.com/ottemo/foundation/app/models"
	"github.com/ottemo/foundation/env"
)

// Get returns object attribute value or nil
func (it *CustomAttributes) Get(attribute string) interface{} {
	return it.values[attribute]
}

// Set sets attribute value to object or returns error
func (it *CustomAttributes) Set(attribute string, value interface{}) error {
	if _, present := it.attributes[attribute]; present {
		it.values[attribute] = value
	} else {
		return env.ErrorNew("attribute '" + attribute + "' invalid")
	}

	return nil
}

// GetAttributesInfo represents object as map[string]interface{}
func (it *CustomAttributes) GetAttributesInfo() []models.StructAttributeInfo {
	info := make([]models.StructAttributeInfo, 0)
	for _, attribute := range it.attributes {
		info = append(info, attribute)
	}
	return info
}

// FromHashMap represents object as map[string]interface{}
func (it *CustomAttributes) FromHashMap(input map[string]interface{}) error {
	it.values = input
	return nil
}

// ToHashMap fills object attributes from map[string]interface{}
func (it *CustomAttributes) ToHashMap() map[string]interface{} {
	return it.values
}
