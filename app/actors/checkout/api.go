package checkout

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/ottemo/foundation/api"
	"github.com/ottemo/foundation/app/models/cart"
	"github.com/ottemo/foundation/app/models/checkout"
	"github.com/ottemo/foundation/app/models/order"
	"github.com/ottemo/foundation/app/models/visitor"

	"github.com/ottemo/foundation/app/utils"
)

func setupAPI() error {

	var err error = nil

	err = api.GetRestService().RegisterAPI("checkout", "GET", "info", restCheckoutInfo)
	if err != nil {
		return err
	}
	err = api.GetRestService().RegisterAPI("checkout", "GET", "payment/methods", restCheckoutPaymentMethods)
	if err != nil {
		return err
	}
	err = api.GetRestService().RegisterAPI("checkout", "GET", "shipping/methods", restCheckoutShippingMethods)
	if err != nil {
		return err
	}
	err = api.GetRestService().RegisterAPI("checkout", "POST", "set/shipping/address", restCheckoutSetShippingAddress)
	if err != nil {
		return err
	}
	err = api.GetRestService().RegisterAPI("checkout", "POST", "set/billing/address", restCheckoutSetBillingAddress)
	if err != nil {
		return err
	}
	err = api.GetRestService().RegisterAPI("checkout", "POST", "set/payment/method/:method", restCheckoutSetPaymentMethod)
	if err != nil {
		return err
	}
	err = api.GetRestService().RegisterAPI("checkout", "POST", "set/shipping/method/:method/:rate", restCheckoutSetShippingMethod)
	if err != nil {
		return err
	}
	err = api.GetRestService().RegisterAPI("checkout", "GET", "submit", restSubmit)
	if err != nil {
		return err
	}
	err = api.GetRestService().RegisterAPI("checkout", "POST", "submit", restSubmit)
	if err != nil {
		return err
	}

	return nil
}

// WEB REST API function to get current checkout process status
func restCheckoutInfo(params *api.T_APIHandlerParams) (interface{}, error) {

	currentCheckout, err := utils.GetCurrentCheckout(params)
	if err != nil {
		return nil, err
	}

	result := map[string]interface{}{
		"billing_address":  nil,
		"shipping_address": nil,

		"payment_method_name": nil,
		"payment_method_code": nil,

		"shipping_method_name": nil,
		"shipping_method_code": nil,

		"shipping_rate":   nil,
		"shipping_amount": nil,

		"discounts":       nil,
		"discount_amount": nil,

		"taxes":      nil,
		"tax_amount": nil,

		"subtotal":   nil,
		"grandtotal": nil,
	}

	if billingAddress := currentCheckout.GetBillingAddress(); billingAddress != nil {
		result["billing_address"] = billingAddress.ToHashMap()
	}

	if shippingAddress := currentCheckout.GetShippingAddress(); shippingAddress != nil {
		result["shipping_address"] = shippingAddress.ToHashMap()
	}

	if paymentMethod := currentCheckout.GetPaymentMethod(); paymentMethod != nil {
		result["payment_method_name"] = paymentMethod.GetName()
		result["payment_method_code"] = paymentMethod.GetCode()
	}

	if shippingMethod := currentCheckout.GetShippingMethod(); shippingMethod != nil {
		result["shipping_method_name"] = shippingMethod.GetName()
		result["shipping_method_code"] = shippingMethod.GetCode()
	}

	if shippingRate := currentCheckout.GetShippingRate(); shippingRate != nil {
		result["shipping_rate"] = shippingRate
		result["shipping_amount"] = shippingRate.Price
	}

	if checkoutCart := currentCheckout.GetCart(); checkoutCart != nil {
		result["subtotal"] = checkoutCart.GetSubtotal()
	}

	result["discount_amount"], result["discounts"] = currentCheckout.GetDiscounts()

	result["tax_amount"], result["taxes"] = currentCheckout.GetTaxes()

	result["grandtotal"] = currentCheckout.GetGrandTotal()

	return result, nil
}

// WEB REST API function to get possible payment methods for checkout
func restCheckoutPaymentMethods(params *api.T_APIHandlerParams) (interface{}, error) {

	currentCheckout, err := utils.GetCurrentCheckout(params)
	if err != nil {
		return nil, err
	}

	type ResultValue struct {
		Name string
		Code string
		Type string
	}
	result := make([]ResultValue, 0)

	for _, paymentMethod := range checkout.GetRegisteredPaymentMethods() {
		if paymentMethod.IsAllowed(currentCheckout) {
			result = append(result, ResultValue{Name: paymentMethod.GetName(), Code: paymentMethod.GetCode(), Type: paymentMethod.GetType()})
		}
	}

	return result, nil
}

// WEB REST API function to get possible shipping methods for checkout
func restCheckoutShippingMethods(params *api.T_APIHandlerParams) (interface{}, error) {

	currentCheckout, err := utils.GetCurrentCheckout(params)
	if err != nil {
		return nil, err
	}

	type ResultValue struct {
		Name  string
		Code  string
		Rates []checkout.T_ShippingRate
	}
	result := make([]ResultValue, 0)

	for _, shippingMethod := range checkout.GetRegisteredShippingMethods() {
		if shippingMethod.IsAllowed(currentCheckout) {
			result = append(result, ResultValue{Name: shippingMethod.GetName(), Code: shippingMethod.GetCode(), Rates: shippingMethod.GetRates(currentCheckout)})
		}
	}

	return result, nil
}

// internal function for  restCheckoutSetShippingAddress() and restCheckoutSetBillingAddress()
func checkoutObtainAddress(params *api.T_APIHandlerParams) (visitor.I_VisitorAddress, error) {

	reqData, err := api.GetRequestContentAsMap(params)
	if err != nil {
		return nil, err
	}

	if addressId, present := reqData["id"]; present {

		// Address id was specified - trying to load
		visitorAddress, err := visitor.LoadVisitorAddressById(utils.InterfaceToString(addressId))
		if err != nil {
			return nil, err
		}

		currentVisitorId := utils.InterfaceToString(params.Session.Get(visitor.SESSION_KEY_VISITOR_ID))
		if visitorAddress.GetVisitorId() != currentVisitorId {
			return nil, errors.New("wrong address id")
		}

		return visitorAddress, nil
	} else {

		// supposedly address data was specified
		visitorAddressModel, err := visitor.GetVisitorAddressModel()
		if err != nil {
			return nil, err
		}

		for attribute, value := range reqData {
			err := visitorAddressModel.Set(attribute, value)
			if err != nil {
				return nil, err
			}
		}

		visitorId := utils.InterfaceToString(params.Session.Get(visitor.SESSION_KEY_VISITOR_ID))
		visitorAddressModel.Set("visitor_id", visitorId)

		err = visitorAddressModel.Save()
		if err != nil {
			return nil, err
		}

		return visitorAddressModel, nil
	}
}

// WEB REST API function to set shipping address
func restCheckoutSetShippingAddress(params *api.T_APIHandlerParams) (interface{}, error) {
	currentCheckout, err := utils.GetCurrentCheckout(params)
	if err != nil {
		return nil, err
	}

	address, err := checkoutObtainAddress(params)
	if err != nil {
		return nil, err
	}

	err = currentCheckout.SetShippingAddress(address)
	if err != nil {
		return nil, err
	}

	return address.ToHashMap(), nil
}

// WEB REST API function to set billing address
func restCheckoutSetBillingAddress(params *api.T_APIHandlerParams) (interface{}, error) {
	currentCheckout, err := utils.GetCurrentCheckout(params)
	if err != nil {
		return nil, err
	}

	address, err := checkoutObtainAddress(params)
	if err != nil {
		return nil, err
	}

	err = currentCheckout.SetBillingAddress(address)
	if err != nil {
		return nil, err
	}

	return address.ToHashMap(), nil
}

// WEB REST API function to set payment method
func restCheckoutSetPaymentMethod(params *api.T_APIHandlerParams) (interface{}, error) {
	currentCheckout, err := utils.GetCurrentCheckout(params)
	if err != nil {
		return nil, err
	}

	// looking for payment method
	for _, paymentMethod := range checkout.GetRegisteredPaymentMethods() {
		if paymentMethod.GetCode() == params.RequestURLParams["method"] {
			if paymentMethod.IsAllowed(currentCheckout) {

				// updating checkout payment method
				err := currentCheckout.SetPaymentMethod(paymentMethod)
				if err != nil {
					return nil, err
				}

				// checking for additional info
				contentValues, _ := api.GetRequestContentAsMap(params)
				for key, value := range contentValues {
					currentCheckout.SetInfo(key, value)
				}

				return "ok", nil
			} else {
				return nil, errors.New("payment method not allowed")
			}
		}
	}

	return nil, errors.New("payment method not found")
}

// WEB REST API function to set payment method
func restCheckoutSetShippingMethod(params *api.T_APIHandlerParams) (interface{}, error) {
	currentCheckout, err := utils.GetCurrentCheckout(params)
	if err != nil {
		return nil, err
	}

	// looking for shipping method
	for _, shippingMethod := range checkout.GetRegisteredShippingMethods() {
		if shippingMethod.GetCode() == params.RequestURLParams["method"] {
			if shippingMethod.IsAllowed(currentCheckout) {

				// looking for shipping rate
				for _, shippingRate := range shippingMethod.GetRates(currentCheckout) {
					if shippingRate.Code == params.RequestURLParams["rate"] {

						err := currentCheckout.SetShippingMethod(shippingMethod)
						if err != nil {
							return nil, err
						}

						err = currentCheckout.SetShippingRate(shippingRate)
						if err != nil {
							return nil, err
						}

						return "ok", nil
					}
				}

			} else {
				return nil, errors.New("shipping method not allowed")
			}
		}
	}

	return nil, errors.New("shipping method and/or rate were not found")
}

// WEB REST API function to submit checkout information and make order
func restSubmit(params *api.T_APIHandlerParams) (interface{}, error) {

	// TODO: should be splited on smaller functions
	// TODO: order for checkout perhaps should be associated with cart

	// checking for needed information
	//--------------------------------
	currentCheckout, err := utils.GetCurrentCheckout(params)
	if err != nil {
		return nil, err
	}

	if currentCheckout.GetBillingAddress() == nil {
		return nil, errors.New("Billing address is not set")
	}

	if currentCheckout.GetShippingAddress() == nil {
		return nil, errors.New("Shipping address is not set")
	}

	if currentCheckout.GetPaymentMethod() == nil {
		return nil, errors.New("Payment method is not set")
	}

	if currentCheckout.GetShippingMethod() == nil {
		return nil, errors.New("Shipping method is not set")
	}

	currentCart := currentCheckout.GetCart()
	if currentCart == nil {
		return nil, errors.New("Cart is not specified")
	}

	cartItems := currentCart.GetItems()
	if len(cartItems) == 0 {
		return nil, errors.New("Cart have no products inside")
	}

	// checking for additional info
	//-----------------------------
	if params.Request.Method == "POST" {
		contentValues, _ := api.GetRequestContentAsMap(params)
		for key, value := range contentValues {
			currentCheckout.SetInfo(key, value)
		}
	}

	// making new order if needed
	//---------------------------
	currentTime := time.Now()

	checkoutOrder := currentCheckout.GetOrder()
	if checkoutOrder == nil {
		newOrder, err := order.GetOrderModel()
		if err != nil {
			return nil, err
		}

		newOrder.Set("created_at", currentTime)

		checkoutOrder = newOrder
	}

	// updating order information
	//---------------------------
	checkoutOrder.Set("updated_at", currentTime)

	checkoutOrder.Set("status", "new")
	if currentVisitor := currentCheckout.GetVisitor(); currentVisitor != nil {
		checkoutOrder.Set("visitor_id", currentVisitor.GetId())

		checkoutOrder.Set("customer_email", currentVisitor.GetEmail())
		checkoutOrder.Set("customer_name", currentVisitor.GetFullName())
	}

	billingAddress := currentCheckout.GetBillingAddress().ToHashMap()
	checkoutOrder.Set("billing_address", billingAddress)

	shippingAddress := currentCheckout.GetShippingAddress().ToHashMap()
	checkoutOrder.Set("shipping_address", shippingAddress)

	checkoutOrder.Set("cart_id", currentCart.GetId())
	checkoutOrder.Set("payment_method", currentCheckout.GetPaymentMethod().GetCode())
	checkoutOrder.Set("shipping_method", currentCheckout.GetShippingMethod().GetCode()+"/"+currentCheckout.GetShippingRate().Code)

	discountAmount, _ := currentCheckout.GetDiscounts()
	taxAmount, _ := currentCheckout.GetTaxes()

	checkoutOrder.Set("discount", discountAmount)
	checkoutOrder.Set("tax_amount", taxAmount)
	checkoutOrder.Set("shipping_amount", currentCheckout.GetShippingRate().Price)

	generateDescriptionFlag := false
	orderDescription := utils.InterfaceToString(currentCheckout.GetInfo("order_description"))
	if orderDescription == "" {
		generateDescriptionFlag = true
	}

	for _, cartItem := range cartItems {
		orderItem, err := checkoutOrder.AddItem(cartItem.GetProductId(), cartItem.GetQty(), cartItem.GetOptions())
		if err != nil {
			return nil, err
		}

		product := cartItem.GetProduct()
		if product == nil {
			return nil, errors.New("no product for cart item " + strconv.Itoa(cartItem.GetIdx()))
		}

		orderItem.Set("name", product.GetName())
		orderItem.Set("sku", product.GetSku())
		orderItem.Set("short_description", product.GetShortDescription())

		orderItem.Set("price", product.GetPrice())
		orderItem.Set("size", product.GetSize())
		orderItem.Set("weight", product.GetWeight())

		if generateDescriptionFlag {
			if orderDescription != "" {
				orderDescription += ", "
			}
			orderDescription += fmt.Sprintf("%dx %s", cartItem.GetQty(), product.GetName())
		}
	}
	checkoutOrder.Set("description", orderDescription)

	err = checkoutOrder.CalculateTotals()
	if err != nil {
		return nil, err
	}

	err = checkoutOrder.Save()
	if err != nil {
		return nil, err
	}

	currentCheckout.SetOrder(checkoutOrder)
	if err != nil {
		return nil, err
	}

	// trying to process payment
	//--------------------------
	err = currentCheckout.GetPaymentMethod().Authorize(currentCheckout)
	if err != nil {
		return nil, err
	}

	if currentCheckout.GetPaymentMethod().GetType() == checkout.PAYMENT_TYPE_REMOTE {
		if currentCheckout.GetInfo(checkout.CHECKOUT_INFO_KEY_REDIRECT) != nil {
			return api.T_RestRedirect{
				Result:   "redirect",
				Location: utils.InterfaceToString(currentCheckout.GetInfo(checkout.CHECKOUT_INFO_KEY_REDIRECT)),
			}, nil
		}
	}

	// assigning new order increment id after success payment
	//-------------------------------------------------------
	checkoutOrder.NewIncrementId()

	checkoutOrder.Set("status", "pending")

	err = checkoutOrder.Save()
	if err != nil {
		return nil, err
	}

	// cleanup checkout information
	//-----------------------------
	currentCart.Deactivate()
	currentCart.Save()
	params.Session.Set(cart.SESSION_KEY_CURRENT_CART, nil)

	params.Session.Set(checkout.SESSION_KEY_CURRENT_CHECKOUT, nil)

	return checkoutOrder.ToHashMap(), nil
}
