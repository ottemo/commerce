package category

import(
	"github.com/ottemo/foundation/app/models/product"
	"github.com/ottemo/foundation/app/models/category"
)

const (
	CATEGORY_COLLECTION_NAME = "category"
	CATEGORY_PRODUCT_JUNCTION_COLLECTION_NAME = "category_product"
)


type DefaultCategory struct {
	id string

	Name string
	Parent category.I_Category
	Products []product.I_Product
}
