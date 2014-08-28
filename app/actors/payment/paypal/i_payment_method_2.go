package paypal

import (
	"errors"
	"fmt"

	"io/ioutil"

	"net/http"
	"net/url"

	"github.com/ottemo/foundation/app/models/checkout"
)

func (it *PayPalExpress) GetName() string {
	return PAYMENT_NAME_EXPRESS
}

func (it *PayPalExpress) GetCode() string {
	return PAYMENT_CODE_EXPRESS
}

func (it *PayPalExpress) GetType() string {
	return checkout.PAYMENT_TYPE_REMOTE
}

func (it *PayPalExpress) IsAllowed(checkoutInstance checkout.I_Checkout) bool {
	return true
}

func (it *PayPalExpress) Authorize(checkoutInstance checkout.I_Checkout) error {

	grandTotal := checkoutInstance.GetGrandTotal()
	shippingAmount := checkoutInstance.GetShippingRate().Price

	params := "USER=" + PP_EXPRESS_USER +
		"&PWD=" + PP_EXPRESS_PWD +
		"&SIGNATURE=" + PP_EXPRESS_SIGNATURE +
		"&METHOD=SetExpressCheckout&VERSION=78" +
		"&PAYMENTREQUEST_0_PAYMENTACTION=" + PP_EXPRESS_PAYMENTACTION +
		"&PAYMENTREQUEST_0_AMT=" + fmt.Sprintf("%.2f", grandTotal) +
		"&PAYMENTREQUEST_0_SHIPPINGAMT=" + fmt.Sprintf("%.2f", shippingAmount) +
		"&PAYMENTREQUEST_0_ITEMAMT=" + fmt.Sprintf("%.2f", grandTotal-shippingAmount) +
		"&PAYMENTREQUEST_0_DESC=Purchase%20for%20%24" + fmt.Sprintf("%.2f", grandTotal) +
		"&PAYMENTREQUEST_0_CUSTOM=" + checkoutInstance.GetOrder().GetId() +
		"&PAYMENTREQUEST_0_CURRENCYCODE=USD" +
		"&cancelUrl=http://dev.ottemo.com/paypal/cancel" +
		"&returnUrl=http://dev.ottemo.com/paypal/success"

	// println(params)

	request, err := http.NewRequest("GET", PP_EXPRESS_ENDPOINT+"?"+params, nil)
	if err != nil {
		return err
	}

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return err
	}

	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return err
	}

	// println(string(responseData))

	responseValues, err := url.ParseQuery(string(responseData))
	if err != nil {
		return errors.New("payment unexpected response")
	}

	if responseValues.Get("ACK") != "Success" || responseValues.Get("TOKEN") == "" {
		if responseValues.Get("L_ERRORCODE0") != "" {
			errors.New("payment error " + responseValues.Get("L_ERRORCODE0") + ": " + "L_LONGMESSAGE0")
		}
	}

	checkoutInstance.SetInfo("redirect", PP_EXPRESS_REDIRECT+"&token="+responseValues.Get("TOKEN"))

	return nil
}

func (it *PayPalExpress) Capture(checkoutInstance checkout.I_Checkout) error {
	return nil
}

func (it *PayPalExpress) Refund(checkoutInstance checkout.I_Checkout) error {
	return nil
}

func (it *PayPalExpress) Void(checkoutInstance checkout.I_Checkout) error {
	return nil
}
