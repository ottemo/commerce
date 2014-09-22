package sqlite

import (
	"strings"

	"github.com/mxk/go-sqlite/sqlite3"
	"github.com/ottemo/foundation/db"
	"strconv"
)

// returns current DB engine name
func (it *SQLite) GetName() string {
	return "Sqlite3"
}

// checks if collection(table) already exists
func (it *SQLite) HasCollection(CollectionName string) bool {
	CollectionName = strings.ToLower(CollectionName)

	SQL := "SELECT name FROM sqlite_master WHERE type='table' AND name='" + CollectionName + "'"

	stmt, err := it.Connection.Query(SQL)
	defer closeStatement(stmt)

	if err == nil {
		return true
	} else {
		return false
	}
}

// creates cllection(table) by it's name
func (it *SQLite) CreateCollection(CollectionName string) error {
	CollectionName = strings.ToLower(CollectionName)

	SQL := "CREATE TABLE " + CollectionName + "(_id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL)"
	if err := it.Connection.Exec(SQL); err == nil {
		return nil
	} else {
		return err
	}
}

// returns collection(table) by name or creates new one
func (it *SQLite) GetCollection(CollectionName string) (db.I_DBCollection, error) {

	if !it.HasCollection(CollectionName) {
		if err := it.CreateCollection(CollectionName); err != nil {
			return nil, err
		}
	}

	collection := &SQLiteCollection{
		TableName:     CollectionName,
		Connection:    it.Connection,
		Columns:       map[string]string{},
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
		return nil, err
	}

	for ; err == nil; err = stmt.Next() {
		if err := stmt.Scan(row); err == nil {
			if _, present := row["_id"]; present {
				row["_id"] = strconv.FormatInt(row["_id"].(int64), 10)
			}
			result = append(result, row)
		} else {
			return result[0], nil
		}
	}

	return result[0], nil
}
