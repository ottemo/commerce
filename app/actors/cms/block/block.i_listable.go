package block

import (
	"github.com/ottemo/foundation/app/models"
	"github.com/ottemo/foundation/app/models/cms"
)

// returns collection of current instance type
func (it *DefaultCMSBlock) GetCollection() models.I_Collection {
	model, _ := models.GetModel(cms.MODEL_NAME_CMS_BLOCK_COLLECTION)
	if result, ok := model.(cms.I_CMSBlockCollection); ok {
		return result
	}

	return nil
}
