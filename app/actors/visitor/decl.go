package visitor

import (
	"time"

	"github.com/ottemo/foundation/app/models/visitor"
	"github.com/ottemo/foundation/db"
)

const (
	COLLECTION_NAME_VISITOR = "visitor"

	EMAIL_VALIDATE_EXPIRE = 60 * 60 * 24
)

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

	IsAdmin bool

	Birthday  time.Time
	CreatedAt time.Time
}

type DefaultVisitorCollection struct {
	listCollection     db.I_DBCollection
	listExtraAtributes []string
}
