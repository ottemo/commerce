package category

import (
	"errors"

	"github.com/ottemo/foundation/app/models"
	"github.com/ottemo/foundation/app/models/category"
	"github.com/ottemo/foundation/utils"
)

// enumerates items of model type
func (it *DefaultCategoryCollection) List() ([]models.T_ListItem, error) {
	result := make([]models.T_ListItem, 0)

	// loading data from DB
	//---------------------
	dbItems, err := it.listCollection.Load()
	if err != nil {
		return result, err
	}

	// converting db record to T_ListItem
	//-----------------------------------
	for _, dbItemData := range dbItems {
		categoryModel, err := category.GetCategoryModel()
		if err != nil {
			return result, err
		}
		categoryModel.FromHashMap(dbItemData)

		// retrieving minimal data needed for list
		resultItem := new(models.T_ListItem)

		resultItem.Id = categoryModel.GetId()
		resultItem.Name = categoryModel.GetName()
		resultItem.Image = ""
		resultItem.Desc = ""

		// serving extra attributes
		//-------------------------
		if len(it.listExtraAtributes) > 0 {
			resultItem.Extra = make(map[string]interface{})

			for _, attributeName := range it.listExtraAtributes {
				resultItem.Extra[attributeName] = categoryModel.Get(attributeName)
			}
		}

		result = append(result, *resultItem)
	}

	return result, nil
}

// allows to obtain additional attributes from  List() function
func (it *DefaultCategoryCollection) ListAddExtraAttribute(attribute string) error {

	categoryModel, err := category.GetCategoryModel()
	if err != nil {
		return err
	}

	allowedAttributes := make([]string, 0)
	for _, attributeInfo := range categoryModel.GetAttributesInfo() {
		allowedAttributes = append(allowedAttributes, attributeInfo.Attribute)
	}
	allowedAttributes = append(allowedAttributes, "parent")

	if utils.IsInArray(attribute, allowedAttributes) {
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
func (it *DefaultCategoryCollection) ListFilterAdd(Attribute string, Operator string, Value interface{}) error {
	it.listCollection.AddFilter(Attribute, Operator, Value.(string))
	return nil
}

// clears presets made by ListFilterAdd() and ListAddExtraAttribute() functions
func (it *DefaultCategoryCollection) ListFilterReset() error {
	it.listCollection.ClearFilters()
	return nil
}

// specifies selection paging
func (it *DefaultCategoryCollection) ListLimit(offset int, limit int) error {
	return it.listCollection.SetLimit(offset, limit)
}
