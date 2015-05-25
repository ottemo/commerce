package trustpilot

import (
	"github.com/ottemo/foundation/app"
	"github.com/ottemo/foundation/app/models/cart"
	"github.com/ottemo/foundation/app/models/order"
	"github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/media"
	"github.com/ottemo/foundation/utils"

	"bytes"
	"encoding/base64"
	"io/ioutil"
	"net/http"
	"strings"
)

// checkoutSuccessHandler on checkout Success make request to Trust Pilot with order information
func checkoutSuccessHandler(event string, eventData map[string]interface{}) bool {
	go sendOrderInfo(eventData["order"].(order.InterfaceOrder), eventData["cart"].(cart.InterfaceCart))

	return true
}

// make asynchronously a request to TrustPilot
func sendOrderInfo(checkoutOrder order.InterfaceOrder, currentCart cart.InterfaceCart) error {
	if utils.InterfaceToBool(env.ConfigGetValue(ConstConfigPathTrustPilotEnabled)) {
		// get all config values of trust pilot settings to variables
		trustPilotAPIKey := utils.InterfaceToString(env.ConfigGetValue(ConstConfigPathTrustPilotAPIKey))
		trustPilotAPISecret := utils.InterfaceToString(env.ConfigGetValue(ConstConfigPathTrustPilotAPISecret))
		trustPilotBusinessUnitID := utils.InterfaceToString(env.ConfigGetValue(ConstConfigPathTrustPilotBusinessUnitID))
		trustPilotUsername := utils.InterfaceToString(env.ConfigGetValue(ConstConfigPathTrustPilotUsername))
		trustPilotPassword := utils.InterfaceToString(env.ConfigGetValue(ConstConfigPathTrustPilotPassword))
		trustPilotAccessTokenURL := utils.InterfaceToString(env.ConfigGetValue(ConstConfigPathTrustPilotAccessTokenURL))
		trustPilotProductReviewURL := utils.InterfaceToString(env.ConfigGetValue(ConstConfigPathTrustPilotProductReviewURL))

		// check config data for presence
		if trustPilotAPIKey != "" && trustPilotAPISecret != "" && trustPilotBusinessUnitID != "" && trustPilotBusinessUnitID != "" &&
			trustPilotUsername != "" && trustPilotPassword != "" && trustPilotAccessTokenURL != "" && trustPilotProductReviewURL != "" {
			currentToken := make(map[string]interface{})

			// do request to get token for authentication
			// create request with all necessary info
			bodyString := "grant_type=password&username=" + trustPilotUsername + "&password=" + trustPilotPassword
			buffer := bytes.NewBuffer([]byte(bodyString))

			valueAMIKeySecret := []byte(trustPilotAPIKey + ":" + trustPilotAPISecret)
			encodedString := base64.StdEncoding.EncodeToString(valueAMIKeySecret)

			request, err := http.NewRequest("POST", trustPilotAccessTokenURL, buffer)
			request.Header.Set("Authorization", "Basic "+encodedString)
			request.Header.Set("Content-Type", "application/x-www-form-urlencoded")

			if err != nil {
				env.LogError(err)
				return nil
			}

			client := &http.Client{}
			response, err := client.Do(request)
			if err != nil {
				env.LogError(err)
				return nil
			}

			defer response.Body.Close()

			responseBody, err := ioutil.ReadAll(response.Body)
			if err != nil {
				env.LogError(err)
				return nil
			}

			jsonResponse, err := utils.DecodeJSONToStringKeyMap(responseBody)
			if err != nil {
				env.LogError(err)
				return nil
			}

			if len(jsonResponse) > 0 {
				currentToken = jsonResponse
			}

			if accessToken, present := currentToken["access_token"]; present {
				cartItems := currentCart.GetItems()

				requestData := make(map[string]interface{})
				customerEmail := utils.InterfaceToString(checkoutOrder.Get("customer_email"))

				testModeOn := true
				if testModeOn {
					customerEmail = strings.Replace(customerEmail, "@", "test@", 1)
				}

				requestData["consumer"] = map[string]interface{}{
					"email": customerEmail,
					"name":  checkoutOrder.Get("customer_name"),
				}

				requestData["referenceId"] = checkoutOrder.GetID()
				requestData["locale"] = "en-US"

				mediaStorage, err := media.GetMediaStorage()
				if err != nil {
					env.LogError(err)
					return nil
				}

				var productsOrdered []map[string]string

				// get all products data to request
				for _, productItem := range cartItems {
					currentProductID := productItem.GetProductID()
					currentProduct := productItem.GetProduct()

					mediaPath, err := mediaStorage.GetMediaPath("product", currentProductID, "image")
					if err != nil {
						env.LogError(err)
						return nil
					}

					productOptions := productItem.GetOptions()
					productBrand := ""
					if brand, present := productOptions["brand"]; present {
						productBrand = utils.InterfaceToString(brand)
					}

					productInfo := map[string]string{
						"productUrl": app.GetStorefrontURL("product/" + currentProductID),
						"imageUrl":   app.GetStorefrontURL(mediaPath + currentProduct.GetDefaultImage()),
						"name":       currentProduct.GetName(),
						"sku":        currentProduct.GetSku(),
						"brand":      productBrand,
					}

					productsOrdered = append(productsOrdered, productInfo)
				}

				requestData["products"] = productsOrdered

				trustPilotProductReviewURL = strings.Replace(trustPilotProductReviewURL, "{businessUnitId}", trustPilotBusinessUnitID, 1)

				jsonString := utils.EncodeToJSONString(requestData)
				buffer := bytes.NewBuffer([]byte(jsonString))

				request, err := http.NewRequest("POST", trustPilotProductReviewURL, buffer)
				request.Header.Set("Content-Type", "application/json")
				request.Header.Set("Authorization", "Bearer "+utils.InterfaceToString(accessToken))

				if err != nil {
					env.LogError(err)
					return nil
				}

				client := &http.Client{}
				response, err := client.Do(request)
				if err != nil {
					env.LogError(err)
					return nil
				}

				defer response.Body.Close()

				responseBody, err := ioutil.ReadAll(response.Body)
				if err != nil {
					env.LogError(err)
					return nil
				}

				jsonResponse, err := utils.DecodeJSONToStringKeyMap(responseBody)
				if err != nil {
					env.LogError(err)
					return nil
				}
				if _, ok := jsonResponse["reviewUrl"]; !ok {
					if errorMessage, present := jsonResponse["message"]; present {
						env.LogError(env.ErrorNew(ConstErrorModule, env.ConstErrorLevelActor, "c53fd02f-2f5d-4111-8318-69a2cc2d2259", "Review link empty, error message: "+utils.InterfaceToString(errorMessage)))
					} else {
						env.LogError(env.ErrorNew(ConstErrorModule, env.ConstErrorLevelActor, "c53fd02f-2f5d-4111-8318-69a2cc2d2259", "Review link empty, no error message"))
					}
				}

			} else {
				env.LogError(env.ErrorNew(ConstErrorModule, env.ConstErrorLevelActor, "1293708d-9638-455a-8d49-3a387f086181", "access token is empty"))
			}
		} else {
			env.LogError(env.ErrorNew(ConstErrorModule, env.ConstErrorLevelActor, "22207d49-e001-4666-8501-26bf5ef0926b", "some of trust pilot settings are blank"))
		}
	}
	return nil
}
