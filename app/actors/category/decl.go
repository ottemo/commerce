package category

import (
	"github.com/ottemo/foundation/app/models/category"

	"github.com/ottemo/foundation/db"
)

const (
	COLLECTION_NAME_CATEGORY                  = "category"
	COLLECTION_NAME_CATEGORY_PRODUCT_JUNCTION = "category_product"
)

type DefaultCategory struct {
	id string

	Name       string
	Parent     category.I_Category
	Path       string
	ProductIds []string
}

type DefaultCategoryCollection struct {
	listCollection     db.I_DBCollection
	listExtraAtributes []string
}
