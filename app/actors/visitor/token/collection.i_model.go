package token

import (
	"github.com/ottemo/commerce/app/models"
	"github.com/ottemo/commerce/app/models/visitor"
	"github.com/ottemo/commerce/db"
	"github.com/ottemo/commerce/env"
)

// GetModelName returns the Visitor Address model
func (it *DefaultVisitorCardCollection) GetModelName() string {
	return visitor.ConstModelNameVisitorCard
}

// GetImplementationName returns the Visitor Address implementation name
func (it *DefaultVisitorCardCollection) GetImplementationName() string {
	return "Default" + visitor.ConstModelNameVisitorCard
}

// New creates a new Visitor Address Collection
func (it *DefaultVisitorCardCollection) New() (models.InterfaceModel, error) {
	dbCollection, err := db.GetCollection(ConstCollectionNameVisitorToken)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	return &DefaultVisitorCardCollection{listCollection: dbCollection, listExtraAtributes: make([]string, 0)}, nil
}
