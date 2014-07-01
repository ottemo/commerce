package default_product

import (
	"strings"
	"github.com/ottemo/foundation/models"
)

func (it *DefaultProductModel) Get(attribute string) interface{} {
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

func (it *DefaultProductModel) Set(attribute string, value interface{}) error {
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
		it.Price = value.(float64)
	default:
		if err := it.CustomAttributes.Set(attribute, value); err != nil {
			return err
		}

	}

	return nil
}

func (it *DefaultProductModel) GetAttributesInfo() []models.T_AttributeInfo {
	staticInfo := []models.T_AttributeInfo {
		models.T_AttributeInfo {
			Model: "Product",
			Collection: "product",
			Attribute: "_id",
			Type: "text",
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
			Label: "Price",
			Group: "Prices",
			Editors: "price",
			Options: "",
			Default: "",
		},
	}

	dynamicInfo := it.CustomAttributes.GetAttributesInfo()

	return append(dynamicInfo, staticInfo...)
}
