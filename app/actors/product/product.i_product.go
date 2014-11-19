package product

import (
	"fmt"
	"sort"
	"strings"

	"github.com/ottemo/foundation/app/models/product"
	"github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/utils"
)

// GetSku will return requested sku for the given product
func (it *DefaultProduct) GetSku() string { return it.Sku }

// GetName will return the name of the given product
func (it *DefaultProduct) GetName() string { return it.Name }

// GetShortDescription will return the short description of the requested product
func (it *DefaultProduct) GetShortDescription() string { return it.ShortDescription }

// GetDescription will return the long description of the requested product
func (it *DefaultProduct) GetDescription() string { return it.Description }

// GetDefaultImage will return the imaged identified as defult for the given product
func (it *DefaultProduct) GetDefaultImage() string { return it.DefaultImage }

// GetPrice will return the price as a float64 for the given product
func (it *DefaultProduct) GetPrice() float64 { return it.Price }

// GetWeight will return the weight for the given product
func (it *DefaultProduct) GetWeight() float64 { return it.Weight }

// GetOptions will return the products otions as a map[string]interface{}
func (it *DefaultProduct) GetOptions() map[string]interface{} { return it.Options }

// ApplyOptions is an internal usage function to update order item fields according to options
func (it *DefaultProduct) ApplyOptions(options map[string]interface{}) error {
	// taking item specified options and product options
	productOptions := it.GetOptions()

	// storing start price for a case of percentage price modifier
	startPrice := it.GetPrice()

	// sorting applicable product attributes according to "order" field
	// optionsApplyOrder := make([]string, 0)
	var optionsApplyOrder []string

	for itemOptionName := range options {

		// looking only for options that customer customer set for item
		if productOption, present := productOptions[itemOptionName]; present {
			if productOption, ok := productOption.(map[string]interface{}); ok {

				orderValue := int(^uint(0) >> 1) // default order - max integer
				if optionValue, present := productOption["order"]; present {
					orderValue = utils.InterfaceToInt(optionValue)
				}

				// encoding key order to string "000000000000001 [attribute name]"
				// for future sort as string (16 digits - max for js integer)
				key := fmt.Sprintf("%.16d %s", orderValue, itemOptionName)
				optionsApplyOrder = append(optionsApplyOrder, key)
			}
		}
	}
	sort.Strings(optionsApplyOrder)

	// function to modify orderItem according to option values
	applyOptionModifiers := func(optionToApply map[string]interface{}) {
		// price modifier
		if optionValue, present := optionToApply["price"]; present {
			addPrice := utils.InterfaceToFloat64(optionValue)
			priceType := utils.InterfaceToString(optionToApply["price_type"])

			if priceType == "percent" || priceType == "%" {
				it.Price += startPrice * addPrice / 100
			} else {
				it.Price += addPrice
			}
		}

		// sku modifier
		if optionValue, present := optionToApply["sku"]; present {
			skuModifier := utils.InterfaceToString(optionValue)
			it.Sku += "-" + skuModifier
		}
	}

	// loop over item applied option in right order
	for _, itemOptionName := range optionsApplyOrder {
		// decoding key order after sort
		itemOptionName := itemOptionName[strings.Index(itemOptionName, " ")+1:]
		itemOptionValue := options[itemOptionName]

		if productOption, present := productOptions[itemOptionName]; present {
			if productOptions, ok := productOption.(map[string]interface{}); ok {

				// product option itself can contain price, sku modifiers
				applyOptionModifiers(productOptions)

				// if product option value have predefined option values, then checking their modifiers
				if productOptionValues, present := productOptions["options"]; present {
					if productOptionValues, ok := productOptionValues.(map[string]interface{}); ok {

						// option user set can be single on multi-value
						// making it uniform
						// itemOptionValueSet := make([]string, 0)
						var itemOptionValueSet []string
						switch typedOptionValue := itemOptionValue.(type) {
						case string:
							itemOptionValueSet = append(itemOptionValueSet, typedOptionValue)
						case []string:
							itemOptionValueSet = typedOptionValue
						case []interface{}:
							for _, value := range typedOptionValue {
								if value, ok := value.(string); ok {
									itemOptionValueSet = append(itemOptionValueSet, value)
								}
							}
						default:
							return env.ErrorNew("unexpected option value for " + itemOptionName + " option")
						}

						// loop through option values customer set for product
						for _, itemOptionValue := range itemOptionValueSet {

							if productOptionValue, present := productOptionValues[itemOptionValue]; present {
								if productOptionValue, ok := productOptionValue.(map[string]interface{}); ok {
									applyOptionModifiers(productOptionValue)
								}
							} else {
								return env.ErrorNew("invalid '" + itemOptionName + "' option value: '" + itemOptionValue)
							}

						}

						// cleaning option values were not used by customer
						for productOptionValueName := range productOptionValues {
							if !utils.IsInArray(productOptionValueName, itemOptionValueSet) {
								delete(productOptionValues, productOptionValueName)
							}
						}
					}
				}
			}
		} else {
			return env.ErrorNew("unknown option '" + itemOptionName + "'")
		}
	}

	// cleaning options were not used by customer
	for productOptionName, productOption := range productOptions {
		if _, present := options[productOptionName]; present {
			if productOption, ok := productOption.(map[string]interface{}); ok {
				productOption["value"] = options[productOptionName]
			}
		} else {
			delete(productOptions, productOptionName)
		}
	}

	it.Price = utils.RoundPrice(it.Price)

	return nil
}

// GetRelatedProductIds will return the related product IDs
func (it *DefaultProduct) GetRelatedProductIds() []string {
	// result := make([]string, 0)
	var result []string

	for _, productID := range it.RelatedProductIds {
		productModel, err := product.LoadProductById(productID)
		if err == nil {
			result = append(result, productModel.GetId())
		}
	}

	return result
}

// GetRelatedProducts will return an array of related products
func (it *DefaultProduct) GetRelatedProducts() []product.InterfaceProduct {
	// result := make([]product.InterfaceProduct, 0)
	var result []product.InterfaceProduct

	for _, productID := range it.RelatedProductIds {
		if productID == "" {
			continue
		}

		productModel, err := product.LoadProductById(productID)
		if err == nil {
			result = append(result, productModel)
		}
	}

	return result
}
