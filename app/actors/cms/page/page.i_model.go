package page

import (
	"github.com/ottemo/foundation/app/models"
	"github.com/ottemo/foundation/app/models/cms"
)

// returns model name
func (it *DefaultCMSPage) GetModelName() string {
	return cms.ConstModelNameCMSPage
}

// returns model implementation name
func (it *DefaultCMSPage) GetImplementationName() string {
	return "DefaultCMSPage"
}

// returns new instance of model implementation object
func (it *DefaultCMSPage) New() (models.InterfaceModel, error) {
	return &DefaultCMSPage{}, nil
}
