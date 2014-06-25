package database

type IDBEngine interface {
	GetName() string

	CreateCollection(Name string) error
	GetCollection(Name string) (IDBCollection, error)
	HasCollection(Name string) bool
}

type IDBCollection interface {
	Load() ([]map[string]interface{}, error)
	LoadById(id string) (map[string]interface{}, error)

	Save(map[string]interface{}) (string, error)

	Delete() (int, error)
	DeleteById(id string) error

	AddFilter(ColumnName string, Operator string, Value string) error //TODO: modify (Value string) to (Value interface{})
	ClearFilters() error

	AddSort(ColumnName string, Desc bool) error
	ClearSort() error

	SetLimit(offset int, limit int) error

	ListColumns() map[string]string
	HasColumn(ColumnName string) bool

	AddColumn(ColumnName string, ColumnType string, indexed bool) error
	RemoveColumn(ColumnName string) error
}
