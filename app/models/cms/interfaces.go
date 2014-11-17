// Package cms represents abstraction of business layer cms page and cms block objects
package cms

import (
	"github.com/ottemo/foundation/app/models"
)

// Package global constants
const (
	MODEL_NAME_CMS_PAGE             = "CMSPage"
	MODEL_NAME_CMS_PAGE_COLLECTION  = "CMSPageCollection"
	MODEL_NAME_CMS_BLOCK            = "CMSBlock"
	MODEL_NAME_CMS_BLOCK_COLLECTION = "CMSBlockCollection"
)

// I_CMSPage represents interface to access business layer implementation of cms page object
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

// I_CMSPageCollection represents interface to access business layer implementation of cms page collection
type I_CMSPageCollection interface {
	ListCMSPages() []I_CMSPage

	models.I_Collection
}

// I_CMSBlock represents interface to access business layer implementation of cms block object
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

// I_CMSBlockCollection represents interface to access business layer implementation of cms block collection
type I_CMSBlockCollection interface {
	ListCMSBlocks() []I_CMSBlock

	models.I_Collection
}
