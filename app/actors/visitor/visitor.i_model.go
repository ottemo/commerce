package visitor

import (
	"github.com/ottemo/foundation/app/helpers/attributes"
	"github.com/ottemo/foundation/app/models"
	"github.com/ottemo/foundation/app/models/visitor"
	"github.com/ottemo/foundation/env"
)

// GetModelName returns model name for the Visitor
func (it *DefaultVisitor) GetModelName() string {
	return visitor.MODEL_NAME_VISITOR
}

// GetImplementationName returns model implementation name for the Visitor
func (it *DefaultVisitor) GetImplementationName() string {
	return "Default" + visitor.MODEL_NAME_VISITOR
}

// New returns new instance of model implementation object for the Visitor
func (it *DefaultVisitor) New() (models.I_Model, error) {

	customAttributes, err := new(attributes.CustomAttributes).Init(visitor.MODEL_NAME_VISITOR, COLLECTION_NAME_VISITOR)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	return &DefaultVisitor{CustomAttributes: customAttributes}, nil
}
