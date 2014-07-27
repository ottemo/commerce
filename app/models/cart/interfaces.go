package cart

import (
	"github.com/ottemo/foundation/app/models"

	"github.com/ottemo/foundation/app/models/product"
	"github.com/ottemo/foundation/app/models/visitor"
)

type I_CartItem interface {
	GetId() string
	SetId(newId string) error

	GetProductId() string
	GetProduct() product.I_Product

	GetQty() int
	SetQty(qty int) error

	GetOptions() map[string]interface{}
	SetOption(optionName string, optionValue interface{}) error

	GetCart() I_Cart
}

type I_Cart interface {
	AddItem(qty int) error
	AddItemEx(qty int, options map[string]interface{}) error

	RemoveItem(itemId string) error

	SetQty(itemId string, qty int) error

	ListItems() []I_CartItem

	GetVisitorId() string
	SetVisitorId(string) error

	GetVisitor() visitor.I_Visitor

	SetCartInfo(infoAttribute string, infoValue interface{}) error
	GetCartInfo(infoAttribute string) interface{}

	models.I_Model
	models.I_Storable
}
