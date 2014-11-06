package visitor

import (
	"github.com/ottemo/foundation/app/models"
	"github.com/ottemo/foundation/app/models/visitor"
)

// returns collection of current instance type
func (it *DefaultVisitor) GetCollection() models.I_Collection {
	model, _ := models.GetModel(visitor.MODEL_NAME_VISITOR_COLLECTION)
	if result, ok := model.(visitor.I_VisitorCollection); ok {
		return result
	}

	return nil
}
