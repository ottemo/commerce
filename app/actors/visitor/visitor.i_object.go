package visitor

import (
	"errors"
	"strings"

	"github.com/ottemo/foundation/app/models"
	"github.com/ottemo/foundation/app/models/visitor"
	"github.com/ottemo/foundation/app/utils"
)

// returns object attribute value or nil
func (it *DefaultVisitor) Get(attribute string) interface{} {
	switch strings.ToLower(attribute) {
	case "_id", "id":
		return it.id
	case "fname", "first_name":
		return it.FirstName
	case "lname", "last_name":
		return it.LastName
	case "billing_address_id":
		return it.BillingAddress.GetId()
	case "shipping_address_id":
		return it.ShippingAddress.GetId()
	case "billing_address":
		return it.BillingAddress
	case "shipping_address":
		return it.ShippingAddress
	case "validate":
		return it.ValidateKey
	case "facebook_id":
		return it.FacebookId
	case "google_id":
		return it.GoogleId
	case "birthday":
		return it.Birthday
	case "is_admin":
		return it.IsAdmin
	case "created_at":
		return it.IsAdmin
	}

	return nil
}

// sets attribute value to object or returns error
func (it *DefaultVisitor) Set(attribute string, value interface{}) error {
	attribute = strings.ToLower(attribute)

	switch attribute {
	case "_id", "id":
		it.id = utils.InterfaceToString(value)
	case "email", "e_mail", "e-mail":
		it.Email = utils.InterfaceToString(value)
	case "fname", "first_name":
		it.FirstName = utils.InterfaceToString(value)
	case "lname", "last_name":
		it.LastName = utils.InterfaceToString(value)
	case "password", "passwd":
		it.SetPassword(utils.InterfaceToString(value))
	case "validate":
		it.ValidateKey = utils.InterfaceToString(value)
	case "facebook_id":
		it.FacebookId = utils.InterfaceToString(value)
	case "google_id":
		it.GoogleId = utils.InterfaceToString(value)
	case "birthday":
		it.Birthday = utils.InterfaceToTime(value)
	case "is_admin":
		it.IsAdmin = utils.InterfaceToBool(value)
	case "created_at":
		it.CreatedAt = utils.InterfaceToTime(value)

	// only address id coming - trying to get it from DB
	case "billing_address_id", "shipping_address_id":
		address, err := visitor.LoadVisitorAddressById(utils.InterfaceToString(value))
		if err != nil {
			return err
		}

		if address != nil && address.GetId() != "" {

			if attribute == "billing_address_id" {
				it.BillingAddress = address
			} else {
				it.ShippingAddress = address
			}

		}

	// address with details coming
	case "billing_address", "shipping_address":
		switch value := value.(type) {

		// we have already have structure
		case visitor.I_VisitorAddress:
			if attribute == "billing_address" {
				it.BillingAddress = value
			} else {
				it.ShippingAddress = value
			}

		// we have sub-map, supposedly I_VisitorAddress capable
		case map[string]interface{}:
			addressModel, err := visitor.GetVisitorAddressModel()

			err = addressModel.FromHashMap(value)
			if err != nil {
				return err
			}

			if attribute == "billing_address" {
				it.BillingAddress = addressModel
			} else {
				it.ShippingAddress = addressModel
			}

		default:
			return errors.New("unsupported billing or shipping address value")
		}
	}
	return nil
}

// fills object attributes from map[string]interface{}
func (it *DefaultVisitor) FromHashMap(input map[string]interface{}) error {

	for attribute, value := range input {
		if err := it.Set(attribute, value); err != nil {
			return err
		}
	}

	return nil
}

// represents object as map[string]interface{}
func (it *DefaultVisitor) ToHashMap() map[string]interface{} {

	result := make(map[string]interface{})

	result["_id"] = it.id

	result["email"] = it.Email
	result["first_name"] = it.FirstName
	result["last_name"] = it.LastName

	result["birthday"] = it.Birthday
	result["created_at"] = it.CreatedAt

	result["billing_address"] = nil
	result["shipping_address"] = nil

	//result["billing_address_id"] = it.BillingAddressId
	//result["shipping_address_id"] = it.ShippingAddressId

	if it.BillingAddress != nil {
		result["billing_address"] = it.BillingAddress.ToHashMap()
	}

	if it.ShippingAddress != nil {
		result["shipping_address"] = it.ShippingAddress.ToHashMap()
	}

	return result
}

// returns information about object attributes
func (it *DefaultVisitor) GetAttributesInfo() []models.T_AttributeInfo {

	info := []models.T_AttributeInfo{
		models.T_AttributeInfo{
			Model:      visitor.MODEL_NAME_VISITOR,
			Collection: COLLECTION_NAME_VISITOR,
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
			Model:      visitor.MODEL_NAME_VISITOR,
			Collection: COLLECTION_NAME_VISITOR,
			Attribute:  "email",
			Type:       "text",
			IsRequired: true,
			IsStatic:   true,
			Label:      "E-mail",
			Group:      "General",
			Editors:    "line_text",
			Options:    "",
			Default:    "",
		},
		models.T_AttributeInfo{
			Model:      visitor.MODEL_NAME_VISITOR,
			Collection: COLLECTION_NAME_VISITOR,
			Attribute:  "first_name",
			Type:       "text",
			IsRequired: true,
			IsStatic:   true,
			Label:      "First Name",
			Group:      "General",
			Editors:    "line_text",
			Options:    "",
			Default:    "",
		},
		models.T_AttributeInfo{
			Model:      visitor.MODEL_NAME_VISITOR,
			Collection: COLLECTION_NAME_VISITOR,
			Attribute:  "last_name",
			Type:       "text",
			IsRequired: true,
			IsStatic:   true,
			Label:      "Last Name",
			Group:      "General",
			Editors:    "line_text",
			Options:    "",
			Default:    "",
		},
		models.T_AttributeInfo{
			Model:      visitor.MODEL_NAME_VISITOR,
			Collection: COLLECTION_NAME_VISITOR,
			Attribute:  "password",
			Type:       "text",
			IsRequired: false,
			IsStatic:   true,
			Label:      "Password",
			Group:      "Password",
			Editors:    "password",
			Options:    "",
			Default:    "",
		},
		models.T_AttributeInfo{
			Model:      visitor.MODEL_NAME_VISITOR,
			Collection: COLLECTION_NAME_VISITOR,
			Attribute:  "billing_address_id",
			Type:       "text",
			IsRequired: false,
			IsStatic:   true,
			Label:      "Billing Address",
			Group:      "General",
			Editors:    "model_selector",
			Options:    "model=VisitorAddress",
			Default:    "",
		},
		models.T_AttributeInfo{
			Model:      visitor.MODEL_NAME_VISITOR,
			Collection: COLLECTION_NAME_VISITOR,
			Attribute:  "shipping_address_id",
			Type:       "text",
			IsRequired: false,
			IsStatic:   true,
			Label:      "Shipping Address",
			Group:      "General",
			Editors:    "model_selector",
			Options:    "model:VisitorAddress",
			Default:    "",
		},
	}

	return info
}
