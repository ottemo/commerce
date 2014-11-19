package product

import (
	"github.com/ottemo/foundation/app/models"
	"github.com/ottemo/foundation/app/models/product"
	"github.com/ottemo/foundation/db"
	"github.com/ottemo/foundation/env"
)

// returns model name
func (it *DefaultProductCollection) GetModelName() string {
	return product.ConstModelNameProductCollection
}

// returns model implementation name
func (it *DefaultProductCollection) GetImplementationName() string {
	return "Default" + product.ConstModelNameProductCollection
}

// returns new instance of model implementation object
func (it *DefaultProductCollection) New() (models.InterfaceModel, error) {
	dbCollection, err := db.GetCollection(ConstCollectionNameProduct)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	return &DefaultProductCollection{listCollection: dbCollection, listExtraAtributes: make([]string, 0)}, nil
}
