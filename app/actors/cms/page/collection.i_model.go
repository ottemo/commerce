package page

import (
	"github.com/ottemo/commerce/app/models"
	"github.com/ottemo/commerce/app/models/cms"
	"github.com/ottemo/commerce/db"
	"github.com/ottemo/commerce/env"
)

// GetModelName returns model name
func (it *DefaultCMSPageCollection) GetModelName() string {
	return cms.ConstModelNameCMSPageCollection
}

// GetImplementationName returns model implementation name
func (it *DefaultCMSPageCollection) GetImplementationName() string {
	return "Default" + cms.ConstModelNameCMSPageCollection
}

// New returns new instance of model implementation object
func (it *DefaultCMSPageCollection) New() (models.InterfaceModel, error) {
	dbCollection, err := db.GetCollection(ConstCmsPageCollectionName)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	return &DefaultCMSPageCollection{listCollection: dbCollection, listExtraAtributes: make([]string, 0)}, nil
}
