package models

type I_Model interface {
	GetModelName() string
	GetImplementationName() string
	New() (I_Model, error)
}

type I_Storable interface {
	Save() error
	Load(id string) error
}

type I_Object interface {
	GetId() string

	Has(Attribute string) bool
	Get(Attribute string) interface{}
	Set(Attribute string, Value interface{}) error

	ListAttributes() []T_AttributeInfo
}

type I_Mapable interface {
	FromHashMap(HashMap map[string]interface{}) error
	ToHashMap() map[string]interface{}
}

type I_CustomAttributes interface {
	AddNewAttribute( newAttribute T_AttributeInfo ) error
	RemoveAttribute( attributeName string ) error
}
