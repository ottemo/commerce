// Package sqlite is a "SQLite" implementation of interfaces declared in
// "github.com/ottemo/foundation/db" package
package sqlite

import (
	"regexp"
	"sync"

	"github.com/mxk/go-sqlite/sqlite3"
)

const (
	UUID_ID   = true  // flag which indicates to use UUID "_id" column type instead of default integer
	DEBUG_SQL = false // flag which indicates to perform log on each SQL operation

	FILTER_GROUP_STATIC  = "static"  // name for static filter, ref. to AddStaticFilter(...)
	FILTER_GROUP_DEFAULT = "default" // name for default filter, ref. to by AddFilter(...)

	COLLECTION_NAME_COLUMN_INFO = "collection_column_info" // table name to hold Ottemo types of columns
)

// regex expression used to check names used within SQL queries
var SQL_NAME_VALIDATOR = regexp.MustCompile("^[A-Za-z_][A-Za-z0-9_]*$")

// structure to hold information of named collection filter
type T_DBFilterGroup struct {
	Name         string
	FilterValues []string
	ParentGroup  string
	OrSequence   bool
}

// I_DBCollection implementer class
type SQLiteCollection struct {
	Name string

	ResultColumns []string
	FilterGroups  map[string]*T_DBFilterGroup
	Order         []string

	Limit string
}

// I_DBEngine implementer class
type SQLite struct {
	connection      *sqlite3.Conn
	connectionMutex sync.RWMutex

	attributeTypes      map[string]map[string]string
	attributeTypesMutex sync.RWMutex
}
