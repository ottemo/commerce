package env

const (
	CONFIG_ITEM_GROUP_TYPE = "group"
)

type I_IniConfig interface {
	GetValue(Name string, Default string) string
	ListItems() []string
}

type I_Config interface {
	RegisterItem(Item T_ConfigItem, Validator F_ConfigValueValidator) error
	UnregisterItem(Path string) error

	ListPathes() []string
	GetValue(Path string) interface{}
	SetValue(Path string, Value interface{}) error

	GetGroupItems() []T_ConfigItem
	GetItemsInfo(Path string) []T_ConfigItem

	Load() error
}

type F_ConfigValueValidator func(interface{}) (interface{}, error)

type T_ConfigItem struct {
	Path  string
	Value interface{}

	Type string

	Editor  string
	Options interface{}

	Label       string
	Description string

	Image string
}
