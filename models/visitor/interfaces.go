package visitor

import("github.com/ottemo/foundation/models")

type IVisitorAddress interface {
	GetStreet() string
	GetCity() string
	GetState() string
	GetPhone() string
	GetZipCode() string

	models.IModel
	models.IObject
	models.IStorable
	models.IMapable
}

type IVisitor interface {
	GetEmail() string

	GetFullName() string
	GetFirstName() string
	GetLastName() string

	GetShippingAddress() IVisitorAddress
	GetBillingAddress() IVisitorAddress

	SetShippingAddress(address IVisitorAddress) error
	SetBillingAddress(address IVisitorAddress) error

	models.IModel
	models.IObject
	models.IStorable
	models.IMapable
}
