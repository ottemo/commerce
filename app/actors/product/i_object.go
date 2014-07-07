package product

import (
	"strings"
	"strconv"
	"errors"

	"github.com/ottemo/foundation/app/models"
)

func (it *DefaultProduct) Get(attribute string) interface{} {
	switch strings.ToLower(attribute) {
	case "_id", "id":
		return it.id
	case "sku":
		return it.Sku
	case "name":
		return it.Name
	case "description":
		return it.Description
	case "default_image", "defaultimage":
		return it.DefaultImage
	case "price":
		return it.Price
	default:
		return it.CustomAttributes.Get(attribute)
	}

	return nil
}

func (it *DefaultProduct) Set(attribute string, value interface{}) error {
	switch strings.ToLower(attribute) {
	case "_id", "id":
		it.id = value.(string)
	case "sku":
		it.Sku = value.(string)
	case "name":
		it.Name = value.(string)
	case "description":
		it.Description = value.(string)
	case "default_image", "defaultimage":
		it.DefaultImage = value.(string)

	case "price":
		switch value := value.(type) {
		case float64:
			it.Price = value
		case string:
			newPrice, err := strconv.ParseFloat( value , 64)
			if err != nil { return err }

			it.Price = newPrice
		default:
			return errors.New("wrong price format")
		}

	default:
		err := it.CustomAttributes.Set(attribute, value)
		if err != nil { return err }
	}

	return nil
}

func (it *DefaultProduct) GetAttributesInfo() []models.T_AttributeInfo {
	result := []models.T_AttributeInfo {
		models.T_AttributeInfo {
			Model: "Product",
			Collection: "product",
			Attribute: "_id",
			Type: "text",
			IsRequired: false,
			IsStatic: true,
			Label: "ID",
			Group: "General",
			Editors: "not_editable",
			Options: "",
			Default: "",
		},
		models.T_AttributeInfo {
			Model: "Product",
			Collection: "product",
			Attribute: "sku",
			Type: "text",
			IsRequired: true,
			IsStatic: true,
			Label: "SKU",
			Group: "General",
			Editors: "line_text",
			Options: "",
			Default: "",
		},
		models.T_AttributeInfo {
			Model: "Product",
			Collection: "product",
			Attribute: "name",
			Type: "text",
			IsRequired: true,
			IsStatic: true,
			Label: "Name",
			Group: "General",
			Editors: "line_text",
			Options: "",
			Default: "",
		},
		models.T_AttributeInfo {
			Model: "Product",
			Collection: "product",
			Attribute: "description",
			Type: "text",
			IsRequired: false,
			IsStatic: true,
			Label: "Description",
			Group: "General",
			Editors: "multiline_text",
			Options: "",
			Default: "",
		},
		models.T_AttributeInfo {
			Model: "Product",
			Collection: "product",
			Attribute: "default_image",
			Type: "text",
			IsRequired: false,
			IsStatic: true,
			Label: "DefaultImage",
			Group: "Pictures",
			Editors: "image_selector",
			Options: "",
			Default: "",
		},
		models.T_AttributeInfo {
			Model: "Product",
			Collection: "product",
			Attribute: "price",
			Type: "text",
			IsRequired: false,
			IsStatic: true,
			Label: "Price",
			Group: "Prices",
			Editors: "price",
			Options: "",
			Default: "",
		},
	}

	dynamicInfo := it.CustomAttributes.GetAttributesInfo()

	for _, dynamicAttribute := range dynamicInfo {
		result = append(result, dynamicAttribute)
	}

	return result
}
