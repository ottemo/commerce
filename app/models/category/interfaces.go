package category

import (
	"github.com/ottemo/foundation/app/models"
	"github.com/ottemo/foundation/app/models/product"
)

type I_Category interface {
	GetName() string

	GetParent() I_Category
	GetProducts() []product.I_Product

	AddProduct(productId string) error
	RemoveProduct(productId string) error

	models.I_Model
	models.I_Object
	models.I_Storable
	models.I_Listable
}
