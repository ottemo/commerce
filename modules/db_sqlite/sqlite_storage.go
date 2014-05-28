package db_sqlite

import (
	"errors"
	"github.com/ottemo/platform/interfaces/database"
	"github.com/ottemo/platform/interfaces/config"

	"github.com/ottemo/platform/tools/module_manager"

	"code.google.com/p/go-sqlite/go1/sqlite3"
)

func init() {
	module_manager.RegisterModule( new(SQLiteEngine) )
}

const(
	configURI = "sqlite.uri"
)

type SQLiteEngine struct {
	connection *sqlite3.Conn
}

// I_Module interface implementation
//----------------------------------
func (it *SQLiteEngine) GetModuleName() string { return "Sqlite3" }
func (it *SQLiteEngine) GetModuleDepends() []string { return make([]string, 0) }

func (it *SQLiteEngine) ModuleMakeSysInit() error { return nil }

func (it *SQLiteEngine) ModuleMakeConfig() error {
	cfg := config.GetConfig()
	return cfg.RegisterItem(configURI, "string", config.ConfigEmptyValueValidator, "ottemo.db" )
}

func (it *SQLiteEngine) ModuleMakeInit() error {
	return database.RegisterDatabaseEngine("Sqlite3", it)
}

func (it *SQLiteEngine) ModuleMakeVerify() error {
	if uri, ok := config.GetConfig().GetValue(configURI).(string); ok {
		if connection, err := sqlite3.Open(uri); err == nil {
			it.connection = connection
		} else {
			return errors.New("can't open sqlite connection to '" + uri + "': " + err.Error() )
		}
	} else {
		return errors.New("can't get config value '" + configURI + "'")
	}
	return nil
}

func (it *SQLiteEngine) ModuleMakeLoad() error { return nil }
func (it *SQLiteEngine) ModuleMakeInstall() error { return nil }
func (it *SQLiteEngine) ModuleMakePostInstall() error { return nil }



// I_DBStorage interface implementation
//----------------------------------

func (it *SQLiteEngine) GetCollection(Name string) (database.I_Collection, error) {
	return new(SQLiteCollection), nil
}

func (it *SQLiteEngine) GetCollectionFor(Object database.I_DBObject) (database.I_Collection, error) {
	return new(SQLiteCollection), nil
}
