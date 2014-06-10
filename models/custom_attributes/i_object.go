package custom_attributes

import (
	"errors"
	"github.com/ottemo/foundation/models"
)

func (it *CustomAttributes) GetId() bool {
	return it.id
}

func (it *CustomAttributes) Has(attribute string) bool {
	_, present := it.attributes[attribute]
	return present
}

func (it *CustomAttributes) Get(attribute string) interface{} {
	return it.values[attribute]
}

func (it *CustomAttributes) Set(attribute string, value interface{}) error {
	if it.Has(attribute) {
		it.values[attribute] = value
	} else {
		return errors.New("attribute '" + attribute + "' invalid")
	}

	return nil
}

func (it *CustomAttributes) ListAttributes() []models.T_AttributeInfo {
	returnValue := make([]models.T_AttributeInfo, len(it.attributes))
	for _, attribute := range it.attributes {
		returnValue = append(returnValue, attribute)
	}
	return returnValue
}
