package models

var (
	declaredModels = map[string]I_Model{}
)

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
