package paypal

import (
	"fmt"

	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/ottemo/foundation/api"
	"github.com/ottemo/foundation/app"
	"github.com/ottemo/foundation/app/models/checkout"
	"github.com/ottemo/foundation/app/models/order"
	"github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/utils"
)

// setupAPI setups package related API endpoint routines
func setupAPI() error {

	var err error

	err = api.GetRestService().RegisterAPI("paypal", "GET", "success", restSuccess)
	if err != nil {
		return err
	}

	err = api.GetRestService().RegisterAPI("paypal", "GET", "cancel", restCancel)
	if err != nil {
		return err
	}

	return nil
}

// Completes will finalizie the PayPal transaction when provided an Order, Token and Payer ID.
func Completes(orderInstance order.InterfaceOrder, token string, payerID string) (map[string]string, error) {

	// getting order information
	//--------------------------
	grandTotal := orderInstance.GetGrandTotal()
	shippingPrice := orderInstance.GetShippingAmount()

	// getting request param values
	//-----------------------------
	user := utils.InterfaceToString(env.ConfigGetValue(ConstConfigPathUser))
	password := utils.InterfaceToString(env.ConfigGetValue(ConstConfigPathPass))
	signature := utils.InterfaceToString(env.ConfigGetValue(ConstConfigPathSignature))
	action := utils.InterfaceToString(env.ConfigGetValue(ConstConfigPathAction))

	amount := fmt.Sprintf("%.2f", grandTotal)
	shippingAmount := fmt.Sprintf("%.2f", shippingPrice)
	itemAmount := fmt.Sprintf("%.2f", grandTotal-shippingPrice)

	description := "Purchase%20for%20%24" + fmt.Sprintf("%.2f", grandTotal)
	custom := orderInstance.GetID()

	// making NVP request
	//-------------------
	requestParams := "USER=" + user +
		"&PWD=" + password +
		"&SIGNATURE=" + signature +
		"&METHOD=DoExpressCheckoutPayment" +
		"&VERSION=78" +
		"&PAYMENTREQUEST_0_PAYMENTACTION=" + action +
		"&PAYMENTREQUEST_0_AMT=" + amount +
		"&PAYMENTREQUEST_0_SHIPPINGAMT=" + shippingAmount +
		"&PAYMENTREQUEST_0_ITEMAMT=" + itemAmount +
		"&PAYMENTREQUEST_0_DESC=" + description +
		"&PAYMENTREQUEST_0_CUSTOM=" + custom +
		"&PAYMENTREQUEST_0_CURRENCYCODE=USD" +
		"&PAYERID=" + payerID +
		"&TOKEN=" + token

	nvpGateway := utils.InterfaceToString(env.ConfigGetValue(ConstConfigPathNVP))

	request, err := http.NewRequest("GET", nvpGateway+"?"+requestParams, nil)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	// reading/decoding response from NVP
	//-----------------------------------
	responseData, err := ioutil.ReadAll(response.Body)

	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	urlGETParams := make(map[string]string)
	urlParsedParams, err := url.ParseQuery(string(responseData))
	if err == nil {
		for key, value := range urlParsedParams {
			urlGETParams[key] = value[0]
		}
	}

	if urlGETParams["ACK"] != "Success" || urlGETParams["TOKEN"] == "" {
		if urlGETParams["L_ERRORCODE0"] != "" {
			return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "2b3f49d1-b6d7-492f-ac96-854a85b64c2f", "payment confirm error "+utils.InterfaceToString(urlGETParams["L_ERRORCODE0"])+": "+"L_LONGMESSAGE0")
		}
	}

	return urlGETParams, nil
}

// WEB REST API function to process PayPal receipt result
func restSuccess(context api.InterfaceApplicationContext) (interface{}, error) {
	requestData := context.GetRequestParameters()
	if !utils.KeysInMapAndNotBlank(requestData, "token", "PayerID") {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "4b0fcede-8ce7-4f4d-bee4-9cf75a427f59", "tocken or payerID are not set")
	}
	sessionID := waitingTokens[requestData["token"]]

	waitingTokensMutex.Lock()
	delete(waitingTokens, requestData["token"])
	waitingTokensMutex.Unlock()

	sessionInstance, err := api.GetSessionByID(utils.InterfaceToString(sessionID))
	if err != nil {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "04eaf220-5964-4d92-821b-c60e1eb8185e", "Wrong session ID")
	}
	err = context.SetSession(sessionInstance)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	currentCheckout, err := checkout.GetCurrentCheckout(context)
	if err != nil {
		return nil, err
	}

	checkoutOrder := currentCheckout.GetOrder()
	currentCart := currentCheckout.GetCart()
	if currentCart == nil {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "15f11fb8-a15e-4788-a4d3-98c41eac4caf", "Cart is not specified")
	}

	if checkoutOrder != nil {

		completeData, err := Completes(checkoutOrder, requestData["token"], requestData["PayerID"])
		if err != nil {
			env.Log("paypal.log", env.ConstLogPrefixInfo, "TRANSACTION NOT COMPLETED: "+
				"VisitorID - "+utils.InterfaceToString(checkoutOrder.Get("visitor_id"))+", "+
				"OrderID - "+checkoutOrder.GetID()+", "+
				"TOKEN - : "+completeData["TOKEN"])

			return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "af2103b1-8501-4e4b-bc1b-9086ef7c63be", "Transaction not confirmed")
		}

		checkoutOrder.NewIncrementID()

		checkoutOrder.Set("status", "pending_shipping")
		checkoutOrder.Set("payment_info", completeData)

		err = currentCheckout.CheckoutSuccess(checkoutOrder, context.GetSession())
		if err != nil {
			return nil, err
		}

		// Send confirmation email
		err = currentCheckout.SendOrderConfirmationMail()
		if err != nil {
			return nil, err
		}

		env.Log("paypal.log", env.ConstLogPrefixInfo, "TRANSACTION COMPLETED: "+
			"VisitorID - "+utils.InterfaceToString(checkoutOrder.Get("visitor_id"))+", "+
			"OrderID - "+checkoutOrder.GetID()+", "+
			"TOKEN - : "+completeData["TOKEN"])

		return api.StructRestRedirect{Location: app.GetStorefrontURL("account/order/" + checkoutOrder.GetID()), DoRedirect: true}, nil
	}

	return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "8d449c1c-ca34-4260-a93b-8af999c1ff04", "Checkout not exist")
}

// WEB REST API function to process PayPal decline result
func restCancel(context api.InterfaceApplicationContext) (interface{}, error) {
	requestData := context.GetRequestParameters()
	if !utils.KeysInMapAndNotBlank(requestData, "token", "PayerID") {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "4b0fcede-8ce7-4f4d-bee4-9cf75a427f59", "tocken or payerID are not set")
	}
	sessionID := waitingTokens[requestData["token"]]

	waitingTokensMutex.Lock()
	delete(waitingTokens, requestData["token"])
	waitingTokensMutex.Unlock()

	session, err := api.GetSessionByID(utils.InterfaceToString(sessionID))
	if err != nil {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "04eaf220-5964-4d92-821b-c60e1eb8185e", "Wrong session ID")
	}
	err = context.SetSession(session)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	currentCheckout, err := checkout.GetCurrentCheckout(context)
	if err != nil {
		return nil, err
	}

	checkoutOrder := currentCheckout.GetOrder()

	env.Log("paypal.log", env.ConstLogPrefixInfo, "CANCELED: "+
		"VisitorID - "+utils.InterfaceToString(checkoutOrder.Get("visitor_id"))+", "+
		"OrderID - "+checkoutOrder.GetID()+", "+
		"TOKEN - : "+requestData["token"])

	return api.StructRestRedirect{Location: app.GetStorefrontURL("checkout"), DoRedirect: true}, nil
}
