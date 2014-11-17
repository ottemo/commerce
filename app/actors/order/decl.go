// Package order is a default implementation of interfaces declared in
// "github.com/ottemo/foundation/app/models/order" package
package order

import (
	"github.com/ottemo/foundation/app/models/order"
	"github.com/ottemo/foundation/db"

	"sync"
	"time"
)

// Package global variables
var (
	lastIncrementId      int = 0
	lastIncrementIdMutex sync.Mutex
)

// Package global constants
const (
	COLLECTION_NAME_ORDER       = "orders"
	COLLECTION_NAME_ORDER_ITEMS = "order_items"

	INCREMENT_ID_FORMAT = "%0.10d"

	CONFIG_PATH_LAST_INCREMENT_ID = "internal.order.increment_id"
)

// DefaultOrderItem is a default implementer of I_OrderItem
type DefaultOrderItem struct {
	id  string
	idx int

	OrderId string

	ProductId string

	Qty int

	Name string
	Sku  string

	ShortDescription string

	Options map[string]interface{}

	Price  float64
	Weight float64
}

// DefaultOrder is a default implementer of I_Order
type DefaultOrder struct {
	id string

	IncrementId string
	Status      string

	VisitorId string
	CartId    string

	Description string
	PaymentInfo map[string]interface{}

	BillingAddress  map[string]interface{}
	ShippingAddress map[string]interface{}

	CustomerEmail string
	CustomerName  string

	PaymentMethod  string
	ShippingMethod string

	Subtotal       float64
	Discount       float64
	TaxAmount      float64
	ShippingAmount float64
	GrandTotal     float64

	CreatedAt time.Time
	UpdatedAt time.Time

	Items map[int]order.I_OrderItem

	maxIdx int
}

// DefaultOrderItemCollection is a default implementer of I_OrderCollection
type DefaultOrderCollection struct {
	listCollection     db.I_DBCollection
	listExtraAtributes []string
}

// DefaultOrderItemCollection is a default implementer of I_OrderItemCollection
type DefaultOrderItemCollection struct {
	listCollection     db.I_DBCollection
	listExtraAtributes []string
}
