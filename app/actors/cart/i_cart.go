package cart

import (
	"github.com/ottemo/foundation/app/models/cart"
	"github.com/ottemo/foundation/app/models/visitor"
)


func (it *DefaultCart) AddItem(qty int) error {
	return nil
}

func (it *DefaultCart) AddItemEx(qty int, options map[string]interface{}) error {
	return nil
}

func (it *DefaultCart) RemoveItem(itemId string) error {
	return nil
}

func (it *DefaultCart) SetQty(itemId string, qty int) error {
	return nil
}

func (it *DefaultCart) ListItems() []cart.I_CartItem {
	return nil
}

func (it *DefaultCart) GetVisitorId() string {
	return it.VisitorId
}

func (it *DefaultCart) SetVisitorId(visitorId string) error {
	it.VisitorId = visitorId
	return nil
}

func (it *DefaultCart) GetVisitor() visitor.I_Visitor {
	return nil
}

func (it *DefaultCart) SetCartInfo(infoAttribute string, infoValue interface{}) error {
	return nil
}

func (it *DefaultCart) GetCartInfo(infoAttribute string) interface{} {
	return nil
}
