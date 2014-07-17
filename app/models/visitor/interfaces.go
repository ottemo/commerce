package visitor

import (
	"github.com/ottemo/foundation/app/models"
)

type I_VisitorAddress interface {
	GetVisitorId() string
	GetStreet() string
	GetCity() string
	GetState() string
	GetPhone() string
	GetZipCode() string

	models.I_Model
	models.I_Object
	models.I_Storable
	models.I_Listable
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

	SetPassword(passwd string) error
	CheckPassword(passwd string) bool

	IsValidated() bool
	Invalidate() error
	Validate(key string) error

	models.I_Model
	models.I_Object
	models.I_Storable
	models.I_Listable
}
