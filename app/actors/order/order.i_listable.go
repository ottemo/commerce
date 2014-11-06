package order

import (
	"github.com/ottemo/foundation/app/models"
	"github.com/ottemo/foundation/app/models/order"
)

// returns collection of current instance type
func (it *DefaultOrder) GetCollection() models.I_Collection {
	model, _ := models.GetModel(order.MODEL_NAME_ORDER_COLLECTION)
	if result, ok := model.(order.I_OrderCollection); ok {
		return result
	}

	return nil
}
