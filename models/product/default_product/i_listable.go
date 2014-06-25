package default_product

import (
	"github.com/ottemo/foundation/database"
	"errors"
)

//---------------------------------
// Implementation specific methods
//---------------------------------

var listCollection database.I_DBCollection = nil

func getListCollection() (database.I_DBCollection, error) {
	if listCollection != nil {
		var err error = nil

		dbEngine := database.GetDBEngine()
		if dbEngine == nil { return nil, errors.New("Can't obtain DBEngine") }

		listCollection, err = dbEngine.GetCollection("Product")

		return listCollection, err
	} else {
		return listCollection, nil
	}
}

//--------------------------
// Interface implementation
//--------------------------

func (it *DefaultProductModel) List() ([]interface{}, error) {
	result := make([]interface{}, 0)

	collection, err := getListCollection()
	if err != nil { return result, err }

	dbItems, err := collection.Load()
	if err != nil { return result, err }

	for _, dbItemData := range dbItems {
		model, err := it.New()
		if err != nil { return result, err }

		product := model.(*DefaultProductModel);
		product.FromHashMap(dbItemData)

		result = append(result, product)
	}

	return result, nil
}

func (it *DefaultProductModel) ListFilterAdd(Attribute string, Operator string, Value interface{}) error {
	collection, err := getListCollection()
	if err != nil { return err }

	collection.AddFilter(Attribute, Operator, Value.(string))
	return nil
}

func (it *DefaultProductModel) ListFilterReset() error {
	collection, err := getListCollection()
	if err != nil { return err }

	collection.ClearFilters()
	return nil
}
