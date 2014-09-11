package category

import (
	"github.com/ottemo/foundation/app/models"
	"github.com/ottemo/foundation/app/models/category"
	"github.com/ottemo/foundation/db"
)

// returns model name
func (it *DefaultCategoryCollection) GetModelName() string {
	return category.MODEL_NAME_CATEGORY
}

// returns model implementation name
func (it *DefaultCategoryCollection) GetImplementationName() string {
	return "Default" + category.MODEL_NAME_CATEGORY
}

// returns new instance of model implementation object
func (it *DefaultCategoryCollection) New() (models.I_Model, error) {
	dbCollection, err := db.GetCollection(COLLECTION_NAME_CATEGORY)
	if err != nil {
		return nil, err
	}

	return &DefaultCategoryCollection{listCollection: dbCollection, listExtraAtributes: make([]string, 0)}, nil
}
