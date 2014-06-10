package database

type I_DBEngine interface {
	GetName() string

	CreateCollection(Name string) error
	GetCollection(Name string) (I_DBCollection, error)
	HasCollection(Name string) bool
}


type I_DBCollection interface {

	Load() ([]map[string]interface{}, error)
	LoadById(id string) (map[string]interface{}, error)

	Save(map[string]interface{}) (string, error)

	Delete() (int, error)
	DeleteById(id string) error

	AddFilter(ColumnName string, Operator string, Value string) error
	ClearFilters() error

	ListColumns() map[string]string
	HasColumn(ColumnName string) bool

	AddColumn(ColumnName string, ColumnType string, indexed bool) error
	RemoveColumn(ColumnName string) error
}
