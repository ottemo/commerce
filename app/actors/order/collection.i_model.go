package order

import (
	"github.com/ottemo/foundation/db"

	"github.com/ottemo/foundation/app/models"
	"github.com/ottemo/foundation/app/models/order"
)

// returns model name
func (it *DefaultOrderCollection) GetModelName() string {
	return order.MODEL_NAME_ORDER_COLLECTION
}

// returns model implementation name
func (it *DefaultOrderCollection) GetImplementationName() string {
	return "Default" + order.MODEL_NAME_ORDER_COLLECTION
}

// returns new instance of model implementation object
func (it *DefaultOrderCollection) New() (models.I_Model, error) {
	dbCollection, err := db.GetCollection(COLLECTION_NAME_ORDER)
	if err != nil {
		return nil, err
	}

	return &DefaultOrderCollection{listCollection: dbCollection, listExtraAtributes: make([]string, 0)}, nil
}
