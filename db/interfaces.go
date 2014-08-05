package db

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

	Count() (int, error)

	AddStaticFilter(ColumnName string, Operator string, Value interface{}) error

	AddFilter(ColumnName string, Operator string, Value interface{}) error
	ClearFilters() error

	AddSort(ColumnName string, Desc bool) error
	ClearSort() error

	SetResultColumns(columns ...string) error

	SetLimit(offset int, limit int) error

	ListColumns() map[string]string
	HasColumn(ColumnName string) bool

	AddColumn(ColumnName string, ColumnType string, indexed bool) error
	RemoveColumn(ColumnName string) error
}
