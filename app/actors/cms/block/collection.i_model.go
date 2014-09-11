package block

import (
	"github.com/ottemo/foundation/app/models"
	"github.com/ottemo/foundation/app/models/cms"
	"github.com/ottemo/foundation/db"
)

// returns model name
func (it *DefaultCMSBlockCollection) GetModelName() string {
	return cms.MODEL_NAME_CMS_BLOCK_COLLECTION
}

// returns model implementation name
func (it *DefaultCMSBlockCollection) GetImplementationName() string {
	return "Default" + cms.MODEL_NAME_CMS_BLOCK_COLLECTION
}

// returns new instance of model implementation object
func (it *DefaultCMSBlockCollection) New() (models.I_Model, error) {
	dbCollection, err := db.GetCollection(CMS_BLOCK_COLLECTION_NAME)
	if err != nil {
		return nil, err
	}

	return &DefaultCMSBlockCollection{listCollection: dbCollection, listExtraAtributes: make([]string, 0)}, nil
}
