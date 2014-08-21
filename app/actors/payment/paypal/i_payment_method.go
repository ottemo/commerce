package paypal

import (
	"bytes"
	"errors"

	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/ottemo/foundation/app/models/checkout"
	"github.com/ottemo/foundation/app/utils"
)

func (it *PayPal) GetName() string {
	return PAYMENT_NAME
}

func (it *PayPal) GetCode() string {
	return PAYMENT_CODE
}

func (it *PayPal) GetType() string {
	return checkout.PAYMENT_TYPE_CREDIT_CARD
}

func (it *PayPal) IsAllowed(checkoutInstance checkout.I_Checkout) bool {
	return true
}

func (it *PayPal) Authorize(checkoutInstance checkout.I_Checkout) error {
	payment := make(map[string]interface{})
	payment["intent"] = "authorize"

	payer := make(map[string]interface{})
	payment["payer"] = payer
	payer["payment_method"] = "credit_card"

	fundingInstruments := make(map[string]interface{})
	payer["funding_instruments"] = fundingInstruments

	creditCard := make(map[string]interface{})
	fundingInstruments["credit_card"] = creditCard

	creditCard["number"] = "4417119669820331"
	creditCard["type"] = "visa"
	creditCard["expire_month"] = 11
	creditCard["expire_year"] = 2018
	creditCard["cvv2"] = "874"
	creditCard["first_name"] = "Betsy"
	creditCard["last_name"] = "Buyer"

	billingAddress := make(map[string]interface{})
	creditCard["billing_address"] = billingAddress

	billingAddress["line1"] = "111 First Street"
	billingAddress["city"] = "Saratoga"
	billingAddress["state"] = "CA"
	billingAddress["postal_code"] = "95070"
	billingAddress["country_code"] = "US"


	body, err := utils.EncodeToJsonString(payment)
	if err != nil {
		return err
	}

	request, err := http.NewRequest("POST", "https://api.sandbox.paypal.com/v1/oauth2/token", bytes.NewBufferString(body))
	if err != nil {
		return err
	}

	accessToken, err := it.GetAccessToken(checkoutInstance)
	if err != nil {
		return err
	}

	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("Accept", "application/json")
	request.Header.Add("Authorization", "Bearer " + accessToken)

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return err
	}

	buf, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return err
	}

	result := make(map[string]interface{})
	err = json.Unmarshal(buf, &result)
	if err != nil {
		return err
	}

	//TODO: should store information to order

	return nil
}

func (it *PayPal) Capture(checkoutInstance checkout.I_Checkout) error {
	return nil
}

func (it *PayPal) Refund(checkoutInstance checkout.I_Checkout) error {
	return nil
}

func (it *PayPal) Void(checkoutInstance checkout.I_Checkout) error {
	return nil
}

// returns application access token needed for all other requests
func (it *PayPal) GetAccessToken(checkoutInstance checkout.I_Checkout) (string, error) {
	body := "grant_type=client_credentials"
	req, err := http.NewRequest("POST", "https://api.sandbox.paypal.com/v1/oauth2/token", bytes.NewBufferString(body))
	if err != nil {
		return "", err
	}

	req.SetBasicAuth("AbrcnhDi238ke9aG2NIQqVkW90oMJVg3B1QsjC68d2xRBLDq8boIrCaigPli", "EPcLWBCmfM_AwSOO1jC6TEDLCg-xZhFrUmXQnvTQ9yZV5_786xc5OkQ4Gx2-")

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Accept-Language", "en_US")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}

	buf, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	result := make(map[string]interface{})
	err = json.Unmarshal(buf, &result)
	if err != nil {
		return "", err
	}

	if token, present := result["access_token"]; present {
		return  utils.InterfaceToString(token), nil
	}

	return "", errors.New("unexpected response - without access_token")
}
