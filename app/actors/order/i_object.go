package order

import (
	"errors"
	"strings"

	"github.com/ottemo/foundation/app/models"
	"github.com/ottemo/foundation/app/utils"

	"github.com/ottemo/foundation/app/models/order"
)

//------------------
// DefaultOrderItem
//------------------

// returns attribute of OrderItem or nil
func (it *DefaultOrderItem) Get(attribute string) interface{} {

	switch strings.ToLower(attribute) {
	case "_id", "id":
		return it.id

	case "idx":
		return it.idx

	case "order":
		orderInstance, err := order.LoadOrderById(it.OrderId)
		if err == nil {
			return orderInstance
		}
		return nil

	case "order_id":
		return it.OrderId

	case "product_id":
		return it.ProductId

	case "qty":
		return it.Qty

	case "name":
		return it.Name

	case "sku":
		return it.Sku

	case "short_description":
		return it.ShortDescription

	case "options":
		return it.ProductOptions

	case "price":
		return it.Price

	case "weight":
		return it.Weight

	case "size":
		return it.Size
	}

	return nil
}

// sets attribute to OrderItem object, returns error on problems
func (it *DefaultOrderItem) Set(attribute string, value interface{}) error {
	attribute = strings.ToLower(attribute)

	switch attribute {
	case "_id", "id":
		it.id = utils.InterfaceToString(value)

	case "idx":
		it.idx = utils.InterfaceToInt(value)

	case "order_id":
		it.OrderId = utils.InterfaceToString(value)

	case "product_id":
		it.ProductId = utils.InterfaceToString(value)

	case "qty":
		it.Qty = utils.InterfaceToInt(value)

	case "name":
		it.Name = utils.InterfaceToString(value)

	case "sku":
		it.Sku = utils.InterfaceToString(value)

	case "short_description":
		it.ShortDescription = utils.InterfaceToString(value)

	case "options":
		if mapValue, ok := value.(map[string]interface{}); ok {
			it.ProductOptions = mapValue
		} else {
			return errors.New("options should me map[string]interface{} type")
		}

	case "price":
		it.Price = utils.InterfaceToFloat64(value)

	case "weight":
		it.Weight = utils.InterfaceToFloat64(value)

	case "size":
		it.Size = utils.InterfaceToFloat64(value)

	default:
		return errors.New("unknown attribute: " + attribute)
	}

	return nil
}

// fills OrderItem attributes with values provided in input map
func (it *DefaultOrderItem) FromHashMap(input map[string]interface{}) error {

	for attribute, value := range input {
		if err := it.Set(attribute, value); err != nil {
			return err
		}
	}

	return nil
}

// makes map from OrderItem attribute values
func (it *DefaultOrderItem) ToHashMap() map[string]interface{} {

	result := make(map[string]interface{})

	result["_id"] = it.Get("_id")
	result["idx"] = it.Get("idx")

	result["order_id"] = it.Get("order_id")
	result["product_id"] = it.Get("product_id")

	result["qty"] = it.Get("qty")

	result["name"] = it.Get("name")
	result["sku"] = it.Get("sku")
	result["short_description"] = it.Get("short_description")

	result["options"] = it.Get("options")

	result["price"] = it.Get("price")
	result["weight"] = it.Get("weight")
	result["size"] = it.Get("size")

	return result
}

// describes attributes of OrderItem model
func (it *DefaultOrderItem) GetAttributesInfo() []models.T_AttributeInfo {

	info := []models.T_AttributeInfo{
		models.T_AttributeInfo{
			Model:      "OrderItem",
			Collection: "OrderItem",
			Attribute:  "_id",
			Type:       "id",
			IsRequired: false,
			IsStatic:   true,
			Label:      "ID",
			Group:      "General",
			Editors:    "not_editable",
			Options:    "",
			Default:    "",
		},
		models.T_AttributeInfo{
			Model:      "OrderItem",
			Collection: "OrderItem",
			Attribute:  "idx",
			Type:       "int",
			IsRequired: true,
			IsStatic:   true,
			Label:      "Increment ID",
			Group:      "General",
			Editors:    "integer",
			Options:    "",
			Default:    "",
		},
		models.T_AttributeInfo{
			Model:      "OrderItem",
			Collection: "OrderItem",
			Attribute:  "order_id",
			Type:       "id",
			IsRequired: true,
			IsStatic:   true,
			Label:      "Order",
			Group:      "General",
			Editors:    "not_editable",
			Options:    "",
			Default:    "",
		},
		models.T_AttributeInfo{
			Model:      "OrderItem",
			Collection: "OrderItem",
			Attribute:  "product_id",
			Type:       "id",
			IsRequired: false,
			IsStatic:   true,
			Label:      "Visitor",
			Group:      "General",
			Editors:    "not_editable",
			Options:    "",
			Default:    "",
		},
		models.T_AttributeInfo{
			Model:      "OrderItem",
			Collection: "OrderItem",
			Attribute:  "qty",
			Type:       "int",
			IsRequired: true,
			IsStatic:   true,
			Label:      "Qty",
			Group:      "General",
			Editors:    "integer",
			Options:    "",
			Default:    "",
		},
		models.T_AttributeInfo{
			Model:      "OrderItem",
			Collection: "OrderItem",
			Attribute:  "name",
			Type:       "varchar(150)",
			IsRequired: true,
			IsStatic:   true,
			Label:      "Name",
			Group:      "General",
			Editors:    "not_editable",
			Options:    "",
			Default:    "",
		},
		models.T_AttributeInfo{
			Model:      "OrderItem",
			Collection: "OrderItem",
			Attribute:  "sku",
			Type:       "varchar(100)",
			IsRequired: true,
			IsStatic:   true,
			Label:      "Sku",
			Group:      "General",
			Editors:    "not_editable",
			Options:    "",
			Default:    "",
		},
		models.T_AttributeInfo{
			Model:      "OrderItem",
			Collection: "OrderItem",
			Attribute:  "short_description",
			Type:       "varchar(255)",
			IsRequired: false,
			IsStatic:   true,
			Label:      "Short Description",
			Group:      "General",
			Editors:    "text",
			Options:    "",
			Default:    "",
		},
		models.T_AttributeInfo{
			Model:      "OrderItem",
			Collection: "OrderItem",
			Attribute:  "options",
			Type:       "text",
			IsRequired: false,
			IsStatic:   true,
			Label:      "Options",
			Group:      "General",
			Editors:    "product_options",
			Options:    "",
			Default:    "",
		},
		models.T_AttributeInfo{
			Model:      "OrderItem",
			Collection: "OrderItem",
			Attribute:  "price",
			Type:       "decimal(10,2)",
			IsRequired: true,
			IsStatic:   true,
			Label:      "Price",
			Group:      "Prices",
			Editors:    "numeric",
			Options:    "",
			Default:    "",
		},
		models.T_AttributeInfo{
			Model:      "OrderItem",
			Collection: "OrderItem",
			Attribute:  "weight",
			Type:       "decimal(10,2)",
			IsRequired: false,
			IsStatic:   true,
			Label:      "Weight",
			Group:      "Sizes",
			Editors:    "numeric",
			Options:    "",
			Default:    "",
		},
		models.T_AttributeInfo{
			Model:      "OrderItem",
			Collection: "OrderItem",
			Attribute:  "size",
			Type:       "decimal(10,2)",
			IsRequired: false,
			IsStatic:   true,
			Label:      "Size",
			Group:      "Sizes",
			Editors:    "numeric",
			Options:    "",
			Default:    "",
		},
	}

	return info
}

//--------------
// DefaultOrder
//--------------

// returns attribute of Order or nil
func (it *DefaultOrder) Get(attribute string) interface{} {
	switch strings.ToLower(attribute) {
	case "_id", "id":
		return it.id

	case "increment_id":
		return it.IncrementId

	case "status":
		return it.Status

	case "visitor_id":
		return it.VisitorId

	case "cart_id":
		return it.CartId

	case "shipping_address":
		return it.ShippingAddress

	case "billing_address":
		return it.BillingAddress

	case "customer_email":
		return it.CustomerEmail

	case "customer_name":
		return it.CustomerName

	case "payment_method":
		return it.PaymentMethod

	case "shipping_method":
		return it.ShippingMethod

	case "subtotal":
		return it.Subtotal

	case "discount":
		return it.Discount

	case "tax_amount":
		return it.TaxAmount

	case "shipping_amount":
		return it.ShippingAmount

	case "grand_total":
		return it.GrandTotal

	case "created_at":
		return it.CreatedAt

	case "updaed_at":
		return it.UpdatedAt
	}

	return nil
}

// sets attribute to Order object, returns error on problems
func (it *DefaultOrder) Set(attribute string, value interface{}) error {
	attribute = strings.ToLower(attribute)

	switch attribute {
	case "_id", "id":
		it.SetId(utils.InterfaceToString(value))

	case "increment_id":
		it.IncrementId = utils.InterfaceToString(value)

	case "status":
		it.Status = utils.InterfaceToString(value)

	case "visitor_id":
		it.VisitorId = utils.InterfaceToString(value)

	case "cart_id":
		it.CartId = utils.InterfaceToString(value)

	case "customer_email":
		it.CustomerEmail = utils.InterfaceToString(value)

	case "customer_name":
		it.CustomerName = utils.InterfaceToString(value)

	case "billing_address":
		it.BillingAddress = utils.InterfaceToMap(value)

	case "shipping_address":
		it.ShippingAddress = utils.InterfaceToMap(value)

	case "payment_method":
		it.PaymentMethod = utils.InterfaceToString(value)

	case "shipping_method":
		it.ShippingMethod = utils.InterfaceToString(value)

	case "subtotal":
		it.Subtotal = utils.InterfaceToFloat64(value)

	case "discount":
		it.Discount = utils.InterfaceToFloat64(value)

	case "tax_amount":
		it.TaxAmount = utils.InterfaceToFloat64(value)

	case "shipping_amount":
		it.ShippingAmount = utils.InterfaceToFloat64(value)

	case "grand_total":
		it.GrandTotal = utils.InterfaceToFloat64(value)

	case "created_at":
		it.CreatedAt = utils.InterfaceToTime(value)

	case "updated_at":
		it.UpdatedAt = utils.InterfaceToTime(value)

	default:
		return errors.New("unknown attribute: " + attribute)
	}

	return nil
}

// fills Order attributes with values provided in input map
func (it *DefaultOrder) FromHashMap(input map[string]interface{}) error {

	for attribute, value := range input {
		if err := it.Set(attribute, value); err != nil {
			return err
		}
	}

	return nil
}

// makes map from Order attribute values
func (it *DefaultOrder) ToHashMap() map[string]interface{} {

	result := make(map[string]interface{})

	result["_id"] = it.id

	result["increment_id"] = it.Get("increment_id")
	result["status"] = it.Get("status")

	result["visitor_id"] = it.Get("visitor_id")
	result["cart_id"] = it.Get("cart_id")

	result["customer_email"] = it.Get("customer_email")
	result["customer_name"] = it.Get("customer_name")

	result["shipping_address"] = it.Get("shipping_address")
	result["billing_address"] = it.Get("billing_address")

	result["payment_method"] = it.Get("payment_method")
	result["shipping_method"] = it.Get("shipping_method")

	result["subtotal"] = it.Get("subtotal")
	result["discount"] = it.Get("discount")
	result["tax_amount"] = it.Get("tax_amount")
	result["shipping_amount"] = it.Get("shipping_amount")
	result["grand_total"] = it.Get("grand_total")

	result["created_at"] = it.Get("created_at")
	result["updaed_at"] = it.Get("updaed_at")

	return result
}

// describes attributes of Order model
func (it *DefaultOrder) GetAttributesInfo() []models.T_AttributeInfo {

	info := []models.T_AttributeInfo{
		models.T_AttributeInfo{
			Model:      "Order",
			Collection: "Order",
			Attribute:  "_id",
			Type:       "id",
			IsRequired: false,
			IsStatic:   true,
			Label:      "ID",
			Group:      "General",
			Editors:    "not_editable",
			Options:    "",
			Default:    "",
		},
		models.T_AttributeInfo{
			Model:      "Order",
			Collection: "Order",
			Attribute:  "increment_id",
			Type:       "varchar(50)",
			IsRequired: true,
			IsStatic:   true,
			Label:      "Increment ID",
			Group:      "General",
			Editors:    "not_editable",
			Options:    "",
			Default:    "",
		},
		models.T_AttributeInfo{
			Model:      "Order",
			Collection: "Order",
			Attribute:  "status",
			Type:       "varchar(50)",
			IsRequired: true,
			IsStatic:   true,
			Label:      "Status",
			Group:      "General",
			Editors:    "selector",
			Options:    "pending,canceled,complete",
			Default:    "pending",
		},
		models.T_AttributeInfo{
			Model:      "Order",
			Collection: "Order",
			Attribute:  "visitor_id",
			Type:       "id",
			IsRequired: true,
			IsStatic:   true,
			Label:      "Visitor",
			Group:      "General",
			Editors:    "model_selector",
			Options:    "model: visitor",
			Default:    "",
		},
		models.T_AttributeInfo{
			Model:      "Order",
			Collection: "Order",
			Attribute:  "customer_email",
			Type:       "varchar(100)",
			IsRequired: true,
			IsStatic:   true,
			Label:      "Customer Email",
			Group:      "General",
			Editors:    "line_text",
			Options:    "",
			Default:    "",
		},
		models.T_AttributeInfo{
			Model:      "Order",
			Collection: "Order",
			Attribute:  "customer_name",
			Type:       "varchar(100)",
			IsRequired: true,
			IsStatic:   true,
			Label:      "Customer Name",
			Group:      "General",
			Editors:    "line_text",
			Options:    "",
			Default:    "",
		},
		models.T_AttributeInfo{
			Model:      "Order",
			Collection: "Order",
			Attribute:  "shipping_address",
			Type:       "text",
			IsRequired: true,
			IsStatic:   true,
			Label:      "Shipping Address",
			Group:      "General",
			Editors:    "visitor_address",
			Options:    "",
			Default:    "",
		},
		models.T_AttributeInfo{
			Model:      "Order",
			Collection: "Order",
			Attribute:  "billing_address",
			Type:       "text",
			IsRequired: true,
			IsStatic:   true,
			Label:      "Customer Name",
			Group:      "General",
			Editors:    "visitor_address",
			Options:    "",
			Default:    "",
		},
		models.T_AttributeInfo{
			Model:      "Order",
			Collection: "Order",
			Attribute:  "payment_method",
			Type:       "varchar(100)",
			IsRequired: true,
			IsStatic:   true,
			Label:      "Payment Method",
			Group:      "General",
			Editors:    "model_selector",
			Options:    "model: payments",
			Default:    "",
		},
		models.T_AttributeInfo{
			Model:      "Order",
			Collection: "Order",
			Attribute:  "shipping_method",
			Type:       "varchar(100)",
			IsRequired: true,
			IsStatic:   true,
			Label:      "Shipping Method",
			Group:      "General",
			Editors:    "model_selector",
			Options:    "model: shipping",
			Default:    "",
		},
		models.T_AttributeInfo{
			Model:      "Order",
			Collection: "Order",
			Attribute:  "subtotal",
			Type:       "decimal(10,2)",
			IsRequired: true,
			IsStatic:   true,
			Label:      "Totals",
			Group:      "General",
			Editors:    "not_editable",
			Options:    "",
			Default:    "",
		},
		models.T_AttributeInfo{
			Model:      "Order",
			Collection: "Order",
			Attribute:  "discount",
			Type:       "decimal(10,2)",
			IsRequired: true,
			IsStatic:   true,
			Label:      "Discount",
			Group:      "Totals",
			Editors:    "not_editable",
			Options:    "",
			Default:    "",
		},
		models.T_AttributeInfo{
			Model:      "Order",
			Collection: "Order",
			Attribute:  "tax_amount",
			Type:       "decimal(10,2)",
			IsRequired: true,
			IsStatic:   true,
			Label:      "Tax Amount",
			Group:      "Totals",
			Editors:    "not_editable",
			Options:    "",
			Default:    "",
		},
		models.T_AttributeInfo{
			Model:      "Order",
			Collection: "Order",
			Attribute:  "shipping_amount",
			Type:       "decimal(10,2)",
			IsRequired: true,
			IsStatic:   true,
			Label:      "Shipping Amount",
			Group:      "Totals",
			Editors:    "not_editable",
			Options:    "",
			Default:    "",
		},
		models.T_AttributeInfo{
			Model:      "Order",
			Collection: "Order",
			Attribute:  "grand_total",
			Type:       "decimal(10,2)",
			IsRequired: true,
			IsStatic:   true,
			Label:      "Grand Total",
			Group:      "Totals",
			Editors:    "not_editable",
			Options:    "",
			Default:    "",
		},
		models.T_AttributeInfo{
			Model:      "Order",
			Collection: "Order",
			Attribute:  "created_at",
			Type:       "datetime",
			IsRequired: true,
			IsStatic:   true,
			Label:      "Created At",
			Group:      "General",
			Editors:    "not_editable",
			Options:    "",
			Default:    "",
		},
		models.T_AttributeInfo{
			Model:      "Order",
			Collection: "Order",
			Attribute:  "updated_at",
			Type:       "datetime",
			IsRequired: true,
			IsStatic:   true,
			Label:      "Updated At",
			Group:      "General",
			Editors:    "not_editable",
			Options:    "",
			Default:    "",
		},
	}

	return info
}
