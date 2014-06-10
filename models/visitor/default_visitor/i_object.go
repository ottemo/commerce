package default_visitor

import (
	"strings"
	"github.com/ottemo/foundation/models"

	"github.com/ottemo/foundation/models/visitor"
	"errors"
)

func (it *DefaultVisitor) GetId() bool {
	return it.id
}

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
	case "billing_address":
		return it.GetBillingAddress()
	case "shipping_address":
		return it.GetShippingAddress()
	}

	return nil
}

func (it *DefaultVisitor) Set(attribute string, value interface{}) error {
	attribute := strings.ToLower(attribute)

	switch  attribute {
	case "_id", "id":
		it.id = value.(string)
	case "email", "e_mail", "e-mail":
		it.Email = value.(string)
	case "fname", "first_name":
		it.Fname = value.(string)
	case "lname", "last_name":
		it.Lname = value.(string)

	//case "billing_address_id":
	//	it.BillingAddressId = value.(string)
	//case "shipping_address_id":
	//	it.ShippingAddressId = value.(string)

	case "billing_address", "shipping_address":
		switch value := value.(type) {
		case visitor.I_VisitorAddress:
			if attribute == "billing_address" {
				it.SetBillingAddress(value)
			} else {
				it.SetShippingAddress(value)
			}

		case map[string]interface{}:
			model, err := models.GetModel("VisitorAddress")
			if err != nil { return err }

			if address, ok := model.(visitor.I_VisitorAddress); ok {
				err := address.FromHashMap(value)
				if err != nil { return err }

				if attribute == "billing_address" {
					it.SetBillingAddress(address)
				} else {
					it.SetShippingAddress(address)
				}
			} else {
				errors.New("unsupported visitor addres model " + model.GetImplementationName())
			}

		default:
			return errors.New("unsupported 'billing_address' value")
		}
	}
	return nil
}

func (it *DefaultVisitor) ListAttributes() []models.T_AttributeInfo {
	return make([]models.T_AttributeInfo, 0)
}
