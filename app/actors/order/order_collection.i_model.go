package order

import (
	"github.com/ottemo/commerce/app/models"
	"github.com/ottemo/commerce/app/models/order"
	"github.com/ottemo/commerce/db"
	"github.com/ottemo/commerce/env"
)

// GetModelName returns model name
func (it *DefaultOrderCollection) GetModelName() string {
	return order.ConstModelNameOrderCollection
}

// GetImplementationName returns model implementation name
func (it *DefaultOrderCollection) GetImplementationName() string {
	return "Default" + order.ConstModelNameOrderCollection
}

// New returns new instance of model implementation object
func (it *DefaultOrderCollection) New() (models.InterfaceModel, error) {
	dbCollection, err := db.GetCollection(ConstCollectionNameOrder)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	return &DefaultOrderCollection{listCollection: dbCollection, listExtraAtributes: make([]string, 0)}, nil
}
