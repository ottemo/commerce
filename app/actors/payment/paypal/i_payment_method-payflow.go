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
	"time"
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

// Authorize	TRXTYPE=S&TENDER=C&USER=paypal-facilitator_api1.ottemo.io&PWD=CRR8GX34BMEVMS3B&PARTNER=PayPal&ACCT=5105105105105100&EXPDATE=1215&AMT=0&COMMENT1=Airport+Shuttle&BILLTOFIRSTNAME=Jamie&BILLTOLASTNAME=Miller&BILLTOSTREET=123+Main+St.&BILLTOCITY=San+Jose&BILLTOSTATE=CA&BILLTOZIP=951311234&BILLTOCOUNTRY=840&CVV2=123&CUSTIP=0.0.0.0&VERBOSITY=HIGH&CREATESECURETOKEN=Y&SILENTTRAN=TRUE&SECURETOKENID=9a9ea8208de1413abc3d60c86cb1f4c5
// TRXTYPE=S&TENDER=C&USER=xom94ok&PWD=KOL78963&PARTNER=PayPal&ACCT=5105105105105100&EXPDATE=1215&AMT=10&COMMENT1=Airport+Shuttle&BILLTOFIRSTNAME=Jamie&BILLTOLASTNAME=Miller&BILLTOSTREET=123+Main+St.&BILLTOCITY=San+Jose&BILLTOSTATE=CA&BILLTOZIP=951311234&BILLTOCOUNTRY=840&&SILENTTRAN=TRUE&SECURETOKENID=1a9ea8208de1413abc3d60c86cb1f4c5&SECURETOKEN=9eZ53i5WnH0WMTC6Ur6oobQUb

func (it *PayFlowAPI) Authorize(orderInstance order.InterfaceOrder, paymentInfo map[string]interface{}) (interface{}, error) {

	if ccInfo, present := paymentInfo["cc"]; !present || !utils.StrKeysInMap(utils.InterfaceToMap(ccInfo), "type", "number", "expire_month", "expire_year") {
		return nil, env.ErrorNew(ConstErrorModule, ConstErrorLevel, "0aff2b85-0db7-4b38-a7d7-d30faee88357", "credit card info was not specified")
	}

	ccInfo := utils.InterfaceToMap(paymentInfo["cc"])

	// getting order information
	//--------------------------
	grandTotal := orderInstance.GetGrandTotal()
	shippingPrice := orderInstance.GetShippingAmount()
	taxesPrice := orderInstance.GetTaxAmount()
	discountAmount := orderInstance.GetDiscountAmount()
	discount :=fmt.Sprintf("%.2f", discountAmount)
	amount := fmt.Sprintf("%.2f", grandTotal)
	shippingAmount := fmt.Sprintf("%.2f", shippingPrice)
	taxesAmount := fmt.Sprintf("%.2f", taxesPrice)
	itemAmount := fmt.Sprintf("%.2f", grandTotal-shippingPrice-taxesPrice+discountAmount)

	billingAddress := orderInstance.GetBillingAddress()
	if billingAddress == nil {
		return nil, env.ErrorNew(ConstErrorModule, ConstErrorLevel, "9c57fbf9-cb29-4472-901e-f77a9399dda6", "no billing address information")
	}

	if orderInstance == nil {
		return nil, env.ErrorNew(ConstErrorModule, ConstErrorLevel, "6219944f-d62b-40d7-8495-f844be0e562d", "no created order")
	}

	user := utils.InterfaceToString(env.ConfigGetValue(ConstConfigPathPayPalPayflowUser))
	password := utils.InterfaceToString(env.ConfigGetValue(ConstConfigPathPayPalPayflowPass))
	vendor := utils.InterfaceToString(env.ConfigGetValue(ConstConfigPathPayPalPayflowVendor))

	accessTokenInfo, err := it.GetAccessToken()
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	// PayFlow Request Fields
	requestParams := "USER=" + user +
			"&PWD=" + password +
			"&VENDOR=" + vendor +
			"&PARTNER=PayPal" +
			"&VERSION=122" +
			"&TRXTYPE=S" + // Sale

			// Credit Card Details Fields
			"&TENDER=C" +
			"&ACCT=" + utils.InterfaceToString(ccInfo["number"]) +
			"&EXPDATE=" + utils.InterfaceToString(ccInfo["expire_month"]) + utils.InterfaceToString(ccInfo["expire_year"]) +

			// Payer Information Fields
			"&EMAIL=" + utils.InterfaceToString(orderInstance.Get("customer_email")) +
			"&BILLTOFIRSTNAME=" + billingAddress.GetFirstName() +
			"&BILLTOLASTNAME=" + billingAddress.GetLastName() +

			// Payment Details Fields
			"&AMT=" + amount +
			"&CURRENCY=USD" +
			"&ITEMAMT=" + itemAmount +
			"&FREIGHTAMT=" + shippingAmount +
			"&TAXAMT=" + taxesAmount +
			"&DISCOUNT=" + discount +

			"&INVNUM=" + orderInstance.GetID()+
			"&" + accessTokenInfo

	fmt.Println(requestParams)

	nvpGateway := utils.InterfaceToString(env.ConfigGetValue(ConstConfigPathPayPalPayflowURL))
	request, err := http.NewRequest("POST", nvpGateway, bytes.NewBufferString(requestParams))
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	request.Header.Add("Content-Type", "text/name value")
	request.Header.Add("Host", utils.InterfaceToString(env.ConfigGetValue(ConstConfigPathPayPalPayflowHost)))

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	defer response.Body.Close()

	buf, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	// https://developer.paypal.com/webapps/developer/docs/classic/payflow/integration-guide/#credit-card-transaction-responses
	fmt.Println(string(buf))
	// PNREF - Transaction ID

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

	// TRXTYPE=S&AMT=0&USER=xom94ok&PWD=KOL78963&PARTNER=PayPal&SILENTTRAN=TRUE&SECURETOKENID=123448308de1413abc3d60c86cb1f4c5&CREATESECURETOKEN=Y
	user := utils.InterfaceToString(env.ConfigGetValue(ConstConfigPathPayPalPayflowUser))
	password := utils.InterfaceToString(env.ConfigGetValue(ConstConfigPathPayPalPayflowPass))
	vendor := utils.InterfaceToString(env.ConfigGetValue(ConstConfigPathPayPalPayflowVendor))

	secureTokenID := utils.InterfaceToString(time.Now().UnixNano())
	// making NVP request
	//-------------------.
	// PayFlow Request Fields
	requestParams := "USER=" + user +
			"&PWD=" + password +
			"&VENDOR=" + vendor +
			"&PARTNER=PayPal" +
			"&VERSION=122" +
			"&TRXTYPE=S" + // Sale
			"&CREATESECURETOKEN=Y" +
			"&SILENTTRAN=TRUE"+
			"&SECURETOKENID=" + secureTokenID

	nvpGateway := utils.InterfaceToString(env.ConfigGetValue(ConstConfigPathPayPalPayflowURL))

	request, err := http.NewRequest("POST", nvpGateway, bytes.NewBufferString(requestParams))
	if err != nil {
		return "", env.ErrorDispatch(err)
	}

	request.Header.Add("Content-Type", "text/name value")
	request.Header.Add("Host", utils.InterfaceToString(env.ConfigGetValue(ConstConfigPathPayPalPayflowHost)))

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

	token, present := result["SECURETOKEN"]
	tokenID, ok := result["SECURETOKENID"]
	if tokenID != secureTokenID {
		return "", env.ErrorNew(ConstErrorModule, ConstErrorLevel, "9b095f62-b371-4eaf-965f-98eb24206e53", "unexpected response, SECURETOKENID value changed")
	}

	if present && ok {
		return "SECURETOKEN=" + utils.InterfaceToString(token) +"&SECURETOKENID="+ utils.InterfaceToString(tokenID), nil
	}

	return "", env.ErrorNew(ConstErrorModule, ConstErrorLevel, "96d95546-1595-4fe1-9156-ef8d6e43f172", "unexpected response - without access_token")
}
