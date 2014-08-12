package checkout

import (
	"errors"
	"time"

	"github.com/ottemo/foundation/api"
	"github.com/ottemo/foundation/app/models/checkout"
	"github.com/ottemo/foundation/app/models/visitor"
	"github.com/ottemo/foundation/app/models/order"
	"github.com/ottemo/foundation/app/utils"
	"strconv"
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
	err = api.GetRestService().RegisterAPI("checkout", "POST", "set/payment/method/:method/:rate", restCheckoutSetPaymentMethod)
	if err != nil {
		return err
	}
	err = api.GetRestService().RegisterAPI("checkout", "POST", "set/shipping/method/:method", restCheckoutSetShippingMethod)
	if err != nil {
		return err
	}
	err = api.GetRestService().RegisterAPI("checkout", "GET", "submit", restSubmit)
	if err != nil {
		return err
	}

	return nil
}



// returns cart for current session
func getCurrentCheckout(params *api.T_APIHandlerParams) (checkout.I_Checkout, error) {
	sessionObject := params.Session.Get(checkout.CHECKOUT_MODEL_NAME)

	if checkoutInstance, ok := sessionObject.(checkout.I_Checkout); ok {
		return checkoutInstance, nil
	} else {
		newCheckoutInstance, err := checkout.GetCheckoutModel()
		if err != nil {
			return nil, err
		}

		params.Session.Set(checkout.CHECKOUT_MODEL_NAME, newCheckoutInstance)

		return newCheckoutInstance, nil
	}
}



// WEB REST API function to get current checkout process status
func restCheckoutInfo(params *api.T_APIHandlerParams) (interface{}, error) {

	currentCheckout, err := getCurrentCheckout(params)
	if err != nil {
		return nil, err
	}

	result := map[string]interface{} {
		"billing_address": nil,
		"shipping_address": nil,

		"payment_method_name": nil,
		"payment_method_code": nil,

		"shipping_method_name": nil,
		"shipping_method_code": nil,

		"shipping_rate": nil,
	}



	if billingAddress := currentCheckout.GetBillingAddress();  billingAddress != nil {
		result["billing_address"] = billingAddress.ToHashMap()
	}


	if shippingAddress := currentCheckout.GetShippingAddress();  shippingAddress != nil {
		result["shipping_address"] = shippingAddress.ToHashMap()
	}


	if paymentMethod := currentCheckout.GetPaymentMethod();  paymentMethod != nil {
		result["payment_method_name"] = paymentMethod.GetName()
		result["payment_method_code"] = paymentMethod.GetCode()
	}

	if shippingMethod := currentCheckout.GetShippingMethod();  shippingMethod != nil {
		result["shipping_method_name"] = shippingMethod.GetName()
		result["shipping_method_code"] = shippingMethod.GetCode()
	}

	if shippingRate := currentCheckout.GetShippingRate();  shippingRate != nil {
		result["shipping_rate"] = shippingRate
	}

	return result, nil
}



// WEB REST API function to get possible payment methods for checkout
func restCheckoutPaymentMethods(params *api.T_APIHandlerParams) (interface{}, error) {

	currentCheckout, err := getCurrentCheckout(params)
	if err != nil {
		return nil, err
	}

	type ResultValue struct{Name string; Code string}
	result := make([]ResultValue, 0)

	for _, paymentMethod := range checkout.GetRegisteredPaymentMethods() {
		if paymentMethod.IsAllowed(currentCheckout) {
			result = append(result, ResultValue{Name: paymentMethod.GetName(), Code: paymentMethod.GetCode()})
		}
	}

	return result, nil
}




// WEB REST API function to get possible shipping methods for checkout
func restCheckoutShippingMethods(params *api.T_APIHandlerParams) (interface{}, error) {

	currentCheckout, err := getCurrentCheckout(params)
	if err != nil {
		return nil, err
	}

	type ResultValue struct{Name string; Code string; Rates []checkout.T_ShippingRate}
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
		visitorAddress, err := visitor.LoadVisitorAddressById( utils.InterfaceToString(addressId) )
		if err != nil {
			return nil, err
		}

		currentVisitorId := utils.InterfaceToString( params.Session.Get( visitor.SESSION_KEY_VISITOR_ID ) )
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

		err = visitorAddressModel.Save()
		if err != nil {
			return nil, err
		}

		return visitorAddressModel, nil
	}
}



// WEB REST API function to set shipping address
func restCheckoutSetShippingAddress(params *api.T_APIHandlerParams) (interface{}, error) {
	currentCheckout, err := getCurrentCheckout(params)
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
	currentCheckout, err := getCurrentCheckout(params)
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
	currentCheckout, err := getCurrentCheckout(params)
	if err != nil {
		return nil, err
	}

	// looking for mayment method
	for _, paymentMethod := range checkout.GetRegisteredPaymentMethods() {
		if paymentMethod.GetCode() == params.RequestURLParams["method"] {
			if paymentMethod.IsAllowed(currentCheckout) {

				err := currentCheckout.SetPaymentMethod(paymentMethod)
				if err != nil {
					return nil, err
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
	currentCheckout, err := getCurrentCheckout(params)
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
	currentCheckout, err := getCurrentCheckout(params)
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

	cartItems := currentCart.ListItems()
	if len(cartItems) == 0 {
		return nil, errors.New("Cart have no products inside")
	}

	newOrder, err := order.GetOrderModel()
	if err != nil {
		return nil ,err
	}



	// newOrder.Set("increment_id",)

	currentTime := time.Now()
	newOrder.Set("created_at", currentTime)
	newOrder.Set("created_at", currentTime)

	newOrder.Set("status", "pending")
	if currentVisitor := currentCheckout.GetVisitor(); currentVisitor != nil {
		newOrder.Set("visitor_id", currentVisitor.GetId())

		newOrder.Set("customer_email", currentVisitor.GetEmail() )
		newOrder.Set("customer_name", currentVisitor.GetFullName() )
	}

	newOrder.Set("cart_id", currentCart.GetId() )
	newOrder.Set("payment_method", currentCheckout.GetPaymentMethod().GetCode() )
	newOrder.Set("shipping_method", currentCheckout.GetShippingMethod().GetCode() + "/" + currentCheckout.GetShippingRate().Code)

	var discountAmount float64 = 0.0
	for _, discount := range currentCheckout.GetDiscounts() {
		discountAmount += discount.Amount
	}
	var taxAmount float64 = 0.0
	for _, taxRate := range currentCheckout.GetTaxes() {
		taxAmount += taxRate.Amount
	}

	newOrder.Set("discount", discountAmount)
	newOrder.Set("tax_amount", taxAmount)
	newOrder.Set("shipping_amount", currentCheckout.GetShippingRate().Price)

	for _, cartItem := range cartItems {
		orderItem, err := newOrder.AddItem(cartItem.GetProductId(), cartItem.GetQty(), cartItem.GetOptions())
		if err != nil {
			return nil, err
		}

		product := cartItem.GetProduct()
		if product == nil {
			return nil, errors.New("no product for cart item " + strconv.Itoa(cartItem.GetIdx()))
		}

		orderItem.Set("name",  product.GetName())
		orderItem.Set("sku",   product.GetSku())
		orderItem.Set("short_description", product.GetShortDescription())

		orderItem.Set("price", product.GetPrice())
		orderItem.Set("size",  product.GetSize())
		orderItem.Set("weight",product.GetWeight())
	}

	newOrder.CalculateTotals()
	newOrder.Save()

	return nil, errors.New("shipping method not found")
}
