package page

import (
	"time"
	"github.com/ottemo/foundation/db"
)

const (
	CMS_PAGE_COLLECTION_NAME = "cms_page"
)

type DefaultCMSPage struct {
	id string

	URL string

	Identifier string

	Title   string
	Content string

	MetaKeywords     string
	MetaDescription  string

	CreatedAt time.Time
	UpdatedAt time.Time

	listCollection     db.I_DBCollection
	listExtraAtributes []string
}
