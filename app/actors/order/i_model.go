package order

import (
	"github.com/ottemo/foundation/app/models"
	"github.com/ottemo/foundation/app/models/order"
)

// returns model name we have implementation for
func (it *DefaultOrder) GetModelName() string {
	return order.ORDER_MODEL_NAME
}

// returns name of current model implementation
func (it *DefaultOrder) GetImplementationName() string {
	return "DefaultOrder"
}

// makes new instance of model
func (it *DefaultOrder) New() (models.I_Model, error) {
	return &DefaultOrder { Items: make(map[int]order.I_OrderItem) }, nil
}
