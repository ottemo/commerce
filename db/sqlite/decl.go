package sqlite

import (
	"regexp"
	"sync"

	"github.com/mxk/go-sqlite/sqlite3"
)

const (
	DEBUG_SQL = false

	FILTER_GROUP_STATIC  = "static"
	FILTER_GROUP_DEFAULT = "default"

	COLLECTION_NAME_COLUMN_INFO = "collection_column_info"
)

var (
	attributeTypes      = make(map[string]map[string]string)
	attributeTypesMutex sync.RWMutex

	SQL_NAME_VALIDATOR = regexp.MustCompile("^[A-Za-z_][A-Za-z0-9_]*$")
)

type T_DBFilterGroup struct {
	Name         string
	FilterValues []string
	ParentGroup  string
	OrSequence   bool
}

type SQLiteCollection struct {
	Connection *sqlite3.Conn

	Name string

	ResultColumns []string
	FilterGroups  map[string]*T_DBFilterGroup
	Order         []string

	Limit string
}

type SQLite struct {
	Connection *sqlite3.Conn
}
