package product

import (
	"github.com/ottemo/foundation/db"
	"github.com/ottemo/foundation/app/actors/attributes"
)

type DefaultProduct struct {
	id string

	Sku  string
	Name string

	Description string
	DefaultImage string

	Price float64

	*attributes.CustomAttributes

	listCollection db.I_DBCollection
	listExtraAtributes []string
}
