package saleprice

// DefaultSalePriceCollection type implements:
// 	- InterfaceModel
//	- InterfaceCollection

import (
	"github.com/ottemo/commerce/db"
	"github.com/ottemo/commerce/env"
	"github.com/ottemo/commerce/utils"

	"github.com/ottemo/commerce/app/models"
	"github.com/ottemo/commerce/app/models/discount/saleprice"
)

// ---------------------------------------------------------------------------------------------------------------------
// InterfaceModel implementation (package "github.com/ottemo/commerce/app/models/interfaces")
// ---------------------------------------------------------------------------------------------------------------------

// GetModelName returns model name
func (it *DefaultSalePriceCollection) GetModelName() string {
	return saleprice.ConstSalePriceDbCollectionName
}

// GetImplementationName default model default implementation name
func (it *DefaultSalePriceCollection) GetImplementationName() string {
	return "Default" + saleprice.ConstSalePriceDbCollectionName
}

// New returns new instance of model implementation object
func (it *DefaultSalePriceCollection) New() (models.InterfaceModel, error) {
	dbCollection, err := db.GetCollection(saleprice.ConstSalePriceDbCollectionName)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	return &DefaultSalePriceCollection{listCollection: dbCollection, listExtraAtributes: make([]string, 0)}, nil
}

// ---------------------------------------------------------------------------------------------------------------------
// InterfaceCollection implementation (package "github.com/ottemo/commerce/app/models/interfaces")
// ---------------------------------------------------------------------------------------------------------------------

// GetDBCollection returns database collection
func (it *DefaultSalePriceCollection) GetDBCollection() db.InterfaceDBCollection {
	return it.listCollection
}

// List returns list of StructListItem items
func (it *DefaultSalePriceCollection) List() ([]models.StructListItem, error) {
	var result []models.StructListItem

	// loading data from DB
	//---------------------
	dbItems, err := it.listCollection.Load()
	if err != nil {
		return result, env.ErrorDispatch(err)
	}

	// converting db record to StructListItem
	//-----------------------------------
	for _, dbItemData := range dbItems {
		salePriceModel, err := saleprice.GetSalePriceModel()
		if err != nil {
			return result, env.ErrorDispatch(err)
		}
		if err := salePriceModel.FromHashMap(dbItemData); err != nil {
			_ = env.ErrorNew(ConstErrorModule, ConstErrorLevel, "b6b2d648-28ab-4895-bebb-49310eb74435", err.Error())
		}

		// retrieving minimal data needed for list
		resultItem := new(models.StructListItem)

		resultItem.ID = salePriceModel.GetID()

		// serving extra attributes
		//-------------------------
		if len(it.listExtraAtributes) > 0 {
			resultItem.Extra = make(map[string]interface{})

			for _, attributeName := range it.listExtraAtributes {
				resultItem.Extra[attributeName] = salePriceModel.Get(attributeName)
			}
		}

		result = append(result, *resultItem)
	}

	return result, nil
}

// ListAddExtraAttribute adds attribute to sale price collection
func (it *DefaultSalePriceCollection) ListAddExtraAttribute(attribute string) error {

	salePriceModel, err := saleprice.GetSalePriceModel()
	if err != nil {
		return env.ErrorDispatch(err)
	}

	var allowedAttributes []string
	for _, attributeInfo := range salePriceModel.GetAttributesInfo() {
		allowedAttributes = append(allowedAttributes, attributeInfo.Attribute)
	}

	if utils.IsInArray(attribute, allowedAttributes) {
		if !utils.IsInListStr(attribute, it.listExtraAtributes) {
			it.listExtraAtributes = append(it.listExtraAtributes, attribute)
		} else {
			return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "7df7cc4f-eace-4fb8-865a-3146ec310383", "attribute already in list")
		}
	} else {
		return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "0fc8658d-8755-4608-9b26-7ab8910f5b01", "not allowed attribute")
	}

	return nil
}

// ListFilterAdd adds filter to sale price collection
func (it *DefaultSalePriceCollection) ListFilterAdd(attribute string, operator string, value interface{}) error {
	if err := it.listCollection.AddFilter(attribute, operator, value.(string)); err != nil {
		_ = env.ErrorNew(ConstErrorModule, ConstErrorLevel, "60b83d74-7c12-4302-ae75-4f6877524c88", err.Error())
	}
	return nil
}

// ListFilterReset resets sale price collection filters
func (it *DefaultSalePriceCollection) ListFilterReset() error {
	if err := it.listCollection.ClearFilters(); err != nil {
		_ = env.ErrorNew(ConstErrorModule, ConstErrorLevel, "9695f278-5d3c-4d5d-b657-988aa745a10b", err.Error())
	}
	return nil
}

// ListLimit limits sale price collection selected records
func (it *DefaultSalePriceCollection) ListLimit(offset int, limit int) error {
	return it.listCollection.SetLimit(offset, limit)
}

// ---------------------------------------------------------------------------------------------------------------------
//  implementation (package "github.com/ottemo/commerce/app/models/interfaces")
// ---------------------------------------------------------------------------------------------------------------------

// ListSalePrices returns list of sale price model items
func (it *DefaultSalePriceCollection) ListSalePrices() []saleprice.InterfaceSalePrice {
	var result []saleprice.InterfaceSalePrice

	dbRecords, err := it.listCollection.Load()
	if err != nil {
		return result
	}

	for _, recordData := range dbRecords {
		salePriceModel, err := saleprice.GetSalePriceModel()
		if err != nil {
			return result
		}
		if err := salePriceModel.FromHashMap(recordData); err != nil {
			_ = env.ErrorNew(ConstErrorModule, ConstErrorLevel, "c300b838-0439-4076-b62a-4c5003a6bdc2", err.Error())
		}

		result = append(result, salePriceModel)
	}

	return result
}
