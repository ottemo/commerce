package authorize

import (
	"bytes"

	"io/ioutil"

	"github.com/ottemo/foundation/api"
	"github.com/ottemo/foundation/app/models/checkout"
	"github.com/ottemo/foundation/env"
)

// startup API registration
func setupAPI() error {

	var err error = nil

	err = api.GetRestService().RegisterAPI("authorizenet", "POST", "receipt", restReceipt)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	return nil
}

// WEB REST API function to process Authorize.Net receipt result
func restReceipt(params *api.T_APIHandlerParams) (interface{}, error) {

	body, err := ioutil.ReadAll(params.Request.Body)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	if bytes.Contains(body, []byte("Thank you for your order!")) {
		currentCheckout, err := checkout.GetCurrentCheckout(params)
		if err != nil {
			return nil, env.ErrorDispatch(err)
		}

		checkoutOrder := currentCheckout.GetOrder()
		if checkoutOrder != nil {
			checkoutOrder.NewIncrementId()

			checkoutOrder.Set("status", "pending")

			err = checkoutOrder.Save()
			if err != nil {
				return nil, env.ErrorDispatch(err)
			}

			return checkoutOrder.ToHashMap(), nil
		}
	} else {
		return nil, env.ErrorNew("Response error: " + string(body))
	}

	return nil, env.ErrorNew("can't process authorize.net response")
}
