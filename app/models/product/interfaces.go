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
	GetSku() string
	GetName() string

	GetShortDescription() string
	GetDescription() string

	GetDefaultImage() string

	GetPrice() float64
	GetWeight() float64

	ApplyOptions(map[string]interface{}) error
	GetOptions() map[string]interface{}

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
