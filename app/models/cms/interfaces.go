package cms

import (
	"github.com/ottemo/foundation/app/models"
)

const (
	CMS_PAGE_MODEL_NAME  = "CMSPage"
	CMS_BLOCK_MODEL_NAME = "CMSBlock"
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
