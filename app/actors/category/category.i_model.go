package category

import (
	"github.com/ottemo/foundation/app/models"
	"github.com/ottemo/foundation/app/models/category"
)

// returns model name
func (it *DefaultCategory) GetModelName() string {
	return category.MODEL_NAME_CATEGORY
}

// returns model implementation name
func (it *DefaultCategory) GetImplementationName() string {
	return "Default" + category.MODEL_NAME_CATEGORY
}

// returns new instance of model implementation object
func (it *DefaultCategory) New() (models.I_Model, error) {
	return &DefaultCategory{ProductIds: make([]string, 0)}, nil
}
