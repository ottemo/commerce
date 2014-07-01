package env

type I_IniConfig interface {
	GetValue(Name string) string
	ListItems() []string
}

type I_Config interface {
	RegisterItem(Name string, Validator func(interface{}) (interface{}, bool), Default interface{} ) error
	UnregisterItem(Name string) error

	GetValue(Name string) interface{}
	SetValue(Name string, Value interface{}) error

	ListItems() []string

	Load() error
	Save() error
}
