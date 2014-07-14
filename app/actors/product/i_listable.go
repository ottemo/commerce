package product

import (
	"errors"
	"github.com/ottemo/foundation/db"

	"github.com/ottemo/foundation/app/models"
	"github.com/ottemo/foundation/app/utils"
)

//---------------------------------
// IMPLEMENTATION SPECIFIC METHODS
//---------------------------------

// initializes and returns shared among couple functions collection
func (it *DefaultProduct) getListCollection() (db.I_DBCollection, error) {
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
func (it *DefaultProduct) List() ([]models.T_ListItem, error) {
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

		// assigning loaded DB data to model
		model, err := it.New()
		if err != nil {
			return result, err
		}

		product := model.(*DefaultProduct)
		product.FromHashMap(dbItemData)

		// retrieving minimal data needed for list
		resultItem := new(models.T_ListItem)

		mediaPath, err := product.GetMediaPath("image")
		if err != nil {
			return result, err
		}

		resultItem.Id    = product.GetId()
		resultItem.Name  = "[" + product.GetSku() + "] "  + product.GetName()
		resultItem.Image = ""
		resultItem.Desc  = product.GetDescription()

		if product.GetDefaultImage() != "" {
			resultItem.Image = mediaPath + product.GetDefaultImage()
		}

		// if extra attributes were required
		if len(it.listExtraAtributes) > 0 {
			resultItem.Extra = make(map[string]interface{})

			for _, attributeName := range it.listExtraAtributes {
				resultItem.Extra[attributeName] = product.Get(attributeName)
			}
		}

		result = append(result, *resultItem)
	}

	return result, nil
}



// allows to obtain additional attributes from  List() function
func (it *DefaultProduct) ListAddExtraAttribute(attribute string) error {

	if utils.IsAmongStr(attribute, "sku", "name", "description", "price", "default_image") {
		if utils.IsInListStr(attribute, it.listExtraAtributes) {
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
func (it *DefaultProduct) ListFilterAdd(Attribute string, Operator string, Value interface{}) error {
	collection, err := it.getListCollection()
	if err != nil {
		return err
	}

	collection.AddFilter(Attribute, Operator, Value.(string))
	return nil
}



// clears presets made by ListFilterAdd() and ListAddExtraAttribute() functions
func (it *DefaultProduct) ListFilterReset() error {
	collection, err := it.getListCollection()
	if err != nil {
		return err
	}

	collection.ClearFilters()

	it.listExtraAtributes = make([]string, 0)

	return nil
}



// specifies selection paging
func (it *DefaultProduct) ListLimit(offset int, limit int) error {
	collection, err := it.getListCollection()
	if err != nil {
		return err
	}

	return collection.SetLimit(offset, limit)
}
