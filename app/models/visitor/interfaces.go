package visitor

import (
	"github.com/ottemo/foundation/app/models"
)

type I_VisitorAddress interface {
	GetStreet() string
	GetCity() string
	GetState() string
	GetPhone() string
	GetZipCode() string

	models.I_Model
	models.I_Object
	models.I_Storable
}

type I_Visitor interface {
	GetEmail() string

	GetFullName() string
	GetFirstName() string
	GetLastName() string

	GetShippingAddress() I_VisitorAddress
	GetBillingAddress() I_VisitorAddress

	SetShippingAddress(address I_VisitorAddress) error
	SetBillingAddress(address I_VisitorAddress) error

	models.I_Model
	models.I_Object
	models.I_Storable
}
