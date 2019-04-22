package order

import (
	"github.com/ottemo/commerce/db"
)

// GetDBCollection returns database collection
func (it *DefaultOrderItemCollection) GetDBCollection() db.InterfaceDBCollection {
	return it.listCollection
}
