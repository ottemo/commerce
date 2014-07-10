package address

import (
	"errors"
	"github.com/ottemo/foundation/db"
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
		if dbEngine == nil {
			return nil, errors.New("Can't obtain DBEngine")
		}

		listCollection, err = dbEngine.GetCollection(VISITOR_ADDRESS_COLLECTION_NAME)

		return listCollection, err
	}
}

//--------------------------
// Interface implementation
//--------------------------

func (it *DefaultVisitorAddress) List() ([]interface{}, error) {
	result := make([]interface{}, 0)

	collection, err := getListCollection()
	if err != nil {
		return result, err
	}

	dbItems, err := collection.Load()
	if err != nil {
		return result, err
	}

	for _, dbItemData := range dbItems {
		model, err := it.New()
		if err != nil {
			return result, err
		}

		address := model.(*DefaultVisitorAddress)
		address.FromHashMap(dbItemData)

		result = append(result, address)
	}

	return result, nil
}

func (it *DefaultVisitorAddress) ListFilterAdd(Attribute string, Operator string, Value interface{}) error {
	collection, err := getListCollection()
	if err != nil {
		return err
	}

	collection.AddFilter(Attribute, Operator, Value.(string))
	return nil
}

func (it *DefaultVisitorAddress) ListFilterReset() error {
	collection, err := getListCollection()
	if err != nil {
		return err
	}

	collection.ClearFilters()
	return nil
}
