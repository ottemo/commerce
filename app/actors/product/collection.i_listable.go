package product

import (
	"github.com/ottemo/foundation/app/models"
	"github.com/ottemo/foundation/app/models/product"
)

// internal delegate function for ListableHelper, converts db record to T_ListItem
func listableRecordToListItemFunc(recordData map[string]interface{}, extraAttributes []string) (models.T_ListItem, bool) {
	result := models.T_ListItem{}

	productModel, err := product.GetProductModel()
	if err != nil {
		return result, false
	}

	err = productModel.FromHashMap(recordData)
	if err != nil {
		return result, false
	}

	result.Id = productModel.GetId()
	result.Name = "[" + productModel.GetSku() + "] " + productModel.GetName()
	result.Image = ""
	result.Desc = productModel.GetShortDescription()

	if productModel.GetDefaultImage() != "" {
		mediaPath, err := productModel.GetMediaPath("image")
		if err == nil {
			result.Image = mediaPath + productModel.GetDefaultImage()
		}
	}

	// serving extra attributes
	//-------------------------
	if len(extraAttributes) > 0 {
		result.Extra = make(map[string]interface{})

		for _, attributeName := range extraAttributes {
			result.Extra[attributeName] = productModel.Get(attributeName)
		}
	}

	return result, true
}

/*
import (
	"errors"
	"github.com/ottemo/foundation/db"

	"github.com/ottemo/foundation/app/models"
	"github.com/ottemo/foundation/app/utils"

	"github.com/ottemo/foundation/app/models/product"
)

//---------------------------------
// IMPLEMENTATION SPECIFIC METHODS
//---------------------------------

// initializes and returns shared among couple functions collection
func (it *DefaultProductCollection) getListCollection() (db.I_DBCollection, error) {
	if it.listCollection != nil {
		return it.listCollection, nil
	} else {
		var err error = nil

		dbEngine := db.GetDBEngine()
		if dbEngine == nil {
			return nil, errors.New("Can't obtain DBEngine")
		}

		it.listCollection, err = dbEngine.GetCollection("Product")

		return it.listCollection, err
	}
}

//--------------------------
// INTERFACE IMPLEMENTATION
//--------------------------

// enumerates items of Product model type
func (it *DefaultProductCollection) List() ([]models.T_ListItem, error) {
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

	for _, dbItemData := range dbItems {

		productModel, err := product.GetProductModel()
		if err != nil {
			return result, err
		}
		productModel.FromHashMap(dbItemData)

		// retrieving minimal data needed for list
		resultItem := new(models.T_ListItem)

		mediaPath, err := productModel.GetMediaPath("image")
		if err != nil {
			return result, err
		}

		resultItem.Id = productModel.GetId()
		resultItem.Name = "[" + productModel.GetSku() + "] " + productModel.GetName()
		resultItem.Image = ""
		resultItem.Desc = productModel.GetShortDescription()

		if productModel.GetDefaultImage() != "" {
			resultItem.Image = mediaPath + productModel.GetDefaultImage()
		}

		// if extra attributes were required
		if len(it.listExtraAtributes) > 0 {
			resultItem.Extra = make(map[string]interface{})

			for _, attributeName := range it.listExtraAtributes {
				resultItem.Extra[attributeName] = productModel.Get(attributeName)
			}
		}

		result = append(result, *resultItem)
	}

	return result, nil
}

// allows to obtain additional attributes from  List() function
func (it *DefaultProductCollection) ListAddExtraAttribute(attribute string) error {

	if utils.IsAmongStr(attribute, "sku", "name", "description", "price", "default_image") {
		if !utils.IsInListStr(attribute, it.listExtraAtributes) {
			it.listExtraAtributes = append(it.listExtraAtributes, attribute)
		} else {
			return errors.New("attribute already in list")
		}
	} else {
		return errors.New("not allowed attribute")
	}

	return nil
}

// adds selection filter to List() function
func (it *DefaultProductCollection) ListFilterAdd(Attribute string, Operator string, Value interface{}) error {
	collection, err := it.getListCollection()
	if err != nil {
		return err
	}

	collection.AddFilter(Attribute, Operator, Value.(string))
	return nil
}

// clears presets made by ListFilterAdd() and ListAddExtraAttribute() functions
func (it *DefaultProductCollection) ListFilterReset() error {
	collection, err := it.getListCollection()
	if err != nil {
		return err
	}

	collection.ClearFilters()

	it.listExtraAtributes = make([]string, 0)

	return nil
}

// specifies selection paging
func (it *DefaultProductCollection) ListLimit(offset int, limit int) error {
	collection, err := it.getListCollection()
	if err != nil {
		return err
	}

	return collection.SetLimit(offset, limit)
}
*/
