package address

import (
	"github.com/ottemo/foundation/db"
)

const (
	VISITOR_ADDRESS_COLLECTION_NAME = "visitor_address"
)

type DefaultVisitorAddress struct {
	id         string
	visitor_id string

	FirstName string
	LastName  string

	Company string

	Country string
	State   string
	City    string

	AddressLine1 string
	AddressLine2 string

	Phone   string
	ZipCode string

	listCollection     db.I_DBCollection
	listExtraAtributes []string
}
