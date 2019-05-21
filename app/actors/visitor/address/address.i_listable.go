package address

import (
	"github.com/ottemo/commerce/app/models"
	"github.com/ottemo/commerce/app/models/visitor"
)

// GetCollection returns collection of current instance type
func (it *DefaultVisitorAddress) GetCollection() models.InterfaceCollection {
	model, _ := models.GetModel(visitor.ConstModelNameVisitorAddressCollection)
	if result, ok := model.(visitor.InterfaceVisitorAddressCollection); ok {
		return result
	}

	return nil
}
