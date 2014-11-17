// Package visitor is a default implementation of models/visitor package visitor related interfaces
package visitor

import (
	"github.com/ottemo/foundation/app/helpers/attributes"
	"github.com/ottemo/foundation/app/models/visitor"
	"github.com/ottemo/foundation/db"
	"time"
)

// Package global constants
const (
	COLLECTION_NAME_VISITOR = "visitor"

	EMAIL_VALIDATE_EXPIRE = 60 * 60 * 24
)

// DefaultVisitor is a default implementer of I_Visitor
type DefaultVisitor struct {
	id string

	Email      string
	FacebookId string
	GoogleId   string

	FirstName string
	LastName  string

	BillingAddress  visitor.I_VisitorAddress
	ShippingAddress visitor.I_VisitorAddress

	Password    string
	ValidateKey string

	Admin bool

	Birthday  time.Time
	CreatedAt time.Time

	*attributes.CustomAttributes
}

// DefaultVisitorCollection is a default implementer of I_VisitorCollection
type DefaultVisitorCollection struct {
	listCollection     db.I_DBCollection
	listExtraAtributes []string
}
