package order

import (
	"github.com/ottemo/foundation/app/models"
	"github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/utils"
)

// enumerates items of Product model type
func (it *DefaultOrderItemCollection) List() ([]models.StructListItem, error) {
	result := make([]models.StructListItem, 0)

	dbRecords, err := it.listCollection.Load()
	if err != nil {
		return result, env.ErrorDispatch(err)
	}

	for _, dbRecordData := range dbRecords {

		orderItemModel := new(DefaultOrderItem)
		err = orderItemModel.FromHashMap(dbRecordData)
		if err != nil {
			return result, env.ErrorDispatch(err)
		}
		// retrieving minimal data needed for list
		resultItem := new(models.StructListItem)

		resultItem.ID = orderItemModel.GetID()
		resultItem.Name = orderItemModel.GetName()
		resultItem.Image = ""
		resultItem.Desc = utils.InterfaceToString(orderItemModel.Get("description"))

		// if extra attributes were required
		if len(it.listExtraAtributes) > 0 {
			resultItem.Extra = make(map[string]interface{})

			for _, attributeName := range it.listExtraAtributes {
				resultItem.Extra[attributeName] = orderItemModel.Get(attributeName)
			}
		}

		result = append(result, *resultItem)
	}

	return result, nil
}

// allows to obtain additional attributes from  List() function
func (it *DefaultOrderItemCollection) ListAddExtraAttribute(attribute string) error {

	orderItemModel := new(DefaultOrderItem)

	allowedAttributes := make([]string, 0)
	for _, attributeInfo := range orderItemModel.GetAttributesInfo() {
		allowedAttributes = append(allowedAttributes, attributeInfo.Attribute)
	}

	if utils.IsInArray(attribute, allowedAttributes) {
		if !utils.IsInListStr(attribute, it.listExtraAtributes) {
			it.listExtraAtributes = append(it.listExtraAtributes, attribute)
		} else {
			return env.ErrorNew("attribute already in list")
		}
	} else {
		return env.ErrorNew("not allowed attribute")
	}

	return nil
}

// adds selection filter to List() function
func (it *DefaultOrderItemCollection) ListFilterAdd(Attribute string, Operator string, Value interface{}) error {
	it.listCollection.AddFilter(Attribute, Operator, Value.(string))
	return nil
}

// clears presets made by ListFilterAdd() and ListAddExtraAttribute() functions
func (it *DefaultOrderItemCollection) ListFilterReset() error {
	it.listCollection.ClearFilters()
	return nil
}

// specifies selection paging
func (it *DefaultOrderItemCollection) ListLimit(offset int, limit int) error {
	return it.listCollection.SetLimit(offset, limit)
}
