package address

import (
	"errors"
	"github.com/ottemo/foundation/app/models"
	"github.com/ottemo/foundation/db"

	"github.com/ottemo/foundation/app/utils"
)

//---------------------------------
// IMPLEMENTATION SPECIFIC METHODS
//---------------------------------

// initializes and returns shared among couple functions collection
func (it *DefaultVisitorAddress) getListCollection() (db.I_DBCollection, error) {
	if it.listCollection != nil {
		return it.listCollection, nil
	} else {
		var err error = nil

		dbEngine := db.GetDBEngine()
		if dbEngine == nil {
			return nil, errors.New("Can't obtain DBEngine")
		}

		it.listCollection, err = dbEngine.GetCollection(VISITOR_ADDRESS_COLLECTION_NAME)

		return it.listCollection, err
	}
}


//--------------------------
// INTERFACE IMPLEMENTATION
//--------------------------

// enumerates items of VisitorAddress model type
func (it *DefaultVisitorAddress) List() ([]models.T_ListItem, error) {
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

		address := model.(*DefaultVisitorAddress)
		address.FromHashMap(dbItemData)

		// retrieving minimal data needed for list
		resultItem := new(models.T_ListItem)

		resultItem.Id    = address.GetId()
		resultItem.Name  = address.GetZipCode() + " " + address.GetState() + ", " + address.GetCity() + ", " + address.GetStreet()
		resultItem.Image = ""
		resultItem.Desc  = "Zip: " + address.GetZipCode() + ", State: " + address.GetState() + ", City: " + address.GetCity() + ", Street: " + address.GetStreet() + ", Phone: " + address.GetPhone()

		// if extra attributes were required
		if len(it.listExtraAtributes) > 0 {
			resultItem.Extra = make(map[string]interface{})

			for _, attributeName := range it.listExtraAtributes {
				resultItem.Extra[attributeName] = address.Get(attributeName)
			}
		}

		result = append(result, *resultItem)
	}

	return result, nil
}



// allows to obtain additional attributes from  List() function
func (it *DefaultVisitorAddress) ListAddExtraAttribute(attribute string) error {

	if utils.IsAmongStr(attribute, "billing_address", "shipping_address") {
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
func (it *DefaultVisitorAddress) ListFilterAdd(Attribute string, Operator string, Value interface{}) error {
	collection, err := it.getListCollection()
	if err != nil {
		return err
	}

	collection.AddFilter(Attribute, Operator, Value.(string))
	return nil
}



// clears presets made by ListFilterAdd() and ListAddExtraAttribute() functions
func (it *DefaultVisitorAddress) ListFilterReset() error {
	collection, err := it.getListCollection()
	if err != nil {
		return err
	}

	collection.ClearFilters()
	return nil
}



// specifies selection paging
func (it *DefaultVisitorAddress) ListLimit(offset int, limit int) error {
	collection, err := it.getListCollection()
	if err != nil {
		return err
	}

	return collection.SetLimit(offset, limit)
}
