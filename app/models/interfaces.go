package models

import (
	"github.com/ottemo/foundation/db"
)

type I_Model interface {
	GetModelName() string
	GetImplementationName() string
	New() (I_Model, error)
}

type I_Storable interface {
	GetId() string
	SetId(string) error

	Save() error
	Load(id string) error
	Delete() error
}

type I_Object interface {
	Get(attribute string) interface{}
	Set(attribute string, value interface{}) error

	FromHashMap(hashMap map[string]interface{}) error
	ToHashMap() map[string]interface{}

	GetAttributesInfo() []T_AttributeInfo
}

type I_Listable interface {
	List() ([]T_ListItem, error)

	ListAddExtraAttribute(attribute string) error

	ListFilterAdd(attribute string, operator string, value interface{}) error
	ListFilterReset() error

	ListLimit(offset int, limit int) error
}

type I_CustomAttributes interface {
	AddNewAttribute(newAttribute T_AttributeInfo) error
	RemoveAttribute(attributeName string) error
}

type I_Media interface {
	AddMedia(mediaType string, mediaName string, content []byte) error
	RemoveMedia(mediaType string, mediaName string) error

	ListMedia(mediaType string) ([]string, error)

	GetMedia(mediaType string, mediaName string) ([]byte, error)
	GetMediaPath(mediaType string) (string, error)
}

type I_Collection interface {
	GetDBCollection() db.I_DBCollection
	I_Listable
}

type T_ListItem struct {
	Id    string
	Name  string
	Image string
	Desc  string

	Extra map[string]interface{}
}

type T_AttributeInfo struct {
	Model      string
	Collection string
	Attribute  string
	Type       string
	Label      string
	IsRequired bool
	IsStatic   bool
	Group      string
	Editors    string
	Options    string
	Default    string
	Validators string
	IsLayered  bool
}
