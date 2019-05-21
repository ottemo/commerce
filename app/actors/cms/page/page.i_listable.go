package page

import (
	"github.com/ottemo/commerce/app/models"
	"github.com/ottemo/commerce/app/models/cms"
)

// GetCollection returns collection of current instance type
func (it *DefaultCMSPage) GetCollection() models.InterfaceCollection {
	model, _ := models.GetModel(cms.ConstModelNameCMSPageCollection)
	if result, ok := model.(cms.InterfaceCMSPageCollection); ok {
		return result
	}

	return nil
}
