package sqlite

import (
	"github.com/mxk/go-sqlite/sqlite3"
	"github.com/ottemo/foundation/db"
)

const (
	DEBUG_SQL = false
)

type SQLiteCollection struct {
	Connection *sqlite3.Conn
	TableName  string
	Columns    map[string]string

	ResultColumns []string
	StaticFilters map[string]string

	Filters map[string]string
	Order   []string

	Limit string
}

type SQLite struct {
	Connection *sqlite3.Conn
}

var collections = map[string]db.I_DBCollection{}
