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
	"github.com/ottemo/foundation/app/models/order"
	"github.com/ottemo/foundation/api/session"
	"net/http"
	"io/ioutil"
	"net/url"
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

// completes paypal transaction
/**
TOKEN=EC%2d2KK6368096400733J
SUCCESSPAGEREDIRECTREQUESTED=false
TIMESTAMP=2014%2d10%2d13T07%3a35%3a39Z
CORRELATIONID=2fc2f4867c46a
ACK=Success
VERSION=78
BUILD=13298800
INSURANCEOPTIONSELECTED=false
SHIPPINGOPTIONISDEFAULT=false
PAYMENTINFO_0_TRANSACTIONID=4JG56001NB517953X
PAYMENTINFO_0_TRANSACTIONTYPE=expresscheckout
PAYMENTINFO_0_PAYMENTTYPE=instant
PAYMENTINFO_0_ORDERTIME=2014%2d10%2d13T07%3a35%3a39Z
PAYMENTINFO_0_AMT=43%2e00
PAYMENTINFO_0_FEEAMT=1%2e55
PAYMENTINFO_0_TAXAMT=0%2e00
PAYMENTINFO_0_CURRENCYCODE=USD
PAYMENTINFO_0_PAYMENTSTATUS=Completed
PAYMENTINFO_0_PENDINGREASON=None
PAYMENTINFO_0_REASONCODE=None
PAYMENTINFO_0_PROTECTIONELIGIBILITY=Eligible
PAYMENTINFO_0_PROTECTIONELIGIBILITYTYPE=ItemNotReceivedEligible%2cUnauthorizedPaymentEligible
PAYMENTINFO_0_SECUREMERCHANTACCOUNTID=P2J4W2PSKHBKQ
PAYMENTINFO_0_ERRORCODE=0
PAYMENTINFO_0_ACK=Success
 */
func Completes(orderInstance order.I_Order, token string, payerId string) (map[string]string, error) {

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
			"&PAYERID=" + payerId +
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

		completeData, err := Completes(checkoutOrder, reqData["token"], reqData["PayerID"])
		if err != nil {
			env.Log("paypal.log", env.LOG_PREFIX_INFO, "TRANSACTION NOT COMPLETED: "+
						"VisitorId - "+utils.InterfaceToString(checkoutOrder.Get("visitor_id"))+", "+
						"OrderId - "+checkoutOrder.GetId() + ", " +
						"TOKEN - : " + completeData["TOKEN"])

			return nil, errors.New("Transaction not confirmed")
		}

		checkoutOrder.NewIncrementId()

		checkoutOrder.Set("status", "pending shipping")
		checkoutOrder.Set("payment_info", completeData)

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

		env.Log("paypal.log", env.LOG_PREFIX_INFO, "TRANSACTION COMPLETED: "+
					"VisitorId - "+utils.InterfaceToString(checkoutOrder.Get("visitor_id"))+", "+
					"OrderId - "+checkoutOrder.GetId() + ", " +
					"TOKEN - : " + completeData["TOKEN"])

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
