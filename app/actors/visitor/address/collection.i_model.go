package address

import (
	"github.com/ottemo/foundation/app/models"
	"github.com/ottemo/foundation/app/models/visitor"
	"github.com/ottemo/foundation/db"
	"github.com/ottemo/foundation/env"
)

// GetModelName returns the Visitor Address model
func (it *DefaultVisitorAddressCollection) GetModelName() string {
	return visitor.MODEL_NAME_VISITOR_ADDRESS
}

// GetImplementationName returns the Visitor Address implementation name
func (it *DefaultVisitorAddressCollection) GetImplementationName() string {
	return "Default" + visitor.MODEL_NAME_VISITOR_ADDRESS
}

// New creates a new Visitor Address Collection
func (it *DefaultVisitorAddressCollection) New() (models.I_Model, error) {
	dbCollection, err := db.GetCollection(COLLECTION_NAME_VISITOR_ADDRESS)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	return &DefaultVisitorAddressCollection{listCollection: dbCollection, listExtraAtributes: make([]string, 0)}, nil
}
