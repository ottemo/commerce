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

	Street  string
	City    string
	State   string
	Phone   string
	ZipCode string

	listCollection     db.I_DBCollection
	listExtraAtributes []string
}
