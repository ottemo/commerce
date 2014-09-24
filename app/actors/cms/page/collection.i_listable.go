package page

import (
	"errors"

	"github.com/ottemo/foundation/app/models"
	"github.com/ottemo/foundation/app/models/cms"
	"github.com/ottemo/foundation/utils"
)

// enumerates items of CMS page model
func (it *DefaultCMSPageCollection) List() ([]models.T_ListItem, error) {
	result := make([]models.T_ListItem, 0)

	dbRecords, err := it.listCollection.Load()
	if err != nil {
		return result, err
	}

	for _, dbRecordData := range dbRecords {
		cmsPageModel, err := cms.GetCMSPageModel()
		if err != nil {
			return result, err
		}
		cmsPageModel.FromHashMap(dbRecordData)

		// retrieving minimal data needed for list
		resultItem := new(models.T_ListItem)

		resultItem.Id = cmsPageModel.GetId()
		resultItem.Name = cmsPageModel.GetIdentifier()
		resultItem.Image = ""
		resultItem.Desc = cmsPageModel.GetTitle()

		// if extra attributes were required
		if len(it.listExtraAtributes) > 0 {
			resultItem.Extra = make(map[string]interface{})

			for _, attributeName := range it.listExtraAtributes {
				resultItem.Extra[attributeName] = cmsPageModel.Get(attributeName)
			}
		}

		result = append(result, *resultItem)
	}

	return result, nil
}

// allows to obtain additional attributes from  List() function
func (it *DefaultCMSPageCollection) ListAddExtraAttribute(attribute string) error {

	if utils.IsAmongStr(attribute, "_id", "id", "url", "identifier", "title", "content", "meta_title", "meta_description", "created_at", "updated_at") {
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
func (it *DefaultCMSPageCollection) ListFilterAdd(Attribute string, Operator string, Value interface{}) error {
	it.listCollection.AddFilter(Attribute, Operator, Value.(string))
	return nil
}

// clears presets made by ListFilterAdd() and ListAddExtraAttribute() functions
func (it *DefaultCMSPageCollection) ListFilterReset() error {
	it.listCollection.ClearFilters()
	return nil
}

// sets select pagination
func (it *DefaultCMSPageCollection) ListLimit(offset int, limit int) error {
	return it.listCollection.SetLimit(offset, limit)
}
