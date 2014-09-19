package order

import (
	"errors"
	"strings"

	"github.com/ottemo/foundation/app/models"
	"github.com/ottemo/foundation/utils"

	"github.com/ottemo/foundation/app/models/order"
)

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
		return it.Options

	case "price":
		return it.Price

	case "weight":
		return it.Weight
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
			it.Options = mapValue
		} else {
			return errors.New("options should me map[string]interface{} type")
		}

	case "price":
		it.Price = utils.InterfaceToFloat64(value)

	case "weight":
		it.Weight = utils.InterfaceToFloat64(value)

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
	}

	return info
}
