package block

import (
	"github.com/ottemo/foundation/app/models"
	"github.com/ottemo/foundation/app/models/cms"
	"github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/utils"
)

// enumerates items of CMS block model
func (it *DefaultCMSBlockCollection) List() ([]models.T_ListItem, error) {
	result := make([]models.T_ListItem, 0)

	dbRecords, err := it.listCollection.Load()
	if err != nil {
		return result, env.ErrorDispatch(err)
	}

	for _, dbRecordData := range dbRecords {
		cmsBlockModel, err := cms.GetCMSBlockModel()
		if err != nil {
			return result, env.ErrorDispatch(err)
		}
		cmsBlockModel.FromHashMap(dbRecordData)

		// retrieving minimal data needed for list
		resultItem := new(models.T_ListItem)

		resultItem.Id = cmsBlockModel.GetId()
		resultItem.Name = cmsBlockModel.GetIdentifier()
		resultItem.Image = ""
		resultItem.Desc = ""

		// if extra attributes were required
		if len(it.listExtraAtributes) > 0 {
			resultItem.Extra = make(map[string]interface{})

			for _, attributeName := range it.listExtraAtributes {
				resultItem.Extra[attributeName] = cmsBlockModel.Get(attributeName)
			}
		}

		result = append(result, *resultItem)
	}

	return result, nil
}

// allows to obtain additional attributes from  List() function
func (it *DefaultCMSBlockCollection) ListAddExtraAttribute(attribute string) error {

	if utils.IsAmongStr(attribute, "_id", "id", "identifier", "content", "created_at", "updated_at") {
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
func (it *DefaultCMSBlockCollection) ListFilterAdd(Attribute string, Operator string, Value interface{}) error {
	it.listCollection.AddFilter(Attribute, Operator, Value.(string))
	return nil
}

// clears presets made by ListFilterAdd() and ListAddExtraAttribute() functions
func (it *DefaultCMSBlockCollection) ListFilterReset() error {
	it.listCollection.ClearFilters()
	return nil
}

// sets select pagination
func (it *DefaultCMSBlockCollection) ListLimit(offset int, limit int) error {
	return it.listCollection.SetLimit(offset, limit)
}
