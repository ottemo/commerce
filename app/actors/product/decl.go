package product

import (
	"github.com/ottemo/foundation/app/helpers/attributes"
	"github.com/ottemo/foundation/app/helpers/listable"
)

const (
	DB_COLLECTION_NAME_PRODUCT = "product"
)

type DefaultProduct struct {
	id string

	Sku  string
	Name string

	ShortDescription string
	Description      string

	DefaultImage string

	Price float64

	Weight float64
	Size   float64

	*attributes.CustomAttributes
}

type DefaultProductCollection struct {
	*listable.ListableHelper
}
