package default_product

import (
	"github.com/ottemo/foundation/models/custom_attributes"
)

type DefaultProductModel struct {
	id string

	Sku string
	Name string

	*custom_attributes.CustomAttributes
}
