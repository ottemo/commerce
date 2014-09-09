package listable

import (
	"github.com/ottemo/foundation/db"
)

// returns database collection
func (it *ListableHelper) GetDBCollection() db.I_DBCollection {
	if it.listCollection != nil {
		return it.listCollection
	} else {
		if it.delegate.GetCollection != nil {
			it.listCollection = it.delegate.GetCollection()
		} else {
			it.listCollection, _ = db.GetCollection(it.delegate.CollectionName)
		}

		return it.listCollection
	}
}
