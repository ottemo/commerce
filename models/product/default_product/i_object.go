package default_product

import (
	"strings"
	"github.com/ottemo/foundation/models"
)

func (it *DefaultProductModel) Has(attribute string) bool {
	return it.Get(attribute) == nil
}

func (it *DefaultProductModel) Get(attribute string) interface{} {
	switch strings.ToLower(attribute) {
	case "_id", "id":
		return it.id
	case "sku":
		return it.Sku
	case "name":
		return it.Name
	default:
		return it.CustomAttributes.Get(attribute)
	}

	return nil
}

func (it *DefaultProductModel) Set(attribute string, value interface{}) error {
	switch strings.ToLower(attribute) {
	case "_id", "id":
		it.id = value.(string)
	case "sku":
		it.Sku = value.(string)
	case "name":
		it.Name = value.(string)
	default:
		if err := it.CustomAttributes.Set(attribute, value); err != nil {
			return err
		}

	}

	return nil
}

func (it *DefaultProductModel) ListAttributes() []models.T_AttributeInfo {
	return make([]models.T_AttributeInfo, 0)
}
