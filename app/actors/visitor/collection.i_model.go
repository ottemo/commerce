package visitor

import (
	"github.com/ottemo/commerce/app/models"
	"github.com/ottemo/commerce/app/models/visitor"
	"github.com/ottemo/commerce/db"
	"github.com/ottemo/commerce/env"
)

// GetModelName returns model name for the Visitor Collection
func (it *DefaultVisitorCollection) GetModelName() string {
	return visitor.ConstModelNameVisitorCollection
}

// GetImplementationName returns model implementation name for the Visitor Collection
func (it *DefaultVisitorCollection) GetImplementationName() string {
	return "Default" + visitor.ConstModelNameVisitorCollection
}

// New returns new instance of model implementation object for the Visitor Collection
func (it *DefaultVisitorCollection) New() (models.InterfaceModel, error) {
	dbCollection, err := db.GetCollection(ConstCollectionNameVisitor)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	return &DefaultVisitorCollection{listCollection: dbCollection, listExtraAtributes: make([]string, 0)}, nil
}
