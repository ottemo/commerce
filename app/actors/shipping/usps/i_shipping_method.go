package usps

import (
	"bytes"

	"html"
	"regexp"
	"text/template"

	"net/http"
	"net/url"

	"io/ioutil"

	"github.com/ottemo/foundation/app/models/checkout"

	"github.com/ottemo/foundation/app/utils"
	"github.com/ottemo/foundation/env"

	"launchpad.net/xmlpath"
	"strings"
)

func (it *USPS) GetName() string {
	return SHIPPING_NAME
}

func (it *USPS) GetCode() string {
	return SHIPPING_CODE
}

func (it *USPS) IsAllowed(checkout checkout.I_Checkout) bool {
	return utils.InterfaceToBool(env.ConfigGetValue(CONFIG_PATH_ENABLED))
}

func (it *USPS) GetRates(checkoutObject checkout.I_Checkout) []checkout.T_ShippingRate {

	result := make([]checkout.T_ShippingRate, 0)

	templateValues := make(map[string]interface{})

	templateValues["userid"] = utils.InterfaceToString(env.ConfigGetValue(CONFIG_PATH_USER))       // "133OTTEM1795",
	templateValues["origin"] = utils.InterfaceToString(env.ConfigGetValue(CONFIG_PATH_ORIGIN_ZIP)) // "44106",

	if templateValues["userid"] == "" || templateValues["origin"] == "" {
		return nil
	}

	if shippingAddress := checkoutObject.GetShippingAddress(); shippingAddress != nil || shippingAddress.GetZipCode() != "" {
		templateValues["destination"] = shippingAddress.GetZipCode()
	} else {
		return result
	}

	var pounds float64 = 0
	var ounces float64 = 0
	if checkoutCart := checkoutObject.GetCart(); checkoutCart != nil {

		cartItems := checkoutCart.GetItems()
		if len(cartItems) == 0 {
			return result
		}

		defaultWeight := utils.InterfaceToFloat64(env.ConfigGetValue(CONFIG_PATH_DEFAULT_WEIGHT))

		for _, cartItem := range cartItems {
			cartProduct := cartItem.GetProduct()
			if cartProduct == nil {
				continue
			}

			if cartProduct.GetWeight() == 0 {
				pounds += defaultWeight * utils.InterfaceToFloat64(cartItem.GetQty())
			} else {
				pounds += cartProduct.GetWeight() * utils.InterfaceToFloat64(cartItem.GetQty())
			}
		}
	}

	templateValues["pounds"] = pounds
	templateValues["ounces"] = ounces

	templateValues["container"] = utils.InterfaceToString(env.ConfigGetValue(CONFIG_PATH_CONTAINER))
	templateValues["size"] = utils.InterfaceToString(env.ConfigGetValue(CONFIG_PATH_SIZE))

	templateValues["width"] = 0.1
	templateValues["height"] = 0.1
	templateValues["length"] = 0.1
	templateValues["girth"] = 0.1

	dimensions := strings.Split(utils.InterfaceToString(env.ConfigGetValue(CONFIG_PATH_DEFAULT_DIMENSIONS)), "x")
	for idx, dimensionValue := range dimensions {
		dimensionValue = strings.Trim(dimensionValue, " ")
		switch idx {
		case 0:
			templateValues["width"] = utils.InterfaceToFloat64(dimensionValue)
		case 1:
			templateValues["height"] = utils.InterfaceToFloat64(dimensionValue)
		case 2:
			templateValues["length"] = utils.InterfaceToFloat64(dimensionValue)
		case 3:
			templateValues["girth"] = utils.InterfaceToFloat64(dimensionValue)
		}

	}

	requestTemplate := `<RateV4Request USERID="{{.userid}}">
     <Revision/>
     <Package ID="1">
          <Service>ALL</Service>
          <ZipOrigination>{{.origin}}</ZipOrigination>
          <ZipDestination>{{.destination}}</ZipDestination>
          <Pounds>{{.pounds}}</Pounds>
          <Ounces>{{.ounces}}</Ounces>
          <Container>{{.container}}</Container>
          <Size>{{.size}}</Size>
          {{if eq .size "LARGE" }}
          <Width>{{.width}}</Width>
          <Length>{{.length}}</Length>
          <Height>{{.height}}</Height>
          <Girth>{{.girth}}</Girth>
          {{end}}
          <Machinable>True</Machinable>
    </Package>
    </RateV4Request>`

	var buff bytes.Buffer
	parsedTemplate, _ := template.New("usps").Parse(requestTemplate)
	parsedTemplate.Execute(&buff, templateValues)

	// println(buff.String())

	query := HTTP_ENDPOINT + "?API=RateV4&XML=" + url.QueryEscape(buff.String())
	response, err := http.Get(query)
	if err != nil {
		return result
	}

	var responseData []byte
	if response.ContentLength > 0 {
		responseData = make([]byte, response.ContentLength)
		_, err := response.Body.Read(responseData)
		if err != nil {
			return result
		}
	} else {
		responseData, err = ioutil.ReadAll(response.Body)
		if err != nil {
			return result
		}
	}

	// println(string(responseData))

	root, err := xmlpath.Parse(bytes.NewReader(responseData))
	if err != nil {
		return result
	}

	allowedMethodsArray := utils.InterfaceToArray(env.ConfigGetValue(CONFIG_PATH_ALLOWED_METHODS))

	postage, _ := xmlpath.Compile("//Postage")
	service, _ := xmlpath.Compile("./MailService")
	rate, _ := xmlpath.Compile("./Rate")
	code, _ := xmlpath.Compile("./@CLASSID")

	regexTags := regexp.MustCompile("<[^>]+>")

	for i := postage.Iter(root); i.Next(); {
		postageNode := i.Node()

		stringService, ok := service.String(postageNode)
		if !ok {
			continue
		}

		stringRate, ok := rate.String(postageNode)
		if !ok {
			continue
		}

		stringCode, ok := code.String(postageNode)
		if !ok {
			continue
		}

		if len(allowedMethodsArray) == 0 || utils.IsInArray(stringCode, allowedMethodsArray) {

			rateName := html.UnescapeString(stringService)
			if REMOVE_RATE_NAME_TAGS {
				rateName = regexTags.ReplaceAllString(rateName, "")

			}

			result = append(result,
				checkout.T_ShippingRate{
					Code:  stringCode,
					Name:  rateName,
					Price: utils.InterfaceToFloat64(stringRate),
				})
		}
	}

	return result
}
