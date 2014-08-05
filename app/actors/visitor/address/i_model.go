package address

import (
	"github.com/ottemo/foundation/app/models"
	"github.com/ottemo/foundation/app/models/visitor"
)

func (it *DefaultVisitorAddress) GetModelName() string {
	return visitor.VISITOR_ADDRESS_MODEL_NAME
}

func (it *DefaultVisitorAddress) GetImplementationName() string {
	return "DefaultVisitorAddress"
}

func (it *DefaultVisitorAddress) New() (models.I_Model, error) {
	return &DefaultVisitorAddress{}, nil
}
