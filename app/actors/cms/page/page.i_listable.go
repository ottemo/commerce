package page

import (
	"github.com/ottemo/foundation/app/models"
	"github.com/ottemo/foundation/app/models/cms"
)

// returns collection of current instance type
func (it *DefaultCMSPage) GetCollection() models.I_Collection {
	model, _ := models.GetModel(cms.MODEL_NAME_CMS_PAGE_COLLECTION)
	if result, ok := model.(cms.I_CMSPageCollection); ok {
		return result
	}

	return nil
}
