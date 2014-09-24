package page

import (
	"github.com/ottemo/foundation/app/models"
	"github.com/ottemo/foundation/app/models/cms"
)

// returns model name
func (it *DefaultCMSPage) GetModelName() string {
	return cms.MODEL_NAME_CMS_PAGE
}

// returns model implementation name
func (it *DefaultCMSPage) GetImplementationName() string {
	return "DefaultCMSPage"
}

// returns new instance of model implementation object
func (it *DefaultCMSPage) New() (models.I_Model, error) {
	return &DefaultCMSPage{}, nil
}
