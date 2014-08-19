package paypal

import (
	"github.com/ottemo/foundation/app/models/checkout"
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

func (it *PayPal) Void()	error {
	return nil
}
