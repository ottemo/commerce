package product

import (
	"github.com/ottemo/foundation/db"
	"errors"
)

//---------------------------------
// Implementation specific methods
//---------------------------------

var listCollection db.I_DBCollection = nil

func getListCollection() (db.I_DBCollection, error) {
	if listCollection != nil {
		return listCollection, nil
	} else {
		var err error = nil

		dbEngine := db.GetDBEngine()
		if dbEngine == nil { return nil, errors.New("Can't obtain DBEngine") }

		listCollection, err = dbEngine.GetCollection("Product")

		return listCollection, err
	}
}

//--------------------------
// Interface implementation
//--------------------------

func (it *DefaultProduct) List() ([]interface{}, error) {
	result := make([]interface{}, 0)

	collection, err := getListCollection()
	if err != nil { return result, err }

	dbItems, err := collection.Load()
	if err != nil { return result, err }

	for _, dbItemData := range dbItems {
		model, err := it.New()
		if err != nil { return result, err }

		product := model.(*DefaultProduct);
		product.FromHashMap(dbItemData)

		result = append(result, product)
	}

	return result, nil
}

func (it *DefaultProduct) ListFilterAdd(Attribute string, Operator string, Value interface{}) error {
	collection, err := getListCollection()
	if err != nil { return err }

	collection.AddFilter(Attribute, Operator, Value.(string))
	return nil
}

func (it *DefaultProduct) ListFilterReset() error {
	collection, err := getListCollection()
	if err != nil { return err }

	collection.ClearFilters()
	return nil
}
