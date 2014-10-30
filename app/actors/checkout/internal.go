package checkout

import (
	"github.com/ottemo/foundation/app/models/cart"
	"github.com/ottemo/foundation/api"
	"github.com/ottemo/foundation/app"
	"github.com/ottemo/foundation/app/models/checkout"
	"github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/utils"
	"github.com/ottemo/foundation/app/models/order"
)

// SetSession sets visitor for checkout
func (it *DefaultCheckout) SendOrderConfirmationMail() error {

	checkoutOrder := it.GetOrder()
	if checkoutOrder == nil {
		return env.ErrorNew("checkout order is not exists")
	}

	confirmationEmail := utils.InterfaceToString(env.ConfigGetValue(checkout.CONFIG_PATH_CONFIRMATION_EMAIL))
	if confirmationEmail != "" {
		email := utils.InterfaceToString(checkoutOrder.Get("customer_email"))
		if email == "" {
			return env.ErrorNew("customer email for order is not set")
		}

		confirmationEmail, err := utils.TextTemplate(confirmationEmail,
			map[string]interface{}{
				"Order":   checkoutOrder.ToHashMap(),
				"Visitor": it.GetVisitor().ToHashMap(),
			})
		if err != nil {
			return env.ErrorDispatch(err)
		}

		err = app.SendMail(email, "Order confirmation", confirmationEmail)
		if err != nil {
			return env.ErrorDispatch(err)
		}
	}

	return nil
}

func (it *DefaultCheckout) CheckoutSuccess(checkoutOrder order.I_Order, session api.I_Session) error {

	if checkoutOrder == nil || session == nil {
		return env.ErrorNew("Order or session is null")
	}

	currentCart := it.GetCart()

	err := checkoutOrder.Save()
	if err != nil {
		return err
	}

	// cleanup checkout information
	//-----------------------------
	currentCart.Deactivate()
	currentCart.Save()


	session.Set(cart.SESSION_KEY_CURRENT_CART, nil)
	session.Set(checkout.SESSION_KEY_CURRENT_CHECKOUT, nil)

	eventData := make(map[string]interface{})
	eventData["sessionId"] = session.GetId()
	env.Event("api.purchased", eventData)

	eventData = make(map[string]interface{})

	products := currentCart.GetItems()
	for i, _ := range products {
		eventData[products[i].GetProductId()] = products[i].GetQty()
	}
	env.Event("api.sales", eventData)

	return nil
}
