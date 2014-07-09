package category

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

		listCollection, err = dbEngine.GetCollection("Category")

		return listCollection, err
	}
}

//------------------------------------
// I_Listable interface implementation
//------------------------------------

// returns category list array with usage of filtering if set
func (it *DefaultCategory) List() ([]interface{}, error) {
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

		category := model.(*DefaultCategory)
		category.FromHashMap(dbItemData)

		result = append(result, category)
	}

	return result, nil
}

// adds selection filter to List() function
func (it *DefaultCategory) ListFilterAdd(Attribute string, Operator string, Value interface{}) error {
	collection, err := getListCollection()
	if err != nil {
		return err
	}

	collection.AddFilter(Attribute, Operator, Value.(string))
	return nil
}

// resets all previously set filters
func (it *DefaultCategory) ListFilterReset() error {
	collection, err := getListCollection()
	if err != nil {
		return err
	}

	collection.ClearFilters()
	return nil
}
