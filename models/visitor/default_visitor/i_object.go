package default_visitor

import (
	"strings"
	"github.com/ottemo/foundation/models"
)

func (it *DefaultVisitor) Has(attribute string) bool {
	return it.Get(attribute) == nil
}

func (it *DefaultVisitor) Get(attribute string) interface{} {
	switch strings.ToLower(attribute) {
	case "_id", "id":
		return it.id
	case "fname", "first_name":
		return it.Fname
	case "lname", "last_name":
		return it.Lname
	case "billing_address_id":
		return it.BillingAddressId
	case "shipping_address_id":
		return it.ShippingAddressId
	}

	return nil
}

func (it *DefaultVisitor) Set(attribute string, value interface{}) error {
	switch strings.ToLower(attribute) {
	case "_id", "id":
		it.id = value.(string)
	case "email", "e_mail", "e-mail":
		it.Email = value.(string)
	case "fname", "first_name":
		it.Fname = value.(string)
	case "lname", "last_name":
		it.Lname = value.(string)
	case "billing_address_id":
		it.BillingAddressId = value.(string)
	case "shipping_address_id":
		it.ShippingAddressId = value.(string)
	}

	return nil
}

func (it *DefaultVisitor) ListAttributes() []models.T_AttributeInfo {
	return make([]models.T_AttributeInfo, 0)
}
