// Package paypal is a PayPal implementation of payment method interface declared in
// "github.com/ottemo/foundation/app/models/checkout" package
package paypal

import (
	"sync"
)

// Package global constants
const (
	PAYMENT_CODE = "paypal_express"
	PAYMENT_NAME = "PayPal Express"

	PAYMENT_ACTION_SALE          = "Sale"
	PAYMENT_ACTION_AUTHORIZATION = "Authorization"

	CONFIG_PATH_GROUP = "payment.paypal"

	CONFIG_PATH_ENABLED = "payment.paypal.enabled"
	CONFIG_PATH_TITLE   = "payment.paypal.title"

	CONFIG_PATH_NVP     = "payment.paypal.nvp"
	CONFIG_PATH_GATEWAY = "payment.paypal.gateway"

	CONFIG_PATH_USER = "payment.paypal.user"
	CONFIG_PATH_PASS = "payment.paypal.password"

	CONFIG_PATH_SIGNATURE = "payment.paypal.signature"
	CONFIG_PATH_ACTION    = "payment.paypal.action"
)

// Package global variables
var (
	waitingTokens      = make(map[string]interface{})
	waitingTokensMutex sync.RWMutex
)

// PayPalExpress is a implementer of I_PaymentMethod for a PayPal Express method
type PayPalExpress struct{}

// PayPalExpress is a implementer of I_PaymentMethod for a PayPal REST method (currently not working)
type PayPalRest struct{}
