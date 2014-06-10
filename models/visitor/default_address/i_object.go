package default_address

import (
	"strings"
	"github.com/ottemo/foundation/models"
)

func (it *DefaultVisitorAddress) Has(attribute string) bool {
	return it.Get(attribute) == nil
}

func (it *DefaultVisitorAddress) Get(attribute string) interface{} {
	switch strings.ToLower(attribute) {
	case "_id", "id":
		return it.id
	case "street":
		return it.Street
	case "city":
		return it.City
	case "state":
		return it.State
	case "phone":
		return it.Phone
	case "zip", "zip_code":
		return it.ZipCode
	}

	return nil
}

func (it *DefaultVisitorAddress) Set(attribute string, value interface{}) error {
	switch strings.ToLower(attribute) {
	case "_id", "id":
		it.id = value.(string)
	case "street":
		it.Street = value.(string)
	case "city":
		it.City = value.(string)
	case "state":
		it.State = value.(string)
	case "phone":
		it.Phone = value.(string)
	case "zip", "zip_code":
		it.ZipCode = value.(string)
	}
	return nil
}

func (it *DefaultVisitorAddress) ListAttributes() []models.T_AttributeInfo {
	return make([]models.T_AttributeInfo, 0)
}
