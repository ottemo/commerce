package order

import (
	"time"
	"github.com/ottemo/foundation/app/models/order"
)

const (
	ORDER_COLLECTION_NAME = "orders"
	ORDER_ITEMS_COLLECTION_NAME = "order_items"
)

type DefaultOrderItem struct {
	id string
	idx int

	order *DefaultOrder

	ProductId string

	Qty int

	Name string
	Sku string

	ShortDescription string

	ProductOptions map[string]interface{}

	Price float64
	Weight float64
	Size float64
}



type DefaultOrder struct {
	id string

	IncrementId string
	Status string

	VisitorId string
	CartId string

	CustomerEmail string
	CustomerName string

	PaymentMethod string
	ShippingMethod string

	Subtotal float64
	Discount float64
	TaxAmount float64
	ShippingAmount float64
	GrandTotal float64

	CreatedAt time.Time
	UpdaedAt time.Time

	Items map[int]order.I_OrderItem

	maxIdx int
}
