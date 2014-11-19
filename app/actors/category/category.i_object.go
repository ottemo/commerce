package category

import (
	"strings"

	"github.com/ottemo/foundation/app/models"
	"github.com/ottemo/foundation/app/models/category"
	"github.com/ottemo/foundation/app/models/product"
	"github.com/ottemo/foundation/env"
)

//---------------------------------
// IMPLEMENTATION SPECIFIC METHODS
//---------------------------------

// updates path attribute of model
func (it *DefaultCategory) updatePath() {
	if it.GetID() == "" {
		it.Path = ""
	} else if it.Parent != nil {
		parentPath, ok := it.Parent.Get("path").(string)
		if ok {
			it.Path = parentPath + "/" + it.GetID()
		}
	} else {
		it.Path = "/" + it.GetID()
	}
}

//--------------------------
// INTERFACE IMPLEMENTATION
//--------------------------

func (it *DefaultCategory) Get(attribute string) interface{} {
	switch strings.ToLower(attribute) {
	case "_id", "id":
		return it.id

	case "name":
		return it.Name

	case "path":
		if it.Path == "" {
			it.updatePath()
		}
		return it.Path

	case "parent_id":
		if it.Parent != nil {
			return it.Parent.GetID()
		} else {
			return ""
		}

	case "parent":
		return it.Parent

	case "products":
		result := make([]map[string]interface{}, 0)

		for _, categoryProduct := range it.GetProducts() {
			result = append(result, categoryProduct.ToHashMap())
		}

		return result
	}

	return nil
}

func (it *DefaultCategory) Set(attribute string, value interface{}) error {
	attribute = strings.ToLower(attribute)

	switch attribute {
	case "_id", "id":
		it.id = value.(string)

	case "name":
		it.Name = value.(string)

	case "parent_id":
		if value, ok := value.(string); ok {
			value = strings.TrimSpace(value)
			if value != "" {
				model, err := models.GetModel("Category")
				if err != nil {
					return env.ErrorDispatch(err)
				}
				categoryModel, ok := model.(category.InterfaceCategory)
				if !ok {
					return env.ErrorNew("unsupported category model " + model.GetImplementationName())
				}

				err = categoryModel.Load(value)
				if err != nil {
					return env.ErrorDispatch(err)
				}

				selfID := it.GetID()
				if selfID != "" {
					parentPath, ok := categoryModel.Get("path").(string)
					if categoryModel.GetID() != selfID && ok && !strings.Contains(parentPath, selfID) {
						it.Parent = categoryModel
					} else {
						return env.ErrorNew("category can't have sub-category or itself as parent")
					}
				} else {
					it.Parent = categoryModel
				}
			} else {
				it.Parent = nil
			}
		} else {
			return env.ErrorNew("unsupported id specified")
		}
		it.updatePath()

	case "parent":
		switch value := value.(type) {
		case category.InterfaceCategory:
			it.Parent = value
		case string:
			it.Set("parent_id", value)
		default:
			env.ErrorNew("unsupported 'parent' value")
		}
		// path should be changed as well
		it.updatePath()

	case "products":
		switch typedValue := value.(type) {

		case []interface{}:
			for _, listItem := range typedValue {
				productID, ok := listItem.(string)
				if ok {
					productModel, err := product.LoadProductByID(productID)
					if err != nil {
						return env.ErrorDispatch(err)
					}

					it.ProductIds = append(it.ProductIds, productModel.GetID())
				}
			}

		case []product.InterfaceProduct:
			it.ProductIds = make([]string, 0)
			for _, productItem := range typedValue {
				it.ProductIds = append(it.ProductIds, productItem.GetID())
			}

		default:
			return env.ErrorNew("unsupported 'products' value")
		}
	}
	return nil
}

func (it *DefaultCategory) FromHashMap(input map[string]interface{}) error {

	for attribute, value := range input {
		if err := it.Set(attribute, value); err != nil {
			return env.ErrorDispatch(err)
		}
	}

	return nil
}

func (it *DefaultCategory) ToHashMap() map[string]interface{} {

	result := make(map[string]interface{})

	result["_id"] = it.id

	result["parent_id"] = it.Get("parent_id")
	result["name"] = it.Get("name")
	result["products"] = it.Get("products")
	result["path"] = it.Get("path")

	return result
}

func (it *DefaultCategory) GetAttributesInfo() []models.StructAttributeInfo {

	info := []models.StructAttributeInfo{
		models.StructAttributeInfo{
			Model:      "Category",
			Collection: "Category",
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
		models.StructAttributeInfo{
			Model:      "Category",
			Collection: "Category",
			Attribute:  "name",
			Type:       "text",
			IsRequired: true,
			IsStatic:   true,
			Label:      "Name",
			Group:      "General",
			Editors:    "line_text",
			Options:    "",
			Default:    "",
		},
		models.StructAttributeInfo{
			Model:      "Category",
			Collection: "Category",
			Attribute:  "parent_id",
			Type:       "id",
			IsRequired: false,
			IsStatic:   true,
			Label:      "Parent",
			Group:      "General",
			Editors:    "category_selector",
			Options:    "",
			Default:    "",
		},
		models.StructAttributeInfo{
			Model:      "Category",
			Collection: "Category",
			Attribute:  "products",
			Type:       "id",
			IsRequired: false,
			IsStatic:   true,
			Label:      "Products",
			Group:      "General",
			Editors:    "product_selector",
			Options:    "",
			Default:    "",
		},
	}

	return info
}
