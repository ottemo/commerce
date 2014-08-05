package visitor

import (
	"github.com/ottemo/foundation/app/models"
)

const (
	VISITOR_MODEL_NAME = "Visitor"
	VISITOR_ADDRESS_MODEL_NAME = "VisitorAddress"
	SESSION_KEY_VISITOR_ID = "visitor_id"
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
	GetFacebookId() string
	GetGoogleId() string

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

	LoadByEmail(email string) error
	LoadByFacebookId(facebookId string) error
	LoadByGoogleId(googleId string) error

	models.I_Model
	models.I_Object
	models.I_Storable
	models.I_Listable
}
