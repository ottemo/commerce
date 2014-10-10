package sqlite

import (
	"strconv"

	"github.com/mxk/go-sqlite/sqlite3"
	"github.com/ottemo/foundation/db"
	"github.com/ottemo/foundation/env"
)

// returns current DB engine name
func (it *SQLite) GetName() string {
	return "Sqlite3"
}

// checks if collection(table) already exists
func (it *SQLite) HasCollection(collectionName string) bool {
	// collectionName = strings.ToLower(collectionName)

	SQL := "SELECT name FROM sqlite_master WHERE type='table' AND name='" + collectionName + "'"

	stmt, err := it.Connection.Query(SQL)
	defer closeStatement(stmt)

	if err == nil {
		return true
	} else {
		return false
	}
}

// creates cllection(table) by it's name
func (it *SQLite) CreateCollection(collectionName string) error {
	// collectionName = strings.ToLower(collectionName)

	SQL := "CREATE TABLE " + collectionName + " (_id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL)"
	if UUID_ID {
		SQL = "CREATE TABLE " + collectionName + " (_id NCHAR(24) PRIMARY KEY NOT NULL)"
	}

	if err := it.Connection.Exec(SQL); err == nil {
		return nil
	} else {
		return env.ErrorDispatch(err)
	}
}

// returns collection(table) by name or creates new one
func (it *SQLite) GetCollection(collectionName string) (db.I_DBCollection, error) {
	if !SQL_NAME_VALIDATOR.MatchString(collectionName) {
		return nil, env.ErrorNew("not valid collection name for DB engine")
	}

	if !it.HasCollection(collectionName) {
		if err := it.CreateCollection(collectionName); err != nil {
			return nil, env.ErrorDispatch(err)
		}
	}

	collection := &SQLiteCollection{
		Name:          collectionName,
		Connection:    it.Connection,
		FilterGroups:  make(map[string]*T_DBFilterGroup),
		Order:         make([]string, 0),
		ResultColumns: make([]string, 0),
	}

	return collection, nil
}

// returns collection(table) by name or creates new one
func (it *SQLite) RawQuery(query string) (map[string]interface{}, error) {

	result := make([]map[string]interface{}, 0, 10)

	row := make(sqlite3.RowMap)

	stmt, err := it.Connection.Query(query)
	defer closeStatement(stmt)

	if err == nil {
		return nil, env.ErrorDispatch(err)
	}

	for ; err == nil; err = stmt.Next() {
		if err := stmt.Scan(row); err == nil {

			if UUID_ID {
				if _, present := row["_id"]; present {
					row["_id"] = strconv.FormatInt(row["_id"].(int64), 10)
				}
			}

			result = append(result, row)
		} else {
			return result[0], nil
		}
	}

	return result[0], nil
}
