package paypal

import (
	"errors"
	"fmt"

	"github.com/ottemo/foundation/api"
	"github.com/ottemo/foundation/app"
	"github.com/ottemo/foundation/utils"
	"github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/app/models/checkout"
	"github.com/ottemo/foundation/app/models/cart"
	"github.com/ottemo/foundation/api/session"
	"net/http"
	"io/ioutil"
)

// startup API registration
func setupAPI() error {

	var err error = nil

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

func restSuccess(params *api.T_APIHandlerParams) (interface{}, error) {
	reqData := params.RequestGETParams
	sessionId := waitingTokens[reqData["token"]]

	waitingTokensMutex.Lock()
	delete(waitingTokens, reqData["token"])
	waitingTokensMutex.Unlock()

	session, err := session.GetSessionById(utils.InterfaceToString(sessionId))
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

		///-------------------------///
		grandTotal := checkoutOrder.GetGrandTotal()
		shippingPrice := checkoutOrder.GetShippingAmount()

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
		custom := checkoutOrder.GetId()

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
				"&PAYERID=" + reqData["PayerID"] +
				"&TOKEN=" + reqData["token"]

			println(requestParams)

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
		fmt.Println(responseData)
		fmt.Println(reqData)
		///-------------------------///

		checkoutOrder.NewIncrementId()

		checkoutOrder.Set("status", "pending")
		checkoutOrder.Set("payment_info", reqData)

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

		// Send confirmation email
		err = currentCheckout.SendOrderConfirmationMail()
		if err != nil {
			return nil, err
		}

		env.Log("paypal.log", env.LOG_PREFIX_INFO, "TRANSACTION APPROVED: "+
					"VisitorId - "+utils.InterfaceToString(checkoutOrder.Get("visitor_id"))+", "+
					"OrderId - "+checkoutOrder.GetId() + ", " +
					"TOKEN - : " + reqData["token"])

		return api.T_RestRedirect{Location: app.GetStorefrontUrl("account/order/" + checkoutOrder.GetId()), DoRedirect: true}, nil
	}

	return nil, errors.New("Checkout not exist")
}

// WEB REST API function to process Paypal.com cancel result
func restCancel(params *api.T_APIHandlerParams) (interface{}, error) {
	reqData := params.RequestGETParams
	sessionId := waitingTokens[reqData["token"]]

	waitingTokensMutex.Lock()
	delete(waitingTokens, reqData["token"])
	waitingTokensMutex.Unlock()

	session, err := session.GetSessionById(utils.InterfaceToString(sessionId))
	if err != nil {
		return nil, errors.New("Wrong session ID")
	}
	params.Session = session

	currentCheckout, err := checkout.GetCurrentCheckout(params)
	if err != nil {
		return nil, err
	}

	checkoutOrder := currentCheckout.GetOrder()

	env.Log("paypal.log", env.LOG_PREFIX_INFO, "CANCELED: " +
				"VisitorId - "+utils.InterfaceToString(checkoutOrder.Get("visitor_id"))+", "+
				"OrderId - "+checkoutOrder.GetId() + ", " +
				"TOKEN - : " + reqData["token"])

	return api.T_RestRedirect{Location: app.GetStorefrontUrl("checkout"), DoRedirect: true}, nil
}
