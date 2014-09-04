package page

import (
	"github.com/ottemo/foundation/app/models"
	"github.com/ottemo/foundation/app/models/cms"
)

func (it *DefaultCMSPage) GetModelName() string {
	return cms.CMS_PAGE_MODEL_NAME
}

func (it *DefaultCMSPage) GetImplementationName() string {
	return "DefaultCMSPage"
}

func (it *DefaultCMSPage) New() (models.I_Model, error) {
	return &DefaultCMSPage{}, nil
}
