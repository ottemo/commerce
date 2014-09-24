package sqlite

import (
	"github.com/mxk/go-sqlite/sqlite3"
	"github.com/ottemo/foundation/db"
)

const (
	DEBUG_SQL = false

	FILTER_GROUP_STATIC  = "static"
	FILTER_GROUP_DEFAULT = "default"

	COLLECTION_NAME_COLUMN_INFO = "collection_column_info"
)

type T_DBFilterGroup struct {
	Name         string
	FilterValues []string
	ParentGroup  string
	OrSequence   bool
}

type SQLiteCollection struct {
	Connection *sqlite3.Conn
	TableName  string
	Columns    map[string]string

	ResultColumns []string
	FilterGroups  map[string]*T_DBFilterGroup
	Order         []string

	Limit string
}

type SQLite struct {
	Connection *sqlite3.Conn
}

var collections = map[string]db.I_DBCollection{}
