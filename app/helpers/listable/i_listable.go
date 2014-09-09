package listable

import (
	"errors"

	"github.com/ottemo/foundation/app/models"
	"github.com/ottemo/foundation/app/utils"
)

// enumerates items of model type
func (it *ListableHelper) ListObjects() ([]interface{}, error) {
	result := make([]interface{}, 0)

	// loading data from DB
	collection := it.GetDBCollection()
	if collection == nil {
		return result, errors.New("can't obtain collection")
	}

	dbItems, err := collection.Load()
	if err != nil {
		return result, err
	}

	// transforming record data to object
	for _, dbItemData := range dbItems {
		if it.delegate.RecordToObjectFunc != nil {
			result = append(result, it.delegate.RecordToObjectFunc(dbItemData, it.listExtraAtributes))
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
	collection := it.GetDBCollection()
	if collection == nil {
		return result, errors.New("can't obtain collection")
	}

	dbItems, err := collection.Load()
	if err != nil {
		return result, err
	}

	// transforming record data to ListItem
	for _, dbItemData := range dbItems {

		if it.delegate.RecordToListItemFunc != nil {
			item, ok := it.delegate.RecordToListItemFunc(dbItemData, it.listExtraAtributes)
			if ok {
				result = append(result, item)
			}
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
		if !utils.IsInListStr(attribute, it.listExtraAtributes) {
			it.listExtraAtributes = append(it.listExtraAtributes, attribute)
		} else {
			return errors.New("attribute '" + attribute + "' already in list")
		}
	} else {
		return errors.New("not allowed list attribute " + attribute)
	}

	return nil
}

// adds selection filter to List() function
func (it *ListableHelper) ListFilterAdd(Attribute string, Operator string, Value interface{}) error {
	collection := it.GetDBCollection()
	if collection != nil {
		return collection.AddFilter(Attribute, Operator, Value.(string))
	}

	return errors.New("can't obtain collection")
}

// clears presets made by ListFilterAdd() and ListAddExtraAttribute() functions
func (it *ListableHelper) ListFilterReset() error {
	collection := it.GetDBCollection()
	if collection != nil {
		collection.ClearFilters()
		it.listExtraAtributes = make([]string, 0)
	}

	return nil
}

// specifies selection paging
func (it *ListableHelper) ListLimit(offset int, limit int) error {
	collection := it.GetDBCollection()
	if collection != nil {
		return collection.SetLimit(offset, limit)
	}

	return errors.New("can't obtain collection")
}
