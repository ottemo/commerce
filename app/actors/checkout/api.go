package checkout

import (
	"errors"

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
	err = api.GetRestService().RegisterAPI("checkout", "POST", "set/payment/method/:code", restCheckoutSetPaymentMethod)
	if err != nil {
		return err
	}
	err = api.GetRestService().RegisterAPI("checkout", "POST", "set/shipping/method/:code", restCheckoutSetShippingMethod)
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

	result := map[string]interface{} {"billing_address": nil, "shipping_address": nil}



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

	type ResultValue struct{Name string; Code string}
	result := make([]ResultValue, 0)

	for _, shippingMethod := range checkout.GetRegisteredShippingMethods() {
		if shippingMethod.IsAllowed(currentCheckout) {
			result = append(result, ResultValue{Name: shippingMethod.GetName(), Code: shippingMethod.GetCode()})
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

	for _, paymentMethod := range checkout.GetRegisteredPaymentMethods() {
		if paymentMethod.GetCode() ==  params.RequestURLParams["code"] {
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

	for _, shippingMethod := range checkout.GetRegisteredShippingMethods() {
		if shippingMethod.GetCode() ==  params.RequestURLParams["code"] {
			if shippingMethod.IsAllowed(currentCheckout) {

				err := currentCheckout.SetShippingMethod(shippingMethod)
				if err != nil {
					return nil, err
				}

				return "ok", nil
			} else {
				return nil, errors.New("shipping method not allowed")
			}
		}
	}

	return nil, errors.New("shipping method not found")
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

	newOrder.Save()

	return nil, errors.New("shipping method not found")
}
