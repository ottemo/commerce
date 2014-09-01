package block

import (
	"time"
	"github.com/ottemo/foundation/db"
)

const (
	CMS_BLOCK_COLLECTION_NAME = "cms_block"
)

type DefaultCMSBlock struct {
	id string

	Identifier string
	Content string

	CreatedAt time.Time
	UpdatedAt time.Time

	listCollection     db.I_DBCollection
	listExtraAtributes []string
}
