package category

import(
	"github.com/ottemo/foundation/models"
	"github.com/ottemo/foundation/models/product"
)

type ICategory interface {

	GetName() string

	GetParent() []ICategory
	GetProducts() []product.IProduct

	AddProduct(ProductId string) error
	RemoveProduct(ProductId string) error

	models.IModel
	models.IObject
	models.IStorable
	models.IMapable
}
