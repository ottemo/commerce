package category

import(
	"github.com/ottemo/foundation/models"
	"github.com/ottemo/foundation/models/product"
)

type I_Category interface {

	    GetName() string
	GetProducts() []product.I_Product

	models.I_Model
	models.I_Object
	models.I_Storable
	models.I_Mapable
}
