package token

import (
	"github.com/ottemo/foundation/utils"
	"time"
	"github.com/ottemo/foundation/app/models/checkout"
	"github.com/ottemo/foundation/env"
)

// GetVisitorID returns the Visitor ID for the Visitor Card
func (it *DefaultVisitorCard) GetVisitorID() string { return it.visitorID }

// GetHolderName returns the Holder of the Credit Card
func (it *DefaultVisitorCard) GetHolderName() string { return it.Holder }

// GetPaymentMethod returns the Payment method code of the Visitor Card
func (it *DefaultVisitorCard) GetPaymentMethod() string { return it.Payment }

// GetType will return the Type of the Visitor Card
func (it *DefaultVisitorCard) GetType() string { return it.Type }

// GetNumber will return the Number attribute of the Visitor Card
func (it *DefaultVisitorCard) GetNumber() string { return it.Number }

// GetExpirationDate will return the Expiration date  of the Visitor Card
func (it *DefaultVisitorCard) GetExpirationDate() string {

	if it.ExpirationDate == "" {
		it.ExpirationDate = utils.InterfaceToString(it.ExpirationMonth) + "/" + utils.InterfaceToString(it.ExpirationYear)
	}

	return it.ExpirationDate
}

// GetToken will return the Token of the Visitor Card
func (it *DefaultVisitorCard) GetToken() string { return it.Token }

// IsExpired will return Expired status of the Visitor Card
func (it *DefaultVisitorCard) IsExpired() bool {
	current := time.Now()
	return it.ExpirationYear < utils.InterfaceToInt(current.Year()) || it.ExpirationMonth < utils.InterfaceToInt(current.Month())
}

// UPDATE it please :D
func (it *DefaultVisitorCard) CreateVisitorCard(ccInfo map[string]interface {}) error {

	paymentMethodCode := utils.InterfaceToString(utils.GetFirstMapValue(ccInfo, "payment", "payment_method"))
	if paymentMethodCode == "" {
		return env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "6d1691c8-2d26-44be-b90d-24d920e26301", "payment method not selected")
	}

	value, present := ccInfo["cc"]
	if !present {
		return env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "2e9f1bfc-ec9f-4017-83c6-4d04b95b9c08", "payment info not specified")
	}

	creditCardInfo := utils.InterfaceToMap(value)

	var paymentMethod checkout.InterfacePaymentMethod

	for _, payment := range checkout.GetRegisteredPaymentMethods() {
		if payment.GetCode() == paymentMethodCode {
			if !payment.IsTokenable(nil) {
				return env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "519ef43c-4d07-4b64-90f7-7fdc3657940a", "for selected payment method credit card can't be saved")
			}
			paymentMethod = payment
		}
	}

	if paymentMethod == nil {
		return env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "c80c4106-1208-4d0b-8577-0889f608869b", "such payment method not existing")
	}

	paymentInfo := map[string]interface{}{
		checkout.ConstPaymentActionTypeKey: checkout.ConstPaymentActionTypeCreateToken,
		"cc": creditCardInfo,
	}

	// contains creditCardLastFour, creditCardType, responseMessage, responseResult, transactionID, creditCardExp
	paymentResult, err := paymentMethod.Authorize(nil, paymentInfo)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	cardInfoMap := utils.InterfaceToMap(paymentResult)
	if !utils.KeysInMapAndNotBlank(cardInfoMap, "transactionID", "creditCardLastFour") {
		return env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "22e17290-56f3-452a-8d54-18d5a9eb2833", "transaction can't be obtained")
	}

	// create credit card map with info
	tokenRecord := map[string]interface{}{
		"visitor_id":      visitorID,
		"payment":         paymentMethodCode,
		"type":            cardInfoMap["creditCardType"],
		"number":          cardInfoMap["creditCardLastFour"],
		"expiration_date": cardInfoMap["creditCardExp"],
		"holder":          utils.InterfaceToString(ccInfo["holder"]),
		"token":           cardInfoMap["transactionID"],
		"updated":         time.Now(),
	}

	err = it.FromHashMap(tokenRecord)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	err = it.Save()
	if err != nil {
		return env.ErrorDispatch(err)
	}
	
	
	return nil
}
