package sqlite

import (
	"strings"

	"github.com/ottemo/foundation/db"
)

func (it *SQLite) GetName() string { return "Sqlite3" }

func (it *SQLite) HasCollection(CollectionName string) bool {
	CollectionName = strings.ToLower(CollectionName)

	SQL := "SELECT name FROM sqlite_master WHERE type='table' AND name='" + CollectionName + "'"
	if _, err := it.Connection.Query(SQL); err == nil {
		return true
	} else {
		return false
	}
}

func (it *SQLite) CreateCollection(CollectionName string) error {
	CollectionName = strings.ToLower(CollectionName)

	SQL := "CREATE TABLE " + CollectionName + "(_id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL)"
	if err := it.Connection.Exec(SQL); err == nil {
		return nil
	} else {
		return err
	}
}

func (it *SQLite) GetCollection(CollectionName string) (db.I_DBCollection, error) {
	CollectionName = strings.ToLower(CollectionName)

	if collection, present := collections[CollectionName]; present {
		return collection, nil
	} else {
		if !it.HasCollection(CollectionName) {
			if err := it.CreateCollection(CollectionName); err != nil {
				return nil, err
			}
		}

		collection := &SQLiteCollection{TableName: CollectionName, Connection: it.Connection, Columns: map[string]string{}}
		collections[CollectionName] = collection

		return collection, nil
	}
}
