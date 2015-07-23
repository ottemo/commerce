package coupon

import (
	"time"

	"github.com/ottemo/foundation/app/models/checkout"
	"github.com/ottemo/foundation/db"
	"github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/utils"
	"strings"
)

// GetName returns name of current discount implementation
func (it *DefaultDiscount) GetName() string {
	return "CouponDiscount"
}

// GetCode returns code of current discount implementation
func (it *DefaultDiscount) GetCode() string {
	return "coupon_discount"
}

// CalculateDiscount calculates and returns a set of discounts applied to given checkout
func (it *DefaultDiscount) CalculateDiscount(checkoutInstance checkout.InterfaceCheckout) []checkout.StructDiscount {

	var result []checkout.StructDiscount

	// checking session for applied coupon codes
	if currentSession := checkoutInstance.GetSession(); currentSession != nil {

		appliedCodes := utils.InterfaceToStringArray(currentSession.Get(ConstSessionKeyAppliedDiscountCodes))
		usedCodes := utils.InterfaceToStringArray(currentSession.Get(ConstSessionKeyUsedDiscountCodes))

		if len(appliedCodes) > 0 {

			// loading information about applied discounts
			collection, err := db.GetCollection(ConstCollectionNameCouponDiscounts)
			if err != nil {
				return result
			}
			err = collection.AddFilter("code", "in", appliedCodes)
			if err != nil {
				return result
			}

			records, err := collection.Load()
			if err != nil || len(records) == 0 {
				return result
			}

			// making coupon code map for right apply order ignoring used coupons and limited
			discountCodes := make(map[string]map[string]interface{})
			for _, record := range records {
				discountsUsageQty := discountsUsage(checkoutInstance, record)
				discountCode := utils.InterfaceToString(record["code"])

				if discountCode != "" && !utils.IsInArray(discountCode, usedCodes) && discountsUsageQty > 0 {
					record["usage_qty"] = discountsUsageQty
					discountCodes[discountCode] = record
				}
			}

			productsInCart := make(map[string]float64)

			// collect products by ID and price with priority to lower (for calculating % discount)
			for _, productInCart := range checkoutInstance.GetCart().GetItems() {
				if cartProduct := productInCart.GetProduct(); cartProduct != nil {
					cartProduct.ApplyOptions(productInCart.GetOptions())
					productPrice := cartProduct.GetPrice()
					productID := productInCart.GetProductID()

					if storedPrice, present := productsInCart[productID]; !present || productPrice < storedPrice {
						productsInCart[productID] = productPrice
					}
				}
			}

			productsDiscountPriorityValue := utils.InterfaceToFloat64(env.ConfigGetValue(ConstConfigPathDiscountApplyPriority))
			cartDiscountPriorityValue := productsDiscountPriorityValue + 0.01

			// accumulation of coupon discounts to result
			for appliedCodesIdx, discountCode := range appliedCodes {
				if discountCoupon, ok := discountCodes[discountCode]; ok {

					applyTimes := utils.InterfaceToInt(discountCoupon["times"])
					workSince := utils.InterfaceToTime(discountCoupon["since"])
					workUntil := utils.InterfaceToTime(discountCoupon["until"])

					currentTime := time.Now()

					// to be applicable coupon should satisfy following conditions:
					//   [applyTimes] should be -1 or >0 and [workSince] >= currentTime <= [workUntil] if set
					if (applyTimes == -1 || applyTimes > 0) &&
						(utils.IsZeroTime(workSince) || workSince.Unix() <= currentTime.Unix()) &&
						(utils.IsZeroTime(workUntil) || workUntil.Unix() >= currentTime.Unix()) {

						// calculating coupon discount amount
						discountAmount := utils.InterfaceToFloat64(discountCoupon["amount"])
						discountPercent := utils.InterfaceToFloat64(discountCoupon["percent"])
						discountTarget := utils.InterfaceToString(discountCoupon["target"])
						discountUsageQty := utils.InterfaceToFloat64(discountCoupon["usage_qty"])

						discountPercent = discountPercent * discountUsageQty
						discountAmount = discountAmount * discountUsageQty

						// case it's a cart discount we just add them to result
						if strings.Contains(discountTarget, checkout.ConstDiscountObjectCart) || discountTarget == "" {

							if discountPercent > 0 {
								result = append(result, checkout.StructDiscount{
									Name:      utils.InterfaceToString(discountCoupon["name"]),
									Code:      utils.InterfaceToString(discountCoupon["code"]),
									Amount:    discountPercent,
									IsPercent: true,
									Priority:  cartDiscountPriorityValue,
									Object:    checkout.ConstDiscountObjectCart,
								})
								cartDiscountPriorityValue += float64(0.0001)
							}

							if discountAmount > 0 {
								result = append(result, checkout.StructDiscount{
									Name:      utils.InterfaceToString(discountCoupon["name"]),
									Code:      utils.InterfaceToString(discountCoupon["code"]),
									Amount:    discountAmount,
									IsPercent: false,
									Priority:  cartDiscountPriorityValue,
									Object:    checkout.ConstDiscountObjectCart,
								})
								cartDiscountPriorityValue += float64(0.0001)
							}

							continue
						}

						// parse target as array of productIDs on which we will apply discount
						for _, productID := range utils.InterfaceToStringArray(discountTarget) {

							if cartProductPrice, present := productsInCart[productID]; present {
								totalProductDiscountAmount := discountPercent/100*cartProductPrice + discountAmount

								if discountPercent > 0 {
									result = append(result, checkout.StructDiscount{
										Name:      utils.InterfaceToString(discountCoupon["name"]),
										Code:      utils.InterfaceToString(discountCoupon["code"]),
										Amount:    totalProductDiscountAmount,
										IsPercent: false,
										Priority:  productsDiscountPriorityValue,
										Object:    productID,
									})
									productsDiscountPriorityValue += float64(0.0001)
								}

							}
						}

					} else {
						// we have not applicable coupon - removing it from applied coupons list
						newAppliedCodes := make([]string, 0, len(appliedCodes)-1)
						for idx, value := range appliedCodes {
							if idx != appliedCodesIdx {
								newAppliedCodes = append(newAppliedCodes, value)
							}
						}
						currentSession.Set(ConstSessionKeyAppliedDiscountCodes, newAppliedCodes)
					}
				}
			}
		}
	}

	return result
}

// checks discount limiting parameters for correspondence to current checkout values
// return qty of usages if discount is allowed for current checkout and satisfy all conditions
func discountsUsage(checkoutInstance checkout.InterfaceCheckout, couponDiscount map[string]interface{}) int {

	result := -1
	if limits, present := couponDiscount["limits"]; present {
		limitations := utils.InterfaceToMap(limits)
		if len(limitations) > 0 {

			productsInCart := make(map[string]int)
			var productID string
			var productQty int

			// collect products to one map by ID and qty
			for _, productInCart := range checkoutInstance.GetCart().GetItems() {
				productID = productInCart.GetProductID()
				productQty = productInCart.GetQty()

				if qty, present := productsInCart[productID]; present {
					productsInCart[productID] = qty + productQty
					continue
				}
				productsInCart[productID] = productQty
			}

			for limitingKey, limitingValue := range limitations {

				switch strings.ToLower(limitingKey) {
				case "product_in_cart":
					requiredProduct := utils.InterfaceToStringArray(limitingValue)
					for index, productID := range requiredProduct {
						if _, present := productsInCart[productID]; present {
							break
						}
						if index == (len(requiredProduct) - 1) {
							return 0
						}
					}

				case "products_in_cart":
					requiredProducts := utils.InterfaceToStringArray(limitingValue)
					for _, productID := range requiredProducts {
						if _, present := productsInCart[productID]; !present {
							return 0
						}
					}

				case "products_in_qty":
					requiredProducts := utils.InterfaceToMap(limitingValue)
					for requiredProductID, requiredQty := range requiredProducts {
						productQty, present := productsInCart[requiredProductID]
						limitingQty := utils.InterfaceToInt(productQty / utils.InterfaceToInt(requiredQty))

						if !present || limitingQty < 1 {
							return 0
						}

						if result == -1 || limitingQty < result {
							result = limitingQty
						}
					}
				case "max_usage_qty":
					if limitingQty := utils.InterfaceToInt(limitingValue); limitingQty >= 1 && (result == -1 || limitingQty < result) {
						result = limitingQty
					}
				}
			}
		}
	}
	if result == -1 {
		result = 1
	}

	return result
}
