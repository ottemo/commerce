package mssql

import (
	"database/sql"
	"regexp"
	"sync"
	"time"

	"github.com/ottemo/commerce/env"
)

// Global constants for mssql package
const (
	ConstConnectionValidateInterval = time.Second * 10 // timer interval to ping connection and refresh it by perforce

	ConstUseUUIDids = false // flag which indicates to use UUID "_id" column type instead of default integer
	ConstDebugSQL   = true  // flag which indicates to perform log on each SQL operation
	ConstDebugFile  = "mssql.log"

	ConstFilterGroupStatic  = "static"  // name for static filter, ref. to AddStaticFilter(...)
	ConstFilterGroupDefault = "default" // name for default filter, ref. to by AddFilter(...)

	ConstCollectionNameColumnInfo = "collection_column_info" // table name to hold Ottemo types of columns

	ConstErrorModule = "db/mssql"
	ConstErrorLevel  = env.ConstErrorLevelService
)

// Global variables for mssql package
var (
	// dbEngine is an instance of database engine (one per application)
	dbEngine *DBEngine

	// ConstSQLNameValidator is a regex expression used to check names used within SQL queries
	ConstSQLNameValidator = regexp.MustCompile("^[A-Za-z_][A-Za-z0-9_]*$")
)

// RowMap - represents row of data from database
type RowMap map[string]interface{}

// StructDBFilterGroup is a structure to hold information of named collection filter
type StructDBFilterGroup struct {
	Name         string
	FilterValues []string
	ParentGroup  string
	OrSequence   bool
}

// DBCollection is of type InterfaceDBCollection
type DBCollection struct {
	Name string

	ResultColumns []string
	FilterGroups  map[string]*StructDBFilterGroup
	Order         []string

	Limit  int
	Offset int
}

// DBEngine is of type InterfaceDBEngine
type DBEngine struct {
	connection *sql.DB

	attributeTypes      map[string]map[string]string
	attributeTypesMutex sync.RWMutex

	isConnected bool
}

// connectionParamsType describes the parameters required to connect to DB
type connectionParamsType struct {
	uri             string
	dbName          string
	poolConnections int
	maxConnections  int
}
