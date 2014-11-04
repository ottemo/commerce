package order

import (
	"github.com/ottemo/foundation/app/models/order"
	"github.com/ottemo/foundation/db"

	"sync"
	"time"
)

var (
	lastIncrementId      int = 0
	lastIncrementIdMutex sync.Mutex
)

const (
	COLLECTION_NAME_ORDER       = "orders"
	COLLECTION_NAME_ORDER_ITEMS = "order_items"

	INCREMENT_ID_FORMAT = "%0.10d"

	CONFIG_PATH_LAST_INCREMENT_ID = "internal.order.increment_id"
)

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

type DefaultOrderCollection struct {
	listCollection     db.I_DBCollection
	listExtraAtributes []string
}

type DefaultOrderItemCollection struct {
	listCollection     db.I_DBCollection
	listExtraAtributes []string
}
