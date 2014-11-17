// Package discount is a default implementation of discount interface declared in
// "github.com/ottemo/foundation/app/models/checkout" package
package discount

// Package global constants
const (
	SESSION_KEY_APPLIED_DISCOUNT_CODES = "applied_discount_codes"
	COLLECTION_NAME_COUPON_DISCOUNTS   = "coupon_discounts"
)

// DefaultDiscount is a default implementer of I_Discount
type DefaultDiscount struct{}
