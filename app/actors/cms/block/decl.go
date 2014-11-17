// Package block is a default implementation of cms block related interfaces declared in
// "github.com/ottemo/foundation/app/models/csm" package
package block

import (
	"github.com/ottemo/foundation/db"
	"time"
)

// Package global constants
const (
	CMS_BLOCK_COLLECTION_NAME = "cms_block"
)

// DefaultCMSBlock is a default implementer of I_CMSBlock
type DefaultCMSBlock struct {
	id string

	Identifier string
	Content    string

	CreatedAt time.Time
	UpdatedAt time.Time
}

// DefaultCMSBlockCollection is a default implementer of I_CMSBlockCollection
type DefaultCMSBlockCollection struct {
	listCollection     db.I_DBCollection
	listExtraAtributes []string
}
