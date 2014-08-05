package visitor

import (
	"github.com/ottemo/foundation/app/models"
	"github.com/ottemo/foundation/app/models/visitor"
)

func (it *DefaultVisitor) GetModelName() string {
	return visitor.VISITOR_MODEL_NAME
}

func (it *DefaultVisitor) GetImplementationName() string {
	return "DefaultVisitor"
}

func (it *DefaultVisitor) New() (models.I_Model, error) {
	return &DefaultVisitor{}, nil
}
