// Package category is a default implementation of interfaces declared in
// "github.com/ottemo/foundation/app/models/category" package
package category

import (
	"github.com/ottemo/foundation/app/models/category"
	"github.com/ottemo/foundation/db"
)

// Package global constants
const (
	COLLECTION_NAME_CATEGORY                  = "category"
	COLLECTION_NAME_CATEGORY_PRODUCT_JUNCTION = "category_product"
)

// DefaultCategory is a default implementer of I_Category
type DefaultCategory struct {
	id string

	Name       string
	Parent     category.I_Category
	Path       string
	ProductIds []string
}

// DefaultCategoryCollection is a default implementer of I_CategoryCollection
type DefaultCategoryCollection struct {
	listCollection     db.I_DBCollection
	listExtraAtributes []string
}
