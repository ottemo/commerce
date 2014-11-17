// Package address is a default implementation of models/visitor package visitor address related  interfaces
package address

import (
	"github.com/ottemo/foundation/db"
)

// Package global constants
const (
	COLLECTION_NAME_VISITOR_ADDRESS = "visitor_address"
)

// DefaultVisitorAddress is a default implementer of I_VisitorAddress
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

// DefaultVisitorAddressCollection is a default implementer of I_VisitorAddressCollection
type DefaultVisitorAddressCollection struct {
	listCollection     db.I_DBCollection
	listExtraAtributes []string
}
