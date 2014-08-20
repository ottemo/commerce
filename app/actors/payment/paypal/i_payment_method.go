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

func (it *PayPal) IsAllowed(checkout checkout.I_Checkout) bool {
	return true
}

func (it *PayPal) Authorize() error {
	// apiUser := "paypal_api1.ottemo.io"
	// apiPassword := "1407821638"
	// apiSignature := "AFcWxV21C7fd0v3bYYYRCpSSRl31AqosoWhBSaGs-CU45dQ.JdNevqah"

	return nil
}

func (it *PayPal) Capture() error {
	return nil
}

func (it *PayPal) Refund() error {
	return nil
}

func (it *PayPal) Void() error {
	return nil
}

// returns application access token needed for all other requests
func GetAccessToken() (string, error) {
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
		return utils.InterfaceToString(token), nil
	}

	return "", errors.New("unexpected response - without access_token")
}
