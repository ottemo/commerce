package attributes

import (
	"errors"
	"github.com/ottemo/foundation/app/models"
)

func (it *CustomAttributes) Get(attribute string) interface{} {
	return it.values[attribute]
}

func (it *CustomAttributes) Set(attribute string, value interface{}) error {
	if _, present := it.attributes[attribute]; present {
		it.values[attribute] = value
	} else {
		return errors.New("attribute '" + attribute + "' invalid")
	}

	return nil
}

func (it *CustomAttributes) GetAttributesInfo() []models.T_AttributeInfo {
	info := make([]models.T_AttributeInfo, 0)
	for _, attribute := range it.attributes {
		info = append(info, attribute)
	}
	return info
}