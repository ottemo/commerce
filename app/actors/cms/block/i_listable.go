package block

import (
	"errors"
	"github.com/ottemo/foundation/db"

	"github.com/ottemo/foundation/app/models"
	"github.com/ottemo/foundation/app/utils"
)

// Internal functions
//-------------------

// initializes and/or returns shared among interface functions collection
func (it *DefaultCMSBlock) getListCollection() (db.I_DBCollection, error) {
	if it.listCollection != nil {
		return it.listCollection, nil
	} else {
		var err error = nil
		it.listCollection, err = db.GetCollection(CMS_BLOCK_COLLECTION_NAME)

		return it.listCollection, err
	}
}

// Interface implementation
//-------------------------

// enumerates items of Product model type
func (it *DefaultCMSBlock) List() ([]models.T_ListItem, error) {
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

		cmsBlock := model.(*DefaultCMSBlock)
		cmsBlock.FromHashMap(dbItemData)

		// retrieving minimal data needed for list
		resultItem := new(models.T_ListItem)

		resultItem.Id = cmsBlock.GetId()
		resultItem.Name = cmsBlock.GetIdentifier()
		resultItem.Image = ""
		resultItem.Desc = ""

		// if extra attributes were required
		if len(it.listExtraAtributes) > 0 {
			resultItem.Extra = make(map[string]interface{})

			for _, attributeName := range it.listExtraAtributes {
				resultItem.Extra[attributeName] = cmsBlock.Get(attributeName)
			}
		}

		result = append(result, *resultItem)
	}

	return result, nil
}

// allows to obtain additional attributes from  List() function
func (it *DefaultCMSBlock) ListAddExtraAttribute(attribute string) error {

	if utils.IsAmongStr(attribute, "_id", "id", "identifier", "content", "created_at", "updated_at") {
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
func (it *DefaultCMSBlock) ListFilterAdd(Attribute string, Operator string, Value interface{}) error {
	collection, err := it.getListCollection()
	if err != nil {
		return err
	}

	collection.AddFilter(Attribute, Operator, Value.(string))
	return nil
}

// clears presets made by ListFilterAdd() and ListAddExtraAttribute() functions
func (it *DefaultCMSBlock) ListFilterReset() error {
	collection, err := it.getListCollection()
	if err != nil {
		return err
	}

	collection.ClearFilters()
	return nil
}

// select pagination
func (it *DefaultCMSBlock) ListLimit(offset int, limit int) error {
	collection, err := it.getListCollection()
	if err != nil {
		return err
	}

	return collection.SetLimit(offset, limit)
}
