package cart

import (
	"errors"
	"strconv"
	"sort"

	"github.com/ottemo/foundation/db"
	"github.com/ottemo/foundation/app/models/cart"
	"github.com/ottemo/foundation/app/models/visitor"
	"github.com/ottemo/foundation/app/models/product"
)



// adds item to the current cart
//   - returns added item or nil if error happened
func (it *DefaultCart) AddItem(productId string, qty int, options map[string]interface{}) (cart.I_CartItem, error) {
	if qty <= 0 {
		return nil, errors.New("qty can't be zero or less")
	}

	reqProduct, err := product.LoadProductById(productId)
	if err != nil {
		return nil, err
	}

	if options == nil {
		options = make(map[string]interface{})
	}

	it.maxIdx += 1

	cartItem := &DefaultCartItem {
				idx: it.maxIdx,
				ProductId: reqProduct.GetId(),
				Qty: qty,
				Options: options,
				Cart: it }

	it.Items[it.maxIdx] = cartItem

	return cartItem, nil
}



// removes item from cart
//   - you need to know index you can get from ListItems()
func (it *DefaultCart) RemoveItem(itemIdx int) error {
	if _, present := it.Items[itemIdx]; present {
		delete(it.Items, itemIdx)
		return nil
	} else {
		return errors.New("can't find index " + strconv.Itoa(itemIdx))
	}
}



// sets new qty for particular item in cart
//   - you need to it's index, use ListItems() for that
func (it *DefaultCart) SetQty(itemIdx int, qty int) error {
	cartItem, present := it.Items[itemIdx]
	if present {
		return cartItem.SetQty(qty)
	} else {
		return errors.New("there is no item with idx=" + strconv.Itoa(itemIdx) )
	}
}



// enumerates current cart items sorted by item idx
func (it *DefaultCart) ListItems() []cart.I_CartItem {

	result := make([]cart.I_CartItem, 0)

	keys := make([]int, 0)
	for key, _ := range it.Items {
		keys = append(keys, key)
	}

	sort.Ints(keys)

	for _, key := range keys {
		result = append(result, it.Items[key])
	}

	return result
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
func (it *DefaultCart) GetCartInfo() map[string]interface{} {
	return it.Info
}



// loads cart information from DB for visitor
func (it *DefaultCart) MakeCartForVisitor(visitorId string) error {
	dbEngine := db.GetDBEngine()
	if dbEngine != nil {
		return errors.New("can't get DB Engine")
	}

	cartCollection, err := dbEngine.GetCollection( CART_COLLECTION_NAME )
	if err != nil {
		return err
	}

	cartCollection.AddFilter("visitor_id", "=", visitorId)
	cartCollection.AddFilter("active", "=", "true")
	rowsData, err := cartCollection.Load()
	if err != nil {
		return err
	}

	if len(rowsData) < 1 {
		newModel, err := it.New()
		if err != nil {
			return err
		}
		newCart := newModel.(*DefaultCart)
		newCart.SetVisitorId(visitorId)
		newCart.Save()
	} else {
		err := it.Load( rowsData[0]["cart_id"].(string) )
		if err != nil {
			return err
		}
	}

	return nil
}



func (it *DefaultCart) Activate() error {
	it.Active = true
	return nil
}

func (it *DefaultCart) Deactivate() error {
	it.Active = false
	return nil
}

func (it *DefaultCart) IsActive() bool {
	return it.Active
}
