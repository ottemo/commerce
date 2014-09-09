package listable

import (
	"github.com/ottemo/foundation/db"
	"github.com/ottemo/foundation/app/models"
)

type ListableHelperDelegates struct {
	CollectionName string
	GetCollection  func() db.I_DBCollection

	ValidateExtraAttributeFunc func(attribute string) bool

	RecordToObjectFunc   func(recordData map[string]interface{}, extraAttributes []string) interface{}
	RecordToListItemFunc func(recordData map[string]interface{}, extraAttributes []string) (models.T_ListItem, bool)
}

type ListableHelper struct {
	delegate ListableHelperDelegates

	listCollection     db.I_DBCollection
	listExtraAtributes []string
}

// use this function to obtain ListableHelper struct for your object
func NewListableHelper(delegates ListableHelperDelegates) *ListableHelper {
	return &ListableHelper{delegate: delegates}
}
