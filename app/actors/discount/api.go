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
// that are aggregated and using final values
func APIGetDiscounts(context api.InterfaceApplicationContext) (interface{}, error) {

	var result []StructAggregatedDiscount
	groupedDiscounts := make(map[string]StructAggregatedDiscount)

	currentCheckout, err := checkout.GetCurrentCheckout(context, false)
	if err != nil || currentCheckout == nil {
		return nil, env.ErrorDispatch(err)
	}

	currentCheckout.GetGrandTotal()

	usedDiscounts := currentCheckout.GetDiscounts()

	for _, currentDiscount := range usedDiscounts {
		key := currentDiscount.Code + currentDiscount.Type
		var groupedDiscount StructAggregatedDiscount

		if savedDiscount, present := groupedDiscounts[key]; present {
			groupedDiscount = savedDiscount

			groupedDiscount.Amount = groupedDiscount.Amount + currentDiscount.Amount

			if _, present := groupedDiscount.Object[currentDiscount.Object]; present {
				groupedDiscount.Object[currentDiscount.Object] += 1
			} else {
				groupedDiscount.Object[currentDiscount.Object] = 1
			}
		} else {
			groupedDiscount.Code = currentDiscount.Code
			groupedDiscount.Name = currentDiscount.Name
			groupedDiscount.Amount = currentDiscount.Amount
			groupedDiscount.Type = currentDiscount.Type

			groupedDiscount.Object = map[string]int{currentDiscount.Object: 1}
		}

		groupedDiscounts[key] = groupedDiscount
	}

	for _, discount := range groupedDiscounts {
		result = append(result, discount)
	}

	return result, nil
}
