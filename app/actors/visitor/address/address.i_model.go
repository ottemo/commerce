package address

import (
	"github.com/ottemo/foundation/app/models"
	"github.com/ottemo/foundation/app/models/visitor"
)

// GetModelName returns the Visitor Address Model
func (it *DefaultVisitorAddress) GetModelName() string {
	return visitor.MODEL_NAME_VISITOR_ADDRESS
}

// GetImplementationName returns the Implementation name
func (it *DefaultVisitorAddress) GetImplementationName() string {
	return "Default" + visitor.MODEL_NAME_VISITOR_ADDRESS
}

// New creates a new Visitor Address interface
func (it *DefaultVisitorAddress) New() (models.I_Model, error) {
	return &DefaultVisitorAddress{}, nil
}
