// Package page is a default implementation of cms page related interfaces declared in
// "github.com/ottemo/foundation/app/models/csm" package
package page

import (
	"github.com/ottemo/foundation/db"
	"time"
)

// Package global constants
const (
	CMS_PAGE_COLLECTION_NAME = "cms_page"
)

// DefaultCMSPage is a default implementer of I_CMSPage
type DefaultCMSPage struct {
	id string

	URL string

	Identifier string

	Title   string
	Content string

	MetaKeywords    string
	MetaDescription string

	CreatedAt time.Time
	UpdatedAt time.Time
}

// DefaultCMSPageCollection is a default implementer of I_CMSPageCollection
type DefaultCMSPageCollection struct {
	listCollection     db.I_DBCollection
	listExtraAtributes []string
}
