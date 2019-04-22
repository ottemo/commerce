package token

import (
	"github.com/ottemo/commerce/app/models/visitor"
	"github.com/ottemo/commerce/db"
	"github.com/ottemo/commerce/env"
)

// GetDBCollection returns the database collection of the Visitor Cards
func (it *DefaultVisitorCardCollection) GetDBCollection() db.InterfaceDBCollection {
	return it.listCollection
}

// ListVisitorsCards returns list of visitor model items for the Visitor Cards
func (it *DefaultVisitorCardCollection) ListVisitorsCards() []visitor.InterfaceVisitorCard {
	var result []visitor.InterfaceVisitorCard

	dbRecords, err := it.listCollection.Load()
	if err != nil {
		return result
	}

	for _, recordData := range dbRecords {
		visitorCardModel, err := visitor.GetVisitorCardModel()
		if err != nil {
			return result
		}
		if err := visitorCardModel.FromHashMap(recordData); err != nil {
			_ = env.ErrorDispatch(err)
		}

		result = append(result, visitorCardModel)
	}

	return result
}
