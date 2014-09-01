package block

import (
	"github.com/ottemo/foundation/app/models"
	"github.com/ottemo/foundation/app/models/cms"
)

func (it *DefaultCMSBlock) GetModelName() string {
	return cms.CMS_BLOCK_MODEL_NAME
}

func (it *DefaultCMSBlock) GetImplementationName() string {
	return "DefaultCMSBlock"
}

func (it *DefaultCMSBlock) New() (models.I_Model, error) {
	return &DefaultCMSBlock{}, nil
}
