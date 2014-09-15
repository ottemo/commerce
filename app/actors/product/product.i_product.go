package product

import (
	"fmt"
	"sort"
	"strings"
	"errors"
	"github.com/ottemo/foundation/app/utils"
)

func (it *DefaultProduct) GetSku() string  { return it.Sku }
func (it *DefaultProduct) GetName() string { return it.Name }

func (it *DefaultProduct) GetShortDescription() string { return it.ShortDescription }
func (it *DefaultProduct) GetDescription() string      { return it.Description }

func (it *DefaultProduct) GetDefaultImage() string { return it.DefaultImage }

func (it *DefaultProduct) GetPrice() float64 { return it.Price }
func (it *DefaultProduct) GetWeight() float64 { return it.Weight }

func (it *DefaultProduct) GetOptions() map[string]interface{}   { return it.Options }

// internal usage function to update order item fields according to options
func (it *DefaultProduct) ApplyOptions(options map[string]interface{}) error {

	// taking item specified options and product options
	productOptions := it.GetOptions()

	// storing start price for a case of percentage price modifier
	startPrice := it.GetPrice()

	// sorting applicable product attributes according to "order" field
	optionsApplyOrder := make([]string, 0)
	for itemOptionName, _ := range options {

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
						itemOptionValueSet := []string{}
						switch typedOptionValue := itemOptionValue.(type) {
						case string:
							itemOptionValueSet = []string{typedOptionValue}
						case []string:
							itemOptionValueSet = typedOptionValue
						default:
							return errors.New("invalid value for " + itemOptionName)
						}

						// loop through option values customer set for product
						for _, itemOptionValue := range itemOptionValueSet {

							if productOptionValue, present := productOptionValues[itemOptionValue]; present {
								if productOptionValue, ok := productOptionValue.(map[string]interface{}); ok {
									applyOptionModifiers(productOptionValue)
								}
							} else {
								return errors.New("invalid '" + itemOptionName + "' option value: '" + itemOptionValue)
							}

						}

						// cleaning option values were not used by customer
						for productOptionValueName, _ := range productOptionValues {
							if !utils.IsInArray(productOptionValueName, itemOptionValueSet) {
								delete(productOptionValues, productOptionValueName)
							}
						}
					}
				}
			}
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

	return nil
}
