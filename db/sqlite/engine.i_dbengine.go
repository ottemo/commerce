package sqlite

import (
	"strconv"

	"github.com/mxk/go-sqlite/sqlite3"
	"github.com/ottemo/foundation/db"
	"github.com/ottemo/foundation/env"
)

// returns current DB engine name
func (it *DBEngine) GetName() string {
	return "Sqlite3"
}

// checks if collection(table) already exists
func (it *DBEngine) HasCollection(collectionName string) bool {
	// collectionName = strings.ToLower(collectionName)

	SQL := "SELECT name FROM sqlite_master WHERE type='table' AND name='" + collectionName + "'"

	stmt, err := connectionQuery(SQL)
	defer closeStatement(stmt)

	if err == nil {
		return true
	} else {
		return false
	}
}

// creates cllection(table) by it's name
func (it *DBEngine) CreateCollection(collectionName string) error {
	// collectionName = strings.ToLower(collectionName)

	SQL := "CREATE TABLE " + collectionName + " (_id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL)"
	if ConstUseUUIDids {
		SQL = "CREATE TABLE " + collectionName + " (_id NCHAR(24) PRIMARY KEY NOT NULL)"
	}

	if err := connectionExec(SQL); err == nil {
		return nil
	} else {
		return env.ErrorDispatch(err)
	}
}

// returns collection(table) by name or creates new one
func (it *DBEngine) GetCollection(collectionName string) (db.InterfaceDBCollection, error) {
	if !ConstSQLNameValidator.MatchString(collectionName) {
		return nil, env.ErrorNew("not valid collection name for DB engine")
	}

	if !it.HasCollection(collectionName) {
		if err := it.CreateCollection(collectionName); err != nil {
			return nil, env.ErrorDispatch(err)
		}
	}

	collection := &DBCollection{
		Name:          collectionName,
		FilterGroups:  make(map[string]*StructDBFilterGroup),
		Order:         make([]string, 0),
		ResultColumns: make([]string, 0),
	}

	if _, present := it.attributeTypes[collectionName]; !present {
		collection.ListColumns()
	}

	return collection, nil
}

// returns collection(table) by name or creates new one
func (it *DBEngine) RawQuery(query string) (map[string]interface{}, error) {

	result := make([]map[string]interface{}, 0, 10)

	row := make(sqlite3.RowMap)

	stmt, err := connectionQuery(query)
	defer closeStatement(stmt)

	if err == nil {
		return nil, env.ErrorDispatch(err)
	}

	for ; err == nil; err = stmt.Next() {
		if err := stmt.Scan(row); err == nil {

			if ConstUseUUIDids {
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
