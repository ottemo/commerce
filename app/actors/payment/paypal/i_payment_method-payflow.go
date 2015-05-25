package paypal

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/ottemo/foundation/app/models/checkout"
	"github.com/ottemo/foundation/app/models/order"
	"github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/utils"
)

func (it *PayFlowAPI) GetName() string {
	return "PayPal Payflow"
}

func (it *PayFlowAPI) GetCode() string {
	return "paypal_payflow"
}

func (it *PayFlowAPI) GetType() string {
	return checkout.ConstPaymentTypeCreditCard
}

func (it *PayFlowAPI) IsAllowed(checkoutInstance checkout.InterfaceCheckout) bool {
	return true
}

func (it *PayFlowAPI) Authorize(orderInstance order.InterfaceOrder, paymentInfo map[string]interface{}) (interface{}, error) {

	//	if !utils.StrKeysInMap(paymentInfo, "type", "number", "expire_month", "expire_year", "cvv") {
	//		return nil, env.ErrorNew(ConstErrorModule, ConstErrorLevel, "ce41da6b-6051-4b3a-883a-4dd37c1c7a1a", "credit card info was not specified")
	//	}

	billingAddress := orderInstance.GetBillingAddress()
	if billingAddress == nil {
		return nil, env.ErrorNew(ConstErrorModule, ConstErrorLevel, "9c57fbf9-cb29-4472-901e-f77a9399dda6", "no billing address information")
	}

	if orderInstance == nil {
		return nil, env.ErrorNew(ConstErrorModule, ConstErrorLevel, "6219944f-d62b-40d7-8495-f844be0e562d", "no created order")
	}

	user := utils.InterfaceToString(env.ConfigGetValue(ConstConfigPathUser))
	password := utils.InterfaceToString(env.ConfigGetValue(ConstConfigPathPass))

	//	templateValues := map[string]string{
	//		// https://developer.paypal.com/webapps/developer/docs/classic/payflow/integration-guide/#paypal-credit-card-transaction-request-parameters
	//		// If you set up one or more additional users on the account, this value is the ID of the user authorized to process transactions.
	//		// If, however, you have not set up additional users on the account, USER has the same value as VENDOR.
	//		"USER":    user,
	//		"TENDER":  "C",      // The method of payment, such as C for credit card
	//		"VENDOR":  user,     // Your merchant login ID that you created when you registered for the account.
	//		"PWD":     password, // The password that you defined while registering for the account.
	//		"PARTNER": "PayPal",
	//
	//		"ACCT":    utils.InterfaceToString("4417119669820331"), // The buyer’s credit card number
	//		"EXPDATE": utils.InterfaceToString(11) + utils.InterfaceToString(18), // The expiration date of the credit card
	//		"CVV2":    utils.InterfaceToString(874),
	//
	//		"BILLTOFIRSTNAME": billingAddress.GetFirstName(),
	//		"BILLTOLASTNAME":  billingAddress.GetLastName(),
	//		"BILLTOSTREET":    billingAddress.GetAddressLine1(),
	//		"BILLTOCITY":      billingAddress.GetCity(),
	//		"BILLTOSTATE":     billingAddress.GetState(),
	//		"BILLTOZIP":       billingAddress.GetZipCode(),
	//		"BILLTOCOUNTRY":   "US", // https://developer.paypal.com/webapps/developer/docs/classic/api/country_codes/
	//		// TODO: Check is it proper values for request
	//
	//		"AMT":      fmt.Sprintf("%.2f", orderInstance.GetGrandTotal()), // The amount of the sale, including two decimal places and without a comma separator
	//		"CURRENCY": "USD",                                              // AUD - Australian dollar, CAD - Canadian dollar, EUR - Euro, GBP - British pound, JPY - Japanese Yen, USD - US dollar
	//
	//		"INVNUM": "order id - " + orderInstance.GetID(),
	//
	//		"CUSTIP": "0.0.0.0", // (Optional) Account holder’s IP address.
	//	}
	requestParams := "USER=" + user +
		"&PWD=" + password +
		"&VENDOR=" + user +
		"&PARTNER=PayPal" +
		"&ACCT=78" +
		"&EXPDATE=" + utils.InterfaceToString(11) + utils.InterfaceToString(18) +
		"&CVV2=" + utils.InterfaceToString(874) +
		"&BILLTOFIRSTNAME=" + billingAddress.GetFirstName() +
		"&BILLTOLASTNAME=" + billingAddress.GetLastName() +
		"&BILLTOSTREET=" + billingAddress.GetAddressLine1() +
		"&BILLTOCITY=" + billingAddress.GetCity() +
		"&BILLTOSTATE=" + billingAddress.GetState() +
		"&BILLTOZIP=" + billingAddress.GetZipCode() +
		"&BILLTOCOUNTRY=US" +
		"&AMT=" + fmt.Sprintf("%.2f", orderInstance.GetGrandTotal()) +
		"&CURRENCY=USD" +
		"&INVNUM=" + orderInstance.GetID()

	fmt.Println(requestParams)

	request, err := http.NewRequest("POST", ConstPaymentPayFlowURL, bytes.NewBufferString(requestParams))
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}
	accessToken, err := it.GetAccessToken()
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	fmt.Println(accessToken)

	request.Header.Add("Content-Type", "text/name value")
	request.Header.Add("Host", "pilot-payflowpro.paypal.com")
	request.Header.Add("X-VPS-CLIENT-TIMEOUT", "45")
	request.Header.Add("X-VPS-REQUEST-ID", "unique-id123123")

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}
	//	TRXTYPE=S&TENDER=C&USER=paypal-facilitator_api1.ottemo.io&PWD=CRR8GX34BMEVMS3B&PARTNER=PayPal&ACCT=5105105105105100&EXPDATE=1215&AMT=0&COMMENT1=Airport+Shuttle&BILLTOFIRSTNAME=Jamie&BILLTOLASTNAME=Miller&BILLTOSTREET=123+Main+St.&BILLTOCITY=San+Jose&BILLTOSTATE=CA&BILLTOZIP=951311234&BILLTOCOUNTRY=840&CVV2=123&CUSTIP=0.0.0.0&VERBOSITY=HIGH&CREATESECURETOKEN=Y&SILENTTRAN=TRUE&SECURETOKENID=9a9ea8208de1413abc3d60c86cb1f4c5

	buf, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	// https://developer.paypal.com/webapps/developer/docs/classic/payflow/integration-guide/#credit-card-transaction-responses
	fmt.Println(string(buf))

	result := make(map[string]interface{})
	err = json.Unmarshal(buf, &result)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}
	if responseMSG, present := result["RESPMSG"]; present {
		fmt.Println("result message: ", utils.InterfaceToString(responseMSG))
	}

	if response.StatusCode != 201 {
		if response.StatusCode == 400 {
			return nil, env.ErrorNew(ConstErrorModule, ConstErrorLevel, "b022da83-61ac-4198-a72d-70f58f65acf0", utils.InterfaceToString(result["details"]))
		}
		return nil, env.ErrorNew(ConstErrorModule, ConstErrorLevel, "741e875d-0ab5-4e8f-81f9-8abbff10cf78", "payment was not processed")
	}

	//TODO: should store information to order

	return nil, nil
}

func (it *PayFlowAPI) Capture(orderInstance order.InterfaceOrder, paymentInfo map[string]interface{}) (interface{}, error) {
	return nil, env.ErrorNew(ConstErrorModule, ConstErrorLevel, "2dc38587-de12-4bdf-9468-a4cef846afe5", "Not implemented")
}

func (it *PayFlowAPI) Refund(orderInstance order.InterfaceOrder, paymentInfo map[string]interface{}) (interface{}, error) {
	return nil, env.ErrorNew(ConstErrorModule, ConstErrorLevel, "2dc38587-de12-4bdf-9468-a4cef846afe5", "Not implemented")
}

func (it *PayFlowAPI) Void(orderInstance order.InterfaceOrder, paymentInfo map[string]interface{}) (interface{}, error) {
	return nil, env.ErrorNew(ConstErrorModule, ConstErrorLevel, "2dc38587-de12-4bdf-9468-a4cef846afe5", "Not implemented")
}

// returns application access token needed for all other requests
func (it *PayFlowAPI) GetAccessToken() (string, error) {

	// TODO:change body
	body := "grant_type=client_credentials"

	request, err := http.NewRequest("POST", "https://api.sandbox.paypal.com/v1/oauth2/token", bytes.NewBufferString(body))
	if err != nil {
		return "", env.ErrorDispatch(err)
	}

	request.SetBasicAuth("AbrcnhDi238ke9aG2NIQqVkW90oMJVg3B1QsjC68d2xRBLDq8boIrCaigPli", "EPcLWBCmfM_AwSOO1jC6TEDLCg-xZhFrUmXQnvTQ9yZV5_786xc5OkQ4Gx2-")

	request.Header.Add("Content-Type", "text/name value")
	request.Header.Add("Host", "pilot-payflowpro.paypal.com")
	request.Header.Add("X-VPS-CLIENT-TIMEOUT", "45")
	request.Header.Add("X-VPS-REQUEST-ID", "unique-id123123")

	resp, err := http.DefaultClient.Do(request)
	if err != nil {
		return "", env.ErrorDispatch(err)
	}

	buf, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", env.ErrorDispatch(err)
	}

	result := make(map[string]interface{})
	err = json.Unmarshal(buf, &result)
	if err != nil {
		return "", env.ErrorDispatch(err)
	}

	// TODO:check response and getting token
	token, present := result["SECURETOKEN"]
	tokenID, ok := result["SECURETOKENID"]

	if present && ok {
		return utils.InterfaceToString(token) + utils.InterfaceToString(tokenID), nil
	}

	return "", env.ErrorNew(ConstErrorModule, ConstErrorLevel, "96d95546-1595-4fe1-9156-ef8d6e43f172", "unexpected response - without access_token")
}
