package order

import (
	"github.com/ottemo/foundation/app/models/order"
	"github.com/ottemo/foundation/db"
)

// returns database collection
func (it *DefaultOrderCollection) GetDBCollection() db.I_DBCollection {
	return it.listCollection
}

// returns array of products in model instance form
func (it *DefaultOrderCollection) ListOrders() []order.I_Order {
	result := make([]order.I_Order, 0)

	dbRecords, err := it.listCollection.Load()
	if err != nil {
		return result
	}

	for _, dbRecordData := range dbRecords {
		orderModel, err := order.GetOrderModel()
		if err != nil {
			return result
		}
		orderModel.FromHashMap(dbRecordData)

		result = append(result, orderModel)
	}

	return result
}
