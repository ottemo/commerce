package paypal

import (
	"errors"
	"fmt"

	"io/ioutil"

	"net/http"
	"net/url"

	"github.com/ottemo/foundation/api"
	"github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/app/utils"

	"github.com/ottemo/foundation/app/models/order"
	"github.com/ottemo/foundation/app/models/checkout"
)

// returns payment method name
func (it *PayPalExpress) GetName() string {
	return utils.InterfaceToString(env.ConfigGetValue(CONFIG_PATH_TITLE))
}

// returns payment method code
func (it *PayPalExpress) GetCode() string {
	return PAYMENT_CODE
}

// returns type of payment method
func (it *PayPalExpress) GetType() string {
	return checkout.PAYMENT_TYPE_REMOTE
}

// checks for method applicability
func (it *PayPalExpress) IsAllowed(checkoutInstance checkout.I_Checkout) bool {
	return utils.InterfaceToBool(env.ConfigGetValue(CONFIG_PATH_ENABLED))
}

// makes payment method authorize operation
func (it *PayPalExpress) Authorize(orderInstance order.I_Order, paymentInfo map[string]interface{}) (interface{}, error) {

	// getting order information
	//--------------------------
	grandTotal := orderInstance.GetGrandTotal()
	shippingPrice := orderInstance.GetShippingAmount()


	// getting request param values
	//-----------------------------
	user := utils.InterfaceToString( env.ConfigGetValue(CONFIG_PATH_USER) )
	password := utils.InterfaceToString( env.ConfigGetValue(CONFIG_PATH_PASS) )
	signature := utils.InterfaceToString( env.ConfigGetValue(CONFIG_PATH_SIGNATURE) )
	action := utils.InterfaceToString( env.ConfigGetValue(CONFIG_PATH_ACTION) )

	amount := fmt.Sprintf("%.2f", grandTotal)
	shippingAmount := fmt.Sprintf("%.2f", shippingPrice)
	itemAmount := fmt.Sprintf("%.2f", grandTotal - shippingPrice)

	description := "Purchase%20for%20%24" + fmt.Sprintf("%.2f", grandTotal)
	custom := orderInstance.GetId()

	cancelUrl := "http://dev.ottemo.com:3000/paypal/cancel"
	returnUrl := "http://dev.ottemo.com:3000/paypal/success"


	// making NVP request
	//-------------------
	requestParams := "USER=" + user +
		"&PWD=" + password +
		"&SIGNATURE=" + signature +
		"&METHOD=SetExpressCheckout" +
		"&VERSION=78" +
		"&PAYMENTREQUEST_0_PAYMENTACTION=" + action +
		"&PAYMENTREQUEST_0_AMT=" + amount +
		"&PAYMENTREQUEST_0_SHIPPINGAMT=" + shippingAmount +
		"&PAYMENTREQUEST_0_ITEMAMT=" + itemAmount +
		"&PAYMENTREQUEST_0_DESC=" + description +
		"&PAYMENTREQUEST_0_CUSTOM=" + custom +
		"&PAYMENTREQUEST_0_CURRENCYCODE=USD" +
		"&cancelUrl=" + cancelUrl +
		"&returnUrl=" + returnUrl

	// println(requestParams)

	nvpGateway := utils.InterfaceToString( env.ConfigGetValue(CONFIG_PATH_NVP) )

	request, err := http.NewRequest("GET", nvpGateway + "?" + requestParams, nil)
	if err != nil {
		return nil, err
	}

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return nil, err
	}

	// reading/decoding response from NVP
	//-----------------------------------
	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	// println(string(responseData))

	responseValues, err := url.ParseQuery(string(responseData))
	if err != nil {
		return nil, errors.New("payment unexpected response")
	}

	if responseValues.Get("ACK") != "Success" || responseValues.Get("TOKEN") == "" {
		if responseValues.Get("L_ERRORCODE0") != "" {
			return nil, errors.New("payment error " + responseValues.Get("L_ERRORCODE0") + ": " + "L_LONGMESSAGE0")
		}
	}

	// redirecting user to PayPal server for following checkout
	//---------------------------------------------------------
	redirectGateway := utils.InterfaceToString( env.ConfigGetValue(CONFIG_PATH_GATEWAY) ) + "&token=" + responseValues.Get("TOKEN")
	return api.T_RestRedirect {
				Result:   "redirect",
				Location: redirectGateway,
			}, nil
}

// makes payment method capture operation
func (it *PayPalExpress) Capture(orderInstance order.I_Order, paymentInfo map[string]interface{}) (interface{}, error) {
	return nil, errors.New("Not implemented")
}

// makes payment method refund operation
func (it *PayPalExpress) Refund(orderInstance order.I_Order, paymentInfo map[string]interface{}) (interface{}, error) {
	return nil, errors.New("Not implemented")
}

// makes payment method void operation
func (it *PayPalExpress) Void(orderInstance order.I_Order, paymentInfo map[string]interface{}) (interface{}, error) {
	return nil, errors.New("Not implemented")
}
