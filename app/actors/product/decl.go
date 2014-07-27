package product

import (
	"github.com/ottemo/foundation/db"
	"github.com/ottemo/foundation/app/helpers/attributes"
)

type DefaultProduct struct {
	id string

	Sku  string
	Name string

	ShortDescription string
	Description string

	DefaultImage string

	Price float64

	Weight float64
	Size float64

	*attributes.CustomAttributes

	listCollection db.I_DBCollection
	listExtraAtributes []string
}
