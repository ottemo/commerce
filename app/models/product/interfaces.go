package product

import (
	"github.com/ottemo/foundation/app/models"
)

const (
	MODEL_NAME_PRODUCT = "Product"
	MODEL_NAME_PRODUCT_COLLECTION = "ProductCollection"
)

type I_Product interface {
	GetSku() string
	GetName() string

	GetShortDescription() string
	GetDescription() string

	GetDefaultImage() string

	GetPrice() float64

	GetWeight() float64
	GetSize() float64

	models.I_Model
	models.I_Object
	models.I_Storable
	models.I_Media

	models.I_CustomAttributes
}

type I_ProductCollection interface {
	ListProducts() []I_Product

	models.I_Collection
}
