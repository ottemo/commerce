package product

import (
	"github.com/ottemo/foundation/app/helpers/attributes"
	"github.com/ottemo/foundation/db"
)

const (
	COLLECTION_NAME_PRODUCT = "product"
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

	Options map[string]interface{}

	*attributes.CustomAttributes
}

type DefaultProductCollection struct {
	listCollection     db.I_DBCollection
	listExtraAtributes []string
}
