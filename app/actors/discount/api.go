package discount

import (
	"github.com/ottemo/foundation/api"
	"github.com/ottemo/foundation/app/models/checkout"
	"github.com/ottemo/foundation/env"
)

// setupAPI setups package related API endpoint routines
func setupAPI() error {
	var err error

	err = api.GetRestService().RegisterAPI("discounts", api.ConstRESTOperationGet, APIGetDiscounts)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	return nil
}

// APIGetDiscounts returns list of applied discounts
func APIGetDiscounts(context api.InterfaceApplicationContext) (interface{}, error) {

	var result []checkout.StructDiscount
	groupedDiscounts := make(map[string]checkout.StructDiscount)

	currentCheckout, err := checkout.GetCurrentCheckout(context)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	currentCheckout.CalculateAmount(0)

	usedDiscounts := currentCheckout.GetDiscounts()

	for _, currentDiscount := range usedDiscounts {
		key := currentDiscount.Code + currentDiscount.Type

		if savedDiscount, present := groupedDiscounts[key]; present {
			currentDiscount.Amount = savedDiscount.Amount + currentDiscount.Amount
			currentDiscount.Object = savedDiscount.Object + ", " + currentDiscount.Object
			currentDiscount.IsPercent = false
		}

		groupedDiscounts[key] = currentDiscount
	}

	for _, discount := range groupedDiscounts {
		result = append(result, discount)
	}

	return result, nil
}
