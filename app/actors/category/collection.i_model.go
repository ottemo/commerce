package category

import (
	"github.com/ottemo/commerce/app/models"
	"github.com/ottemo/commerce/app/models/category"
	"github.com/ottemo/commerce/db"
	"github.com/ottemo/commerce/env"
)

// GetModelName returns model name
func (it *DefaultCategoryCollection) GetModelName() string {
	return category.ConstModelNameCategory
}

// GetImplementationName returns model implementation name
func (it *DefaultCategoryCollection) GetImplementationName() string {
	return "Default" + category.ConstModelNameCategory
}

// New returns new instance of model implementation object
func (it *DefaultCategoryCollection) New() (models.InterfaceModel, error) {
	dbCollection, err := db.GetCollection(ConstCollectionNameCategory)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	return &DefaultCategoryCollection{listCollection: dbCollection, listExtraAtributes: make([]string, 0)}, nil
}
