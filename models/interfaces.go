package models

type IModel interface {
	GetModelName() string
	GetImplementationName() string
	New() (IModel, error)
}

type IStorable interface {
	GetId() string
	SetId(string) error

	Save() error
	Load(id string) error
	Delete(Id string) error
}

type IObject interface {
	Get(Attribute string) interface{}
	Set(Attribute string, Value interface{}) error

	GetAttributesInfo() []T_AttributeInfo
}

type IMapable interface {
	FromHashMap(HashMap map[string]interface{}) error
	ToHashMap() map[string]interface{}
}

type ICustomAttributes interface {
	AddNewAttribute( newAttribute T_AttributeInfo ) error
	RemoveAttribute( attributeName string ) error
}
