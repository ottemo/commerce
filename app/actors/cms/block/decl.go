// Package block is a default implementation of cms block related interfaces declared in
// "github.com/ottemo/foundation/app/models/csm" package
package block

import (
	"github.com/ottemo/foundation/db"
	"time"
)

// Package global constants
const (
	ConstCmsBlockCollectionName = "cms_block"
)

// DefaultCMSBlock is a default implementer of InterfaceCMSBlock
type DefaultCMSBlock struct {
	id string

	Identifier string
	Content    string

	CreatedAt time.Time
	UpdatedAt time.Time
}

// DefaultCMSBlockCollection is a default implementer of InterfaceCMSBlockCollection
type DefaultCMSBlockCollection struct {
	listCollection     db.InterfaceDBCollection
	listExtraAtributes []string
}
