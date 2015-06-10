package checkout

import (
	"fmt"
	"time"

	"github.com/ottemo/foundation/api"
	"github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/utils"

	"github.com/ottemo/foundation/app/models/cart"
	"github.com/ottemo/foundation/app/models/checkout"
	"github.com/ottemo/foundation/app/models/order"
	"github.com/ottemo/foundation/app/models/visitor"
)

// SetShippingAddress sets shipping address for checkout
func (dc *DefaultCheckout) SetShippingAddress(address visitor.InterfaceVisitorAddress) error {
	if address == nil {
		dc.ShippingAddress = nil
		return nil
	}

	dc.ShippingAddress = address.ToHashMap()
	return nil
}

// GetShippingAddress returns checkout shipping address
func (dc *DefaultCheckout) GetShippingAddress() visitor.InterfaceVisitorAddress {
	if dc.ShippingAddress == nil {
		return nil
	}

	shippingAddress, err := visitor.GetVisitorAddressModel()
	if err != nil {
		env.ErrorDispatch(err)
		return nil
	}

	err = shippingAddress.FromHashMap(dc.ShippingAddress)
	if err != nil {
		env.ErrorDispatch(err)
		return nil
	}

	return shippingAddress
}

// SetBillingAddress sets billing address for checkout
func (dc *DefaultCheckout) SetBillingAddress(address visitor.InterfaceVisitorAddress) error {
	if address == nil {
		dc.BillingAddress = nil
		return nil
	}

	dc.BillingAddress = address.ToHashMap()
	return nil
}

// GetBillingAddress returns checkout billing address
func (dc *DefaultCheckout) GetBillingAddress() visitor.InterfaceVisitorAddress {
	if dc.BillingAddress == nil {
		return nil
	}

	billingAddress, err := visitor.GetVisitorAddressModel()
	if err != nil {
		env.ErrorDispatch(err)
		return nil
	}

	err = billingAddress.FromHashMap(dc.BillingAddress)
	if err != nil {
		env.ErrorDispatch(err)
		return nil
	}

	return billingAddress
}

// SetPaymentMethod sets payment method for checkout
func (dc *DefaultCheckout) SetPaymentMethod(paymentMethod checkout.InterfacePaymentMethod) error {
	dc.PaymentMethodCode = paymentMethod.GetCode()
	return nil
}

// GetPaymentMethod returns checkout payment method
func (dc *DefaultCheckout) GetPaymentMethod() checkout.InterfacePaymentMethod {
	if paymentMethods := checkout.GetRegisteredPaymentMethods(); paymentMethods != nil {
		for _, paymentMethod := range paymentMethods {
			if paymentMethod.GetCode() == dc.PaymentMethodCode {
				return paymentMethod
			}
		}
	}
	return nil
}

// SetShippingMethod sets payment method for checkout
func (dc *DefaultCheckout) SetShippingMethod(shippingMethod checkout.InterfaceShippingMethod) error {
	dc.ShippingMethodCode = shippingMethod.GetCode()
	return nil
}

// GetShippingMethod returns a checkout shipping method
func (dc *DefaultCheckout) GetShippingMethod() checkout.InterfaceShippingMethod {
	if shippingMethods := checkout.GetRegisteredShippingMethods(); shippingMethods != nil {
		for _, shippingMethod := range shippingMethods {
			if shippingMethod.GetCode() == dc.ShippingMethodCode {
				return shippingMethod
			}
		}
	}
	return nil
}

// SetShippingRate sets shipping rate for checkout
func (dc *DefaultCheckout) SetShippingRate(shippingRate checkout.StructShippingRate) error {
	dc.ShippingRate = shippingRate
	return nil
}

// GetShippingRate returns a checkout shipping rate
func (dc *DefaultCheckout) GetShippingRate() *checkout.StructShippingRate {
	return &dc.ShippingRate
}

// SetCart sets cart for checkout
func (dc *DefaultCheckout) SetCart(checkoutCart cart.InterfaceCart) error {
	dc.CartID = checkoutCart.GetID()
	return nil
}

// GetCart returns a shopping cart
func (dc *DefaultCheckout) GetCart() cart.InterfaceCart {
	cartInstance, _ := cart.LoadCartByID(dc.CartID)
	return cartInstance
}

// SetVisitor sets visitor for checkout
func (dc *DefaultCheckout) SetVisitor(checkoutVisitor visitor.InterfaceVisitor) error {
	dc.VisitorID = checkoutVisitor.GetID()

	if dc.BillingAddress == nil && checkoutVisitor.GetBillingAddress() != nil {
		dc.BillingAddress = checkoutVisitor.GetBillingAddress().ToHashMap()
	}

	if dc.ShippingAddress == nil && checkoutVisitor.GetShippingAddress() != nil {
		dc.ShippingAddress = checkoutVisitor.GetShippingAddress().ToHashMap()
	}

	return nil
}

// GetVisitor return checkout visitor
func (dc *DefaultCheckout) GetVisitor() visitor.InterfaceVisitor {
	if dc.VisitorID == "" {
		return nil
	}

	visitorInstance, err := visitor.LoadVisitorByID(dc.VisitorID)
	if err != nil {
		return nil
	}

	return visitorInstance
}

// SetSession sets visitor for checkout
func (dc *DefaultCheckout) SetSession(checkoutSession api.InterfaceSession) error {
	dc.SessionID = checkoutSession.GetID()
	return nil
}

// GetSession return checkout visitor
func (dc *DefaultCheckout) GetSession() api.InterfaceSession {
	sessionInstance, _ := api.GetSessionByID(dc.SessionID)
	return sessionInstance
}

// GetTaxes collects taxes applied for current checkout
func (dc *DefaultCheckout) GetTaxes() (float64, []checkout.StructTaxRate) {

	var amount float64

	if !dc.taxesCalculateFlag {
		dc.taxesCalculateFlag = true

		dc.Taxes = make([]checkout.StructTaxRate, 0)
		for _, tax := range checkout.GetRegisteredTaxes() {
			for _, taxRate := range tax.CalculateTax(dc) {
				dc.Taxes = append(dc.Taxes, taxRate)
				amount += taxRate.Amount
			}
		}

		dc.taxesCalculateFlag = false
	} else {
		for _, taxRate := range dc.Taxes {
			amount += taxRate.Amount
		}
	}

	return amount, dc.Taxes
}

// GetDiscounts collects discounts applied for current checkout
func (dc *DefaultCheckout) GetDiscounts() (float64, []checkout.StructDiscount) {

	var amount float64

	if !dc.discountsCalculateFlag {
		dc.discountsCalculateFlag = true

		dc.Discounts = make([]checkout.StructDiscount, 0)
		for _, discount := range checkout.GetRegisteredDiscounts() {
			for _, discountValue := range discount.CalculateDiscount(dc) {
				dc.Discounts = append(dc.Discounts, discountValue)
				amount += discountValue.Amount
			}
		}

		dc.discountsCalculateFlag = false
	} else {
		for _, discount := range dc.Discounts {
			amount += discount.Amount
		}
	}

	return amount, dc.Discounts
}

// GetSubtotal returns subtotal total for current checkout
func (dc *DefaultCheckout) GetSubtotal() float64 {
	currentCart := dc.GetCart()
	if currentCart != nil {
		return currentCart.GetSubtotal()
	}

	return 0
}

// GetGrandTotal returns grand total for current checkout: [cart subtotal] + [shipping rate] + [taxes] - [discounts]
func (dc *DefaultCheckout) GetGrandTotal() float64 {
	amount := dc.GetSubtotal()

	if shippingRate := dc.GetShippingRate(); shippingRate != nil {
		amount += shippingRate.Price
	}

	taxAmount, _ := dc.GetTaxes()
	amount += taxAmount

	discountAmount, _ := dc.GetDiscounts()
	amount -= discountAmount

	return amount
}

// SetInfo sets additional info for checkout - any values related to checkout process
func (dc *DefaultCheckout) SetInfo(key string, value interface{}) error {
	if value == nil {
		if _, present := dc.Info[key]; present {
			delete(dc.Info, key)
		}
	} else {
		dc.Info[key] = value
	}

	return nil
}

// GetInfo returns additional checkout info value or nil,
//   - use "*" as a key to get all keys and values currently set
func (dc *DefaultCheckout) GetInfo(key string) interface{} {
	if key == "*" {
		return dc.Info
	}

	if value, present := dc.Info[key]; present {
		return value
	}

	return nil
}

// SetOrder sets order for current checkout
func (dc *DefaultCheckout) SetOrder(checkoutOrder order.InterfaceOrder) error {
	dc.OrderID = checkoutOrder.GetID()
	return nil
}

// GetOrder returns current checkout related order or nil if not created yet
func (dc *DefaultCheckout) GetOrder() order.InterfaceOrder {
	if dc.OrderID != "" {
		orderInstance, err := order.LoadOrderByID(dc.OrderID)
		if err == nil {
			return orderInstance
		}
	}
	return nil
}

// Submit creates the order with provided information
func (dc *DefaultCheckout) Submit() (interface{}, error) {

	if dc.GetBillingAddress() == nil {
		return nil, env.ErrorNew(ConstErrorModule, ConstErrorLevel, "080db3c0-dbb5-4398-b1f1-4c3fefef79b4", "Billing address is not set")
	}

	if dc.GetShippingAddress() == nil {
		return nil, env.ErrorNew(ConstErrorModule, ConstErrorLevel, "1c069d54-2847-46cb-bccd-76fc13d229ea", "Shipping address is not set")
	}

	if dc.GetPaymentMethod() == nil {
		return nil, env.ErrorNew(ConstErrorModule, ConstErrorLevel, "c0628038-3e06-47e9-9252-480351d903c0", "Payment method is not set")
	}

	if dc.GetShippingMethod() == nil {
		return nil, env.ErrorNew(ConstErrorModule, ConstErrorLevel, "e22877fe-248d-4b5e-ad2f-10843cb9890c", "Shipping method is not set")
	}

	currentCart := dc.GetCart()
	if currentCart == nil {
		return nil, env.ErrorNew(ConstErrorModule, ConstErrorLevel, "ea211827-2025-46ca-be47-d09841021890", "Cart is not specified")
	}

	cartItems := currentCart.GetItems()
	if len(cartItems) == 0 {
		return nil, env.ErrorNew(ConstErrorModule, ConstErrorLevel, "84ad4db5-29e9-430c-aecf-b675cbafbcec", "Cart is empty")
	}

	currentVisitor := dc.GetVisitor()
	if (dc.VisitorID == "" || currentVisitor == nil) && utils.InterfaceToString(dc.GetInfo("customer_email")) == "" {
		return nil, env.ErrorNew(ConstErrorModule, ConstErrorLevel, "c5f53ede-63b7-40ea-952d-4d4c04337563", "customer e-mail was not specified")
	}

	if err := currentCart.ValidateCart(); err != nil {
		return nil, env.ErrorDispatch(err)
	}

	// making new order if needed
	//---------------------------
	currentTime := time.Now()

	checkoutOrder := dc.GetOrder()
	if checkoutOrder == nil {
		newOrder, err := order.GetOrderModel()
		if err != nil {
			return nil, env.ErrorDispatch(err)
		}

		newOrder.Set("created_at", currentTime)

		checkoutOrder = newOrder
	}

	// updating order information
	//---------------------------
	checkoutOrder.Set("updated_at", currentTime)

	checkoutOrder.SetStatus(order.ConstOrderStatusNew)

	if currentVisitor != nil {
		checkoutOrder.Set("visitor_id", currentVisitor.GetID())

		checkoutOrder.Set("customer_email", currentVisitor.GetEmail())
		checkoutOrder.Set("customer_name", currentVisitor.GetFullName())
	}

	if dc.GetInfo("customer_email") != nil {
		orderCustomerEmail := utils.InterfaceToString(dc.GetInfo("customer_email"))
		checkoutOrder.Set("customer_email", orderCustomerEmail)
	}
	if dc.GetInfo("customer_name") != nil {
		orderCustomerName := utils.InterfaceToString(dc.GetInfo("customer_name"))
		checkoutOrder.Set("customer_name", orderCustomerName)
	}

	billingAddress := dc.GetBillingAddress().ToHashMap()
	checkoutOrder.Set("billing_address", billingAddress)

	shippingAddress := dc.GetShippingAddress().ToHashMap()
	checkoutOrder.Set("shipping_address", shippingAddress)

	checkoutOrder.Set("cart_id", currentCart.GetID())
	checkoutOrder.Set("payment_method", dc.GetPaymentMethod().GetCode())
	checkoutOrder.Set("shipping_method", dc.GetShippingMethod().GetCode()+"/"+dc.GetShippingRate().Code)

	discountAmount, discounts := dc.GetDiscounts()
	checkoutOrder.Set("discount", discountAmount)
	checkoutOrder.Set("discounts", discounts)

	taxAmount, taxes := dc.GetTaxes()
	checkoutOrder.Set("tax_amount", taxAmount)
	checkoutOrder.Set("taxes", taxes)

	checkoutOrder.Set("shipping_amount", dc.GetShippingRate().Price)

	generateDescriptionFlag := false
	orderDescription := utils.InterfaceToString(dc.GetInfo("order_description"))
	if orderDescription == "" {
		generateDescriptionFlag = true
	}

	for _, cartItem := range cartItems {
		orderItem, err := checkoutOrder.AddItem(cartItem.GetProductID(), cartItem.GetQty(), cartItem.GetOptions())
		if err != nil {
			return nil, env.ErrorDispatch(err)
		}

		if generateDescriptionFlag {
			if orderDescription != "" {
				orderDescription += ", "
			}
			orderDescription += fmt.Sprintf("%dx %s", cartItem.GetQty(), orderItem.GetName())
		}
	}
	checkoutOrder.Set("description", orderDescription)

	err := checkoutOrder.CalculateTotals()
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	err = checkoutOrder.Proceed()
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	dc.SetOrder(checkoutOrder)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	// trying to process payment
	//--------------------------
	if dc.GetGrandTotal() > 0 {
		paymentInfo := make(map[string]interface{})
		paymentInfo["sessionID"] = dc.GetSession().GetID()
		paymentInfo["cc"] = dc.GetInfo("cc")

		result, err := dc.GetPaymentMethod().Authorize(checkoutOrder, paymentInfo)
		if err != nil {
			return nil, env.ErrorDispatch(err)
		}

		// if payment.Authorize returns non nil result, that supposing additional operations to complete payment
		if result != nil {
			return result, nil
		}
	}

	err = dc.CheckoutSuccess(checkoutOrder, dc.GetSession())

	return checkoutOrder.ToHashMap(), err
}
