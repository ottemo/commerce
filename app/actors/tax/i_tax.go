package tax

import (
	"github.com/ottemo/foundation/db"
	"github.com/ottemo/foundation/app/utils"
	"github.com/ottemo/foundation/app/models/checkout"
)

func (it *DefaultTax) GetName() string {
	return "Tax"
}

func (it *DefaultTax) GetCode() string {
	return "tax"
}

func processRecords(name string, records []map[string]interface{}, result []checkout.T_TaxRate) []checkout.T_TaxRate {
	for _, record := range records {
		taxRate := checkout.T_TaxRate {
			Name: name,
			Code: utils.InterfaceToString(record["code"]),
			Amount: utils.InterfaceToFloat64(record["rate"]),
		}

		result = append(result, taxRate)
	}

	return result
}

func (it *DefaultTax) CalculateTax(currentCheckout checkout.I_Checkout) []checkout.T_TaxRate {
	result := make([]checkout.T_TaxRate, 0)

	if shippingAddress := currentCheckout.GetShippingAddress(); shippingAddress != nil {
		state := shippingAddress.GetState()
		zip := shippingAddress.GetZipCode()

		if dbEngine := db.GetDBEngine(); dbEngine != nil {
			if collection, err := dbEngine.GetCollection("Taxes"); err == nil {
				collection.AddFilter("state", "=", "*")
				collection.AddFilter("zip", "=", "*")

				if records, err := collection.Load(); err == nil {
					result = processRecords(it.GetName(), records, result)
				}

				collection.ClearFilters()
				collection.AddFilter("state", "=", state)
				collection.AddFilter("zip", "=", "*")

				if records, err := collection.Load(); err == nil {
					result = processRecords(it.GetName(), records, result)
				}

				collection.ClearFilters()
				collection.AddFilter("state", "=", state)
				collection.AddFilter("zip", "=", zip)

				if records, err := collection.Load(); err == nil {
					result = processRecords(it.GetName(), records, result)
				}
			}
		}
	}

	return result
}
