package sqlite

import (
	"github.com/mxk/go-sqlite/sqlite3"
	"github.com/ottemo/foundation/db"
	"github.com/ottemo/foundation/env"
)

func init() {
	instance := new(SQLite)

	env.RegisterOnConfigIniStart(instance.Startup)
	db.RegisterDBEngine(instance)
}

func (it *SQLite) Startup() error {

	// opening connection
	var uri string = "ottemo.db"

	if iniConfig := env.GetIniConfig(); iniConfig != nil {
		if iniValue := iniConfig.GetValue("db.sqlite3.uri", uri); iniValue != "" {
			uri = iniValue
		}
	}

	if newConnection, err := sqlite3.Open(uri); err == nil {
		it.Connection = newConnection
	} else {
		return env.ErrorDispatch(err)
	}

	// making column info table
	SQL := "CREATE TABLE IF NOT EXISTS " + COLLECTION_NAME_COLUMN_INFO + ` (
		_id        INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
		collection VARCHAR(255),
		column     VARCHAR(255),
		type       VARCHAR(255),
		indexed    NUMERIC)`

	err := it.Connection.Exec(SQL)
	if err != nil {
		return sqlError(SQL, err)
	}

	db.OnDatabaseStart()

	return nil
}
