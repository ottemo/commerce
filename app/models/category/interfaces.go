// Package cart represents abstraction of business layer category object
package category

import (
	"github.com/ottemo/foundation/app/models"
	"github.com/ottemo/foundation/app/models/product"
)

// Package global constants
const (
	MODEL_NAME_CATEGORY            = "Category"
	MODEL_NAME_CATEGORY_COLLECTION = "CategoryCollection"
)

// I_Category represents interface to access business layer implementation of category object
type I_Category interface {
	GetName() string

	GetParent() I_Category

	GetProductIds() []string
	GetProductsCollection() product.I_ProductCollection
	GetProducts() []product.I_Product

	AddProduct(productId string) error
	RemoveProduct(productId string) error

	models.I_Model
	models.I_Object
	models.I_Storable
	models.I_Listable
}

// I_CategoryCollection represents interface to access business layer implementation of category collection
type I_CategoryCollection interface {
	ListCategories() []I_Category

	models.I_Collection
}
