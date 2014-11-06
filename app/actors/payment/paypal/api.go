package paypal

import (
	"errors"
	"fmt"

	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/ottemo/foundation/api"
	"github.com/ottemo/foundation/api/session"
	"github.com/ottemo/foundation/app"
	"github.com/ottemo/foundation/app/models/checkout"
	"github.com/ottemo/foundation/app/models/order"
	"github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/utils"
)

// startup API registration
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
func Completes(orderInstance order.I_Order, token string, payerID string) (map[string]string, error) {

	// getting order information
	//--------------------------
	grandTotal := orderInstance.GetGrandTotal()
	shippingPrice := orderInstance.GetShippingAmount()

	// getting request param values
	//-----------------------------
	user := utils.InterfaceToString(env.ConfigGetValue(CONFIG_PATH_USER))
	password := utils.InterfaceToString(env.ConfigGetValue(CONFIG_PATH_PASS))
	signature := utils.InterfaceToString(env.ConfigGetValue(CONFIG_PATH_SIGNATURE))
	action := utils.InterfaceToString(env.ConfigGetValue(CONFIG_PATH_ACTION))

	amount := fmt.Sprintf("%.2f", grandTotal)
	shippingAmount := fmt.Sprintf("%.2f", shippingPrice)
	itemAmount := fmt.Sprintf("%.2f", grandTotal-shippingPrice)

	description := "Purchase%20for%20%24" + fmt.Sprintf("%.2f", grandTotal)
	custom := orderInstance.GetId()

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

	nvpGateway := utils.InterfaceToString(env.ConfigGetValue(CONFIG_PATH_NVP))

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
			return nil, env.ErrorNew("payment confirm error " + utils.InterfaceToString(urlGETParams["L_ERRORCODE0"]) + ": " + "L_LONGMESSAGE0")
		}
	}

	return urlGETParams, nil
}

func restSuccess(params *api.T_APIHandlerParams) (interface{}, error) {
	reqData := params.RequestGETParams
	sessionID := waitingTokens[reqData["token"]]

	waitingTokensMutex.Lock()
	delete(waitingTokens, reqData["token"])
	waitingTokensMutex.Unlock()

	session, err := session.GetSessionById(utils.InterfaceToString(sessionID))
	if err != nil {
		return nil, errors.New("Wrong session ID")
	}
	params.Session = session

	currentCheckout, err := checkout.GetCurrentCheckout(params)
	if err != nil {
		return nil, err
	}

	checkoutOrder := currentCheckout.GetOrder()
	currentCart := currentCheckout.GetCart()
	if currentCart == nil {
		return nil, errors.New("Cart is not specified")
	}

	if checkoutOrder != nil {

		completeData, err := Completes(checkoutOrder, reqData["token"], reqData["PayerID"])
		if err != nil {
			env.Log("paypal.log", env.LOG_PREFIX_INFO, "TRANSACTION NOT COMPLETED: "+
				"VisitorId - "+utils.InterfaceToString(checkoutOrder.Get("visitor_id"))+", "+
				"OrderId - "+checkoutOrder.GetId()+", "+
				"TOKEN - : "+completeData["TOKEN"])

			return nil, errors.New("Transaction not confirmed")
		}

		checkoutOrder.NewIncrementId()

		checkoutOrder.Set("status", "pending_shipping")
		checkoutOrder.Set("payment_info", completeData)

		err = currentCheckout.CheckoutSuccess(checkoutOrder, params.Session)
		if err != nil {
			return nil, err
		}

		// Send confirmation email
		err = currentCheckout.SendOrderConfirmationMail()
		if err != nil {
			return nil, err
		}

		env.Log("paypal.log", env.LOG_PREFIX_INFO, "TRANSACTION COMPLETED: "+
			"VisitorId - "+utils.InterfaceToString(checkoutOrder.Get("visitor_id"))+", "+
			"OrderId - "+checkoutOrder.GetId()+", "+
			"TOKEN - : "+completeData["TOKEN"])

		return api.T_RestRedirect{Location: app.GetStorefrontUrl("account/order/" + checkoutOrder.GetId()), DoRedirect: true}, nil
	}

	return nil, errors.New("Checkout not exist")
}

// WEB REST API function to process Paypal.com cancel result
func restCancel(params *api.T_APIHandlerParams) (interface{}, error) {
	reqData := params.RequestGETParams
	sessionID := waitingTokens[reqData["token"]]

	waitingTokensMutex.Lock()
	delete(waitingTokens, reqData["token"])
	waitingTokensMutex.Unlock()

	session, err := session.GetSessionById(utils.InterfaceToString(sessionID))
	if err != nil {
		return nil, errors.New("Wrong session ID")
	}
	params.Session = session

	currentCheckout, err := checkout.GetCurrentCheckout(params)
	if err != nil {
		return nil, err
	}

	checkoutOrder := currentCheckout.GetOrder()

	env.Log("paypal.log", env.LOG_PREFIX_INFO, "CANCELED: "+
		"VisitorId - "+utils.InterfaceToString(checkoutOrder.Get("visitor_id"))+", "+
		"OrderId - "+checkoutOrder.GetId()+", "+
		"TOKEN - : "+reqData["token"])

	return api.T_RestRedirect{Location: app.GetStorefrontUrl("checkout"), DoRedirect: true}, nil
}
