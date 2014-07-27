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



// sets new qty for particular item in cart
func (it *DefaultCart) SetQty(itemId string, qty int) error {
	return nil
}



// enumerates current cart items
func (it *DefaultCart) ListItems() []cart.I_CartItem {
	return it.Items
}



// returns visitor id this cart belongs to
func (it *DefaultCart) GetVisitorId() string {
	return it.VisitorId
}



// sets new owner of cart
func (it *DefaultCart) SetVisitorId(visitorId string) error {
	it.VisitorId = visitorId
	return nil
}



// returns visitor model represents owner or current cart or nil if visitor was not set to cart
func (it *DefaultCart) GetVisitor() visitor.I_Visitor {
	visitor, _ := visitor.LoadVisitorById(it.VisitorId)
	return visitor
}



// assigns some information to current cart
func (it *DefaultCart) SetCartInfo(infoAttribute string, infoValue interface{}) error {
	if it.Info == nil {
		it.Info = make(map[string]interface{})
	}

	it.Info[infoAttribute] = infoValue

	return nil
}


// returns current cart info assigned
func (it *DefaultCart) GetCartInfo(infoAttribute string) interface{} {
	return it.Info
}
