package category

import (
	"github.com/ottemo/foundation/app/models"
	"github.com/ottemo/foundation/app/models/product"
)

const (
	MODEL_NAME_CATEGORY            = "Category"
	MODEL_NAME_CATEGORY_COLLECTION = "CategoryCollection"
)

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

type I_CategoryCollection interface {
	ListCategories() []I_Category

	models.I_Collection
}
