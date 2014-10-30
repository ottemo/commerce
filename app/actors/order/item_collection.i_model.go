package order

import (
	"github.com/ottemo/foundation/app/models"
	"github.com/ottemo/foundation/app/models/order"
	"github.com/ottemo/foundation/db"
	"github.com/ottemo/foundation/env"
)

// returns model name
func (it *DefaultOrderItemCollection) GetModelName() string {
	return order.MODEL_NAME_ORDER_ITEM_COLLECTION
}

// returns model implementation name
func (it *DefaultOrderItemCollection) GetImplementationName() string {
	return "Default" + order.MODEL_NAME_ORDER_ITEM_COLLECTION
}

// returns new instance of model implementation object
func (it *DefaultOrderItemCollection) New() (models.I_Model, error) {
	dbCollection, err := db.GetCollection(COLLECTION_NAME_ORDER_ITEMS)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	return &DefaultOrderItemCollection{listCollection: dbCollection, listExtraAtributes: make([]string, 0)}, nil
}
