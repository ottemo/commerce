package authorize

import (
	"bytes"
	"errors"
	"io/ioutil"

	"github.com/ottemo/foundation/api"
	"github.com/ottemo/foundation/app/utils"
)

// startup API registration
func setupAPI() error {

	var err error = nil

	err = api.GetRestService().RegisterAPI("authorizenet", "POST", "receipt", restReceipt)
	if err != nil {
		return err
	}

	return nil
}

// WEB REST API function to process Authorize.Net receipt result
func restReceipt(params *api.T_APIHandlerParams) (interface{}, error) {

	body, err := ioutil.ReadAll(params.Request.Body)
	if err != nil {
		return nil, err
	}

	if bytes.Contains(body, []byte("Thank you for your order!")) {
		currentCheckout, err := utils.GetCurrentCheckout(params)
		if err != nil {
			return nil, err
		}

		checkoutOrder := currentCheckout.GetOrder()
		if checkoutOrder != nil {
			checkoutOrder.NewIncrementId()

			checkoutOrder.Set("status", "pending")

			err = checkoutOrder.Save()
			if err != nil {
				return nil, err
			}

			return checkoutOrder.ToHashMap(), nil
		}
	} else {
		return nil, errors.New("Response error: " + string(body))
	}

	return nil, errors.New("can't process authorize.net response")
}
