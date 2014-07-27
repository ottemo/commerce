package listable

import (
	"errors"
	"github.com/ottemo/foundation/db"

	"github.com/ottemo/foundation/app/models"
	"github.com/ottemo/foundation/app/utils"
)

// initializes and returns shared among functions collection
func (it *ListableHelper) getListCollection() (db.I_DBCollection, error) {
	if it.listCollection != nil {
		return it.listCollection, nil
	} else {
		var err error = nil

		if it.delegate.GetCollection != nil {
			it.listCollection = it.delegate.GetCollection()
		} else {
			dbEngine := db.GetDBEngine()
			if dbEngine == nil {
				return nil, errors.New("Can't obtain DBEngine")
			}

			it.listCollection, err = dbEngine.GetCollection(it.delegate.CollectionName)
		}

		return it.listCollection, err
	}
}

//--------------------------
// INTERFACE IMPLEMENTATION
//--------------------------

// returns count of items inside list
func (it *ListableHelper) ListCount() (int, error) {
	// loading data from DB
	collection, err := it.getListCollection()
	if err != nil {
		return 0, err
	}

	return collection.Count()
}

// enumerates items of model type
func (it *ListableHelper) ListObjects() ([]interface{}, error) {
	result := make([]interface{}, 0)

	// loading data from DB
	collection, err := it.getListCollection()
	if err != nil {
		return result, err
	}

	dbItems, err := collection.Load()
	if err != nil {
		return result, err
	}

	// transforming record data to object
	for _, dbItemData := range dbItems {
		if it.delegate.RecordToObjectFunc != nil {
			result = append(result, it.delegate.RecordToObjectFunc(dbItemData, it.listExtraAtributes) )
		} else {
			result = append(result, dbItemData)
		}
	}

	return result, nil
}

// enumerates items of model type in T_ListItem format
func (it *ListableHelper) List() ([]models.T_ListItem, error) {
	result := make([]models.T_ListItem, 0)

	// loading data from DB
	collection, err := it.getListCollection()
	if err != nil {
		return result, err
	}

	dbItems, err := collection.Load()
	if err != nil {
		return result, err
	}

	// transforming record data to ListItem
	for _, dbItemData := range dbItems {

		if it.delegate.RecordToListItemFunc != nil {
			result = append(result, it.delegate.RecordToListItemFunc(dbItemData, it.listExtraAtributes) )
		} else {
			resultItem := new(models.T_ListItem)

			if value, present := dbItemData["id"]; present {
				resultItem.Id, _ = value.(string)
			}
			if value, present := dbItemData["name"]; present {
				resultItem.Name, _ = value.(string)
			}
			if value, present := dbItemData["image"]; present {
				resultItem.Image, _ = value.(string)
			}
			if value, present := dbItemData["desc"]; present {
				resultItem.Desc, _ = value.(string)
			}

			result = append(result, *resultItem)
		}
	}

	return result, nil
}



// allows to obtain additional attributes from  List() function
func (it *ListableHelper) ListAddExtraAttribute(attribute string) error {

	if it.delegate.ValidateExtraAttributeFunc == nil || it.delegate.ValidateExtraAttributeFunc(attribute) {
		if utils.IsInListStr(attribute, it.listExtraAtributes) {
			it.listExtraAtributes = append(it.listExtraAtributes, attribute)
		} else {
			return errors.New("attribute already in list")
		}
	} else {
		return errors.New("not allowed list attribute " + attribute)
	}

	return nil
}


// adds selection filter to List() function
func (it *ListableHelper) ListFilterAdd(Attribute string, Operator string, Value interface{}) error {
	collection, err := it.getListCollection()
	if err != nil {
		return err
	}

	return collection.AddFilter(Attribute, Operator, Value.(string))
}



// clears presets made by ListFilterAdd() and ListAddExtraAttribute() functions
func (it *ListableHelper) ListFilterReset() error {
	collection, err := it.getListCollection()
	if err != nil {
		return err
	}

	collection.ClearFilters()

	it.listExtraAtributes = make([]string, 0)

	return nil
}



// specifies selection paging
func (it *ListableHelper) ListLimit(offset int, limit int) error {
	collection, err := it.getListCollection()
	if err != nil {
		return err
	}

	return collection.SetLimit(offset, limit)
}
