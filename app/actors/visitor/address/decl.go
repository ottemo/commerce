package address

import (
	"github.com/ottemo/foundation/db"
)

// Constants for the Visistor Address collection
const (
	COLLECTION_NAME_VISITOR_ADDRESS = "visitor_address"
)

// DefaultVisitorAddress is a struct which holds the default information representing a Visitor Address
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
}

// DefaultVisitorAddressCollection is a struct which holds the collection information for the Visitor Address
type DefaultVisitorAddressCollection struct {
	listCollection     db.I_DBCollection
	listExtraAtributes []string
}
