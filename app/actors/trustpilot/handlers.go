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
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

// checkoutSuccessHandler on checkout Success make request to Trust Pilot with order information
func checkoutSuccessHandler(event string, eventData map[string]interface{}) bool {

	// make asynchronously a request to TrustPilot
	go func(checkoutOrder order.InterfaceOrder, currentCart cart.InterfaceCart) {
		if utils.InterfaceToBool(env.ConfigGetValue(ConstConfigPathTrustPilotEnabled)) {
			// get all config values of trust pilot settings to variables
			trustPilotAPIKey := utils.InterfaceToString(env.ConfigGetValue(ConstConfigPathTrustPilotAPIKey))
			trustPilotAPISecret := utils.InterfaceToString(env.ConfigGetValue(ConstConfigPathTrustPilotAPISecret))
			trustPilotBusinessUnitID := utils.InterfaceToString(env.ConfigGetValue(ConstConfigPathTrustPilotBusinessUnitID))
			trustPilotUsername := utils.InterfaceToString(env.ConfigGetValue(ConstConfigPathTrustPilotUsername))
			trustPilotPassword := utils.InterfaceToString(env.ConfigGetValue(ConstConfigPathTrustPilotPassword))

			// check config data for presence
			if trustPilotAPIKey != "" && trustPilotAPISecret != "" && trustPilotBusinessUnitID != "" &&
				trustPilotBusinessUnitID != "" && trustPilotUsername != "" && trustPilotPassword != "" {
				currentToken := new(Token)

				// do request to get token for authentication
				if currentToken.Access == "" || currentToken.Expiration <= time.Now().Unix() {
					// create request with all necessary info
					requestURL := "https://api.trustpilot.com/v1/oauth/oauth-business-users-for-applications/accesstoken"
					fmt.Println("send request to: ", requestURL)

					bodyString := "grant_type=password&username=" + trustPilotUsername + "&password=" + trustPilotPassword
					buffer := bytes.NewBuffer([]byte(bodyString))
					fmt.Println(bodyString)

					valueAPIKeySacret := []byte(trustPilotAPIKey + ":" + trustPilotAPISecret)
					encodedString := base64.StdEncoding.EncodeToString(valueAPIKeySacret)
					fmt.Println(encodedString)

					request, err := http.NewRequest("POST", requestURL, buffer)
					request.Header.Set("Authorization", "Basic "+encodedString)
					request.Header.Set("Content-Type", "application/x-www-form-urlencoded")

					if err != nil {
						env.LogError(err)
					}

					client := &http.Client{}
					response, err := client.Do(request)
					if err != nil {
						env.LogError(err)
					}

					defer response.Body.Close()

					responseBody, err := ioutil.ReadAll(response.Body)
					if err != nil {
						env.LogError(err)
					}

					jsonResponse, err := utils.DecodeJSONToStringKeyMap(responseBody)
					if err != nil {
						env.LogError(err)
					}

					fmt.Println("token get response: ", jsonResponse)
					if result, ok := jsonResponse["result"]; ok {
						fmt.Println("result get: ", result)
						tokenData := utils.InterfaceToMap(result)
						currentToken.Access = utils.InterfaceToString(tokenData["access_token"])
						currentToken.Refresh = utils.InterfaceToString(tokenData["refresh_token"])
						currentToken.Expiration = time.Now().Unix() + tokenData["expires_in"].(int64)
					}
					fmt.Println("end function request one")
				}

				if currentToken.Access != "" && currentToken.Expiration > time.Now().Unix() {
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

					requestData["products"] = make([]map[string]string, 0)

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

						productInfo := map[string]string{
							"productUrl": app.GetStorefrontURL("product/" + currentProductID),
							"imageUrl":   app.GetStorefrontURL(mediaPath + currentProduct.GetDefaultImage()),
							"name":       currentProduct.GetName(),
							"sku":        currentProduct.GetSku(),
							"brand":      productBrand,
						}

						requestData["products"] = append(utils.InterfaceToArray(requestData["products"]), productInfo)
					}

					fmt.Println(requestData)

					requestURL := "https://api.trustpilot.com/v1/private/product-reviews/business-units/{businessUnitId}/invitation-links"

					requestURL = strings.Replace(requestURL, "{businessUnitId}", trustPilotBusinessUnitID, 1)
					fmt.Println("send request to: ", requestURL)

					jsonString := utils.EncodeToJSONString(requestData)
					buffer := bytes.NewBuffer([]byte(jsonString))

					request, err := http.NewRequest("POST", requestURL, buffer)
					request.Header.Set("Content-Type", "application/json")
					request.Header.Set("Authorization", "Bearer"+currentToken.Access)

					if err != nil {
						env.LogError(err)
					}

					client := &http.Client{}
					response, err := client.Do(request)
					if err != nil {
						env.LogError(err)
					}

					defer response.Body.Close()

					responseBody, err := ioutil.ReadAll(response.Body)
					if err != nil {
						env.LogError(err)
					}

					jsonResponse, err := utils.DecodeJSONToStringKeyMap(responseBody)
					if err != nil {
						env.LogError(err)
					}

					fmt.Println(jsonResponse)
					if result, ok := jsonResponse["result"]; ok {
						fmt.Println(result)
					}
					fmt.Println("end function")
				}
			} else {
				env.LogError(env.ErrorNew(ConstErrorModule, env.ConstErrorLevelActor, "22207d49-e001-4666-8501-26bf5ef0926b", "some of trust pilot settings are blank"))
			}
		}

	}(eventData["order"].(order.InterfaceOrder), eventData["cart"].(cart.InterfaceCart))

	return true
}
