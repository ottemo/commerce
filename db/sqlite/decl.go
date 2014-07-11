package sqlite

import (
	"github.com/mxk/go-sqlite/sqlite3"
	"github.com/ottemo/foundation/db"
)

type SQLiteCollection struct {
	Connection *sqlite3.Conn
	TableName  string
	Columns    map[string]string

	Filters []string
	Order   []string

	Limit string
}

type SQLite struct {
	Connection *sqlite3.Conn
}

var collections = map[string]db.I_DBCollection{}
