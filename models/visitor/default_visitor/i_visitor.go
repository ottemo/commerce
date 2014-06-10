package default_visitor

import (
	"github.com/ottemo/foundation/models/visitor"
	"github.com/ottemo/foundation/models"
)

func (it *DefaultVisitor) GetEmail() string { return it.Email }

func (it *DefaultVisitor) GetFullName() string { return it.Fname + " " + it.Lname }
func (it *DefaultVisitor) GetFirstName() string { return it.Fname }
func (it *DefaultVisitor) GetLastName() string { return it.Lname }



func getVisitorAddressById(addressId string) visitor.I_VisitorAddress {
	address_model, err := models.GetModel("VisitorAddress")
	if err != nil { return nil }

	if address_model, ok := address_model.(visitor.I_VisitorAddress); ok {
		address_model.Load( addressId )

		return address_model
	}

	return nil
}

func (it *DefaultVisitor) GetShippingAddress() visitor.I_VisitorAddress {
	return getVisitorAddressById(it.ShippingAddressId)
}

func (it *DefaultVisitor) GetBillingAddress() visitor.I_VisitorAddress {
	return getVisitorAddressById(it.BillingAddressId)
}
