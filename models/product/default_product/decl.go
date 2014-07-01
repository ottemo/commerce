package default_product

import (
	"github.com/ottemo/foundation/models/custom_attributes"
)

type DefaultProductModel struct {
	id string

	Sku string
	Name string

	Description string

	DefaultImage string

	Price float64

	*custom_attributes.CustomAttributes
}
