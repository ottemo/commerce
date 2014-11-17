// Package product represents abstraction of business layer product object
package product

import (
	"github.com/ottemo/foundation/app/models"
)

// Package global constants
const (
	MODEL_NAME_PRODUCT            = "Product"
	MODEL_NAME_PRODUCT_COLLECTION = "ProductCollection"
)

// I_Product represents interface to access business layer implementation of product object
type I_Product interface {
	GetSku() string
	GetName() string

	GetShortDescription() string
	GetDescription() string

	GetDefaultImage() string

	GetPrice() float64
	GetWeight() float64

	ApplyOptions(map[string]interface{}) error
	GetOptions() map[string]interface{}

	models.I_Model
	models.I_Object
	models.I_Storable
	models.I_Media
	models.I_Listable
	models.I_CustomAttributes
}

// I_ProductCollection represents interface to access business layer implementation of product collection
type I_ProductCollection interface {
	ListProducts() []I_Product

	models.I_Collection
}
