// Package models represents abstraction of business layer object and basic access interfaces for it
package models

import (
	"github.com/ottemo/foundation/db"
)

// I_Model represents interface for basic business layer implementation object
type I_Model interface {
	GetModelName() string
	GetImplementationName() string
	New() (I_Model, error)
}

// I_Storable represents interface load/store business layer implementation object from database
type I_Storable interface {
	GetId() string
	SetId(string) error

	Save() error
	Load(id string) error
	Delete() error
}

// I_Object represents interface to access business layer implementation object via get/set functions
type I_Object interface {
	Get(attribute string) interface{}
	Set(attribute string, value interface{}) error

	FromHashMap(hashMap map[string]interface{}) error
	ToHashMap() map[string]interface{}

	GetAttributesInfo() []T_AttributeInfo
}

// I_Listable represents interface to access business layer implementation collection via object instance
type I_Listable interface {
	GetCollection() I_Collection
}

// I_CustomAttributes represents interface to access business layer implementation object custom attributes
type I_CustomAttributes interface {
	GetCustomAttributeCollectionName() string

	AddNewAttribute(newAttribute T_AttributeInfo) error
	RemoveAttribute(attributeName string) error
}

// I_Media represents interface to access business layer implementation object assigned media resources
type I_Media interface {
	AddMedia(mediaType string, mediaName string, content []byte) error
	RemoveMedia(mediaType string, mediaName string) error

	ListMedia(mediaType string) ([]string, error)

	GetMedia(mediaType string, mediaName string) ([]byte, error)
	GetMediaPath(mediaType string) (string, error)
}

// I_Collection represents interface to access business layer implementation collection
type I_Collection interface {
	GetDBCollection() db.I_DBCollection

	List() ([]T_ListItem, error)

	ListAddExtraAttribute(attribute string) error

	ListFilterAdd(attribute string, operator string, value interface{}) error
	ListFilterReset() error

	ListLimit(offset int, limit int) error
}
