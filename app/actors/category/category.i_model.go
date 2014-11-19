package category

import (
	"github.com/ottemo/foundation/app/models"
	"github.com/ottemo/foundation/app/models/category"
)

// returns model name
func (it *DefaultCategory) GetModelName() string {
	return category.ConstModelNameCategory
}

// returns model implementation name
func (it *DefaultCategory) GetImplementationName() string {
	return "Default" + category.ConstModelNameCategory
}

// returns new instance of model implementation object
func (it *DefaultCategory) New() (models.InterfaceModel, error) {
	return &DefaultCategory{ProductIds: make([]string, 0)}, nil
}
