package page

import (
	"github.com/ottemo/foundation/app/models"
	"github.com/ottemo/foundation/app/models/cms"
	"github.com/ottemo/foundation/db"
)

// returns model name
func (it *DefaultCMSPageCollection) GetModelName() string {
	return cms.MODEL_NAME_CMS_PAGE_COLLECTION
}

// returns model implementation name
func (it *DefaultCMSPageCollection) GetImplementationName() string {
	return "Default" + cms.MODEL_NAME_CMS_PAGE_COLLECTION
}

// returns new instance of model implementation object
func (it *DefaultCMSPageCollection) New() (models.I_Model, error) {
	dbCollection, err := db.GetCollection(CMS_PAGE_COLLECTION_NAME)
	if err != nil {
		return nil, err
	}

	return &DefaultCMSPageCollection{listCollection: dbCollection, listExtraAtributes: make([]string, 0)}, nil
}
