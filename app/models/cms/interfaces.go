package cms

import (
	"github.com/ottemo/foundation/app/models"
)

const (
	MODEL_NAME_CMS_PAGE             = "CMSPage"
	MODEL_NAME_CMS_PAGE_COLLECTION  = "CMSPageCollection"
	MODEL_NAME_CMS_BLOCK            = "CMSBlock"
	MODEL_NAME_CMS_BLOCK_COLLECTION = "CMSBlockCollection"
)

type I_CMSPage interface {
	GetURL() string
	SetURL(string) error

	GetIdentifier() string
	SetIdentifier(string) error

	GetTitle() string
	SetTitle(string) error

	GetContent() string
	SetContent(string) error

	GetMetaKeywords() string
	SetMetaKeywords(string) error

	GetMetaDescription() string
	SetMetaDescription(string) error

	models.I_Model
	models.I_Object
	models.I_Storable
	models.I_Listable
}

type I_CMSPageCollection interface {
	ListCMSPages() []I_CMSPage

	models.I_Collection
}

type I_CMSBlock interface {
	GetIdentifier() string
	SetIdentifier(string) error

	GetContent() string
	SetContent(string) error

	models.I_Model
	models.I_Object
	models.I_Storable
	models.I_Listable
}

type I_CMSBlockCollection interface {
	ListCMSBlocks() []I_CMSBlock

	models.I_Collection
}
