// Package product is a implementation of interfaces declared in
// "github.com/ottemo/foundation/app/models/product" package
package product

import (
	"github.com/ottemo/foundation/app/helpers/attributes"
	"github.com/ottemo/foundation/db"
)

// Package global constants
const (
	ConstCollectionNameProduct = "product"
)

// DefaultProduct is a default implementer of InterfaceProduct
type DefaultProduct struct {
	id string

	Sku  string
	Name string

	ShortDescription string
	Description      string

	DefaultImage string

	Price float64

	Weight float64

	Options map[string]interface{}

	RelatedProductIds []string

	*attributes.CustomAttributes
}

// DefaultProductCollection is a default implementer of InterfaceProduct
type DefaultProductCollection struct {
	listCollection     db.InterfaceDBCollection
	listExtraAtributes []string
}
