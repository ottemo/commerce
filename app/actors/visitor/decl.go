package visitor

import (
	"github.com/ottemo/foundation/app/models/visitor"

	"github.com/ottemo/foundation/db"
)

const (
	VISITOR_COLLECTION_NAME = "visitor"
)

type DefaultVisitor struct {
	id string

	Email     string
	FirstName string
	LastName  string

	BillingAddress  visitor.I_VisitorAddress
	ShippingAddress visitor.I_VisitorAddress

	Password string
	Validated bool


	listCollection db.I_DBCollection
	listExtraAtributes []string
}
