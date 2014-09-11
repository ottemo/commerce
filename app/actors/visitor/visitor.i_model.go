package visitor

import (
	"github.com/ottemo/foundation/app/models"
	"github.com/ottemo/foundation/app/models/visitor"
)

// returns model name
func (it *DefaultVisitor) GetModelName() string {
	return visitor.MODEL_NAME_VISITOR
}

// returns model implementation name
func (it *DefaultVisitor) GetImplementationName() string {
	return "Default" + visitor.MODEL_NAME_VISITOR
}

// returns new instance of model implementation object
func (it *DefaultVisitor) New() (models.I_Model, error) {
	return &DefaultVisitor{}, nil
}
