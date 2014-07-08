package models

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
	Delete(Id string) error
}

type I_Object interface {
	Get(Attribute string) interface{}
	Set(Attribute string, Value interface{}) error

	FromHashMap(HashMap map[string]interface{}) error
	ToHashMap() map[string]interface{}

	GetAttributesInfo() []T_AttributeInfo
}

type I_Listable interface {
	List() ([]interface{}, error)
	ListFilterAdd(Attribute string, Operator string, Value interface{}) error
	ListFilterReset() error
}

type I_CustomAttributes interface {
	AddNewAttribute(newAttribute T_AttributeInfo) error
	RemoveAttribute(attributeName string) error
}

type I_Media interface {
	AddMedia(mediaType string, mediaName string, content []byte) error
	RemoveMedia(mediaType string, mediaName string) error
	ListMedia(mediaType string ) ([]string, error)
	GetMedia(mediaType string, mediaName string) ([]byte, error)
	GetMediaPath(mediaType string) (string, error)
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
}
