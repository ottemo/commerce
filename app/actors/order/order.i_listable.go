package order

import (
	"github.com/ottemo/commerce/app/models"
	"github.com/ottemo/commerce/app/models/order"
)

// GetCollection returns collection of current instance type
func (it *DefaultOrder) GetCollection() models.InterfaceCollection {
	model, _ := models.GetModel(order.ConstModelNameOrderCollection)
	if result, ok := model.(order.InterfaceOrderCollection); ok {
		return result
	}

	return nil
}
