package order

import (
	"github.com/ottemo/foundation/db"
)

// returns database collection
func (it *DefaultOrderItemCollection) GetDBCollection() db.I_DBCollection {
	return it.listCollection
}
