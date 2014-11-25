// Package product represents abstraction of business layer product object
package product

import (
	"github.com/ottemo/foundation/app/models"
)

// Package global constants
const (
	ConstModelNameProduct           = "Product"
	ConstModelNameProductCollection = "ProductCollection"
)

// InterfaceProduct represents interface to access business layer implementation of product object
type InterfaceProduct interface {
	GetEnabled() bool

	GetSku() string
	GetName() string

	GetShortDescription() string
	GetDescription() string

	GetDefaultImage() string

	GetPrice() float64
	GetWeight() float64

	GetQty() float64

	GetAppliedOptions() map[string]interface{}
	GetOptions() map[string]interface{}

	ApplyOptions(map[string]interface{}) error

	models.InterfaceModel
	models.InterfaceObject
	models.InterfaceStorable
	models.InterfaceMedia
	models.InterfaceListable
	models.InterfaceCustomAttributes
}

// InterfaceProductCollection represents interface to access business layer implementation of product collection
type InterfaceProductCollection interface {
	ListProducts() []InterfaceProduct

	models.InterfaceCollection
}

// InterfaceStock represents interface to access business layer implementation of stock management
type InterfaceStock interface {
	SetProductQty(productID string, options map[string]interface{}, qty float64) error
	GetProductQty(productID string, options map[string]interface{}) float64

	RemoveProductQty(productID string, options map[string]interface{}) error
	UpdateProductQty(productID string, options map[string]interface{}, deltaQty float64) error
}
