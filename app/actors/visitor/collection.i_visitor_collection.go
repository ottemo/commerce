package visitor

import (
	"github.com/ottemo/foundation/app/models/visitor"
	"github.com/ottemo/foundation/db"
)

// GetDBCollection returns database collection for the Visitor
func (it *DefaultVisitorCollection) GetDBCollection() db.I_DBCollection {
	return it.listCollection
}

// ListVisitors returns list of visitor model items in the Visitor Collection
func (it *DefaultVisitorCollection) ListVisitors() []visitor.I_Visitor {
	var result []visitor.I_Visitor

	dbRecords, err := it.listCollection.Load()
	if err != nil {
		return result
	}

	for _, recordData := range dbRecords {
		visitorModel, err := visitor.GetVisitorModel()
		if err != nil {
			return result
		}
		visitorModel.FromHashMap(recordData)

		result = append(result, visitorModel)
	}

	return result
}
