package default_product

import (
	"github.com/ottemo/foundation/models/custom_attributes"
	"github.com/ottemo/foundation/database"
)

type DefaultProductModel struct {
	id string

	Sku string
	Name string

	*custom_attributes.CustomAttributes
}
