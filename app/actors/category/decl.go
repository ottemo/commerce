package category

import (
	"github.com/ottemo/foundation/app/models/category"
	"github.com/ottemo/foundation/app/models/product"

	"github.com/ottemo/foundation/db"
)

const (
	CATEGORY_COLLECTION_NAME                  = "category"
	CATEGORY_PRODUCT_JUNCTION_COLLECTION_NAME = "category_product"
)

type DefaultCategory struct {
	id string

	Name     string
	Parent   category.I_Category
	Path	 string
	Products []product.I_Product


	listCollection db.I_DBCollection
	listExtraAtributes []string
}
