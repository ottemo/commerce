package env

// IINIConfig is an initialization interface for reading INI file values
type IINIConfig interface {
	GetValue(Name string) string
	ListItems() []string
}

// IConfig is an interface for working with configuration entities and values
type IConfig interface {
	RegisterItem(Name string, Validator func(interface{}) (interface{}, bool), Default interface{}) error
	UnregisterItem(Name string) error

	GetValue(Name string) interface{}
	SetValue(Name string, Value interface{}) error

	ListItems() []string

	Load() error
	Save() error
}
