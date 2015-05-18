package checkout

import (
	"fmt"
	"github.com/ottemo/foundation/api"
	"github.com/ottemo/foundation/app"
	"github.com/ottemo/foundation/app/models/cart"
	"github.com/ottemo/foundation/app/models/checkout"
	"github.com/ottemo/foundation/app/models/order"
	"github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/media"
	"github.com/ottemo/foundation/utils"
	"io/ioutil"
	"net/http"
	"strings"
)

// SendOrderConfirmationMail sends an order confirmation email
func (it *DefaultCheckout) SendOrderConfirmationMail() error {

	checkoutOrder := it.GetOrder()
	if checkoutOrder == nil {
		return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "e7c69056-cc28-4632-9524-50d71b909d83", "given checkout order does not exists")
	}

	confirmationEmail := utils.InterfaceToString(env.ConfigGetValue(checkout.ConstConfigPathConfirmationEmail))
	if confirmationEmail != "" {
		email := utils.InterfaceToString(checkoutOrder.Get("customer_email"))
		if email == "" {
			return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "1202fcfb-da3f-4a0f-9a2e-92f288fd3881", "customer email for order is not set")
		}

		visitorMap := make(map[string]interface{})
		if visitorModel := it.GetVisitor(); visitorModel != nil {
			visitorMap = visitorModel.ToHashMap()
		} else {
			visitorMap["first_name"] = checkoutOrder.Get("customer_name")
			visitorMap["email"] = checkoutOrder.Get("customer_email")
		}

		confirmationEmail, err := utils.TextTemplate(confirmationEmail,
			map[string]interface{}{
				"Order":   checkoutOrder.ToHashMap(),
				"Visitor": visitorMap,
			})
		if err != nil {
			return env.ErrorDispatch(err)
		}

		err = app.SendMail(email, "Order confirmation", confirmationEmail)
		if err != nil {
			return env.ErrorDispatch(err)
		}
	}

	return nil
}

// CheckoutSuccess will save the order and clear the shopping in the session.
func (it *DefaultCheckout) CheckoutSuccess(checkoutOrder order.InterfaceOrder, session api.InterfaceSession) error {

	// making sure order and session were specified
	if checkoutOrder == nil || session == nil {
		return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "17d45365-7808-4a1b-ad36-1741a83e820f", "Order or session is null")
	}

	// if payment method did not set status by itself - making this
	if orderStatus := checkoutOrder.GetStatus(); orderStatus == "" || orderStatus == order.ConstOrderStatusNew {
		err := checkoutOrder.SetStatus(order.ConstOrderStatusPending)
		if err != nil {
			return env.ErrorDispatch(err)
		}
	}

	// checkout information cleanup
	//-----------------------------
	currentCart := it.GetCart()

	err := currentCart.Deactivate()
	if err != nil {
		return env.ErrorDispatch(err)
	}

	err = currentCart.Save()
	if err != nil {
		return env.ErrorDispatch(err)
	}

	session.Set(cart.ConstSessionKeyCurrentCart, nil)
	session.Set(checkout.ConstSessionKeyCurrentCheckout, nil)

	// sending notifications
	//----------------------
	eventData := map[string]interface{}{"checkout": it, "order": checkoutOrder, "session": session, "cart": currentCart}
	env.Event("checkout.success", eventData)

	// make asynchronously a request to TrustPilot
	go func(checkoutOrder order.InterfaceOrder, currentCart cart.InterfaceCart) {
		cartItems := currentCart.GetItems()

		requestData := make(map[string]interface{})

		requestData["consumer"] = map[string]interface{}{
			"email": "test" + utils.InterfaceToString(checkoutOrder.Get("customer_email")),
			"name":  checkoutOrder.Get("customer_name"),
		}
		requestData["referenceId"] = checkoutOrder.GetID()
		requestData["locale"] = "en-US"
		requestData["products"] = make([]map[string]string, 0)

		// product media removal
		mediaStorage, err := media.GetMediaStorage()
		if err != nil {
			env.LogError(err)
		}

		for _, productItem := range cartItems {
			currentProductID := productItem.GetProductID()
			currentProduct := productItem.GetProduct()

			mediaPath, err := mediaStorage.GetMediaPath("product", currentProductID, "image")
			if err != nil {
				env.ErrorDispatch(err)
			}
			productOptions := productItem.GetOptions()
			productBrand := ""
			if brand, present := productOptions["brand"]; present {
				productBrand = utils.InterfaceToString(brand)
			}

			productInfo := map[string]string {
				"productUrl": app.GetStorefrontURL("product/" + currentProductID),
				"imageUrl":   app.GetStorefrontURL(mediaPath + currentProduct.GetDefaultImage()),
				"name":       currentProduct.GetName(),
				"sku":        currentProduct.GetSku(),
				"brand":      productBrand,
			}

			requestData["products"] = append(utils.InterfaceToArray(requestData["products"]), productInfo)
		}
		fmt.Println(requestData, currentCart)

		requestURL := "https://api.trustpilot.com/v1/private/product-reviews/business-units/{businessUnitId}/invitation-links"
		businessUnitID := "test"
		requestURL = strings.Replace(requestURL, "{businessUnitId}", businessUnitID, 1)

		request, err := http.NewRequest("POST", requestURL, nil)
		request.Header.Set("Content-Type", "application/json")
		if err != nil {
			env.LogError(err)
		}

		client := &http.Client{}
		response, err := client.Do(request)

		defer response.Body.Close()
		if err != nil {
			env.LogError(err)
		}
		responseBody, err := ioutil.ReadAll(response.Body)
		if err != nil {
			env.LogError(err)
		}

		jsonResponse, err := utils.DecodeJSONToStringKeyMap(responseBody)
		if err != nil {
			env.LogError(err)
		}

		if result, ok := jsonResponse["result"]; ok {
			fmt.Println(result)
		}

	}(checkoutOrder)

	err = it.SendOrderConfirmationMail()
	if err != nil {
		env.ErrorDispatch(err)
	}

	return nil
}
