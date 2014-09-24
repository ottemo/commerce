package mongo

import (
	"errors"
	"github.com/ottemo/foundation/db"
	"github.com/ottemo/foundation/env"

	"labix.org/v2/mgo"
)

// package self initializer
func init() {
	instance := new(MongoDB)

	env.RegisterOnConfigIniStart(instance.Startup)
	db.RegisterDBEngine(instance)
}

// mongo DB engine startup, opens connections to database
func (it *MongoDB) Startup() error {

	var DBUri = "mongodb://localhost:27017/ottemo"
	var DBName = "ottemo"

	if iniConfig := env.GetIniConfig(); iniConfig != nil {
		if iniValue := iniConfig.GetValue("mongodb.uri", DBUri); iniValue != "" {
			DBUri = iniValue
		}

		if iniValue := iniConfig.GetValue("mongodb.db", DBName); iniValue != "" {
			DBName = iniValue
		}
	}

	session, err := mgo.Dial(DBUri)
	if err != nil {
		return errors.New("Can't connect to MongoDB")
	}

	it.session = session
	it.database = session.DB(DBName)
	it.DBName = DBName
	it.collections = map[string]bool{}

	if MONGO_DEBUG {
		mgo.SetDebug(true)
		mgo.SetLogger(it)
	}

	if collectionsList, err := it.database.CollectionNames(); err == nil {
		for _, collection := range collectionsList {
			it.collections[collection] = true
		}
	}

	db.OnDatabaseStart()

	return nil
}

// debug logger mgo.log_Logger implementation
func (it *MongoDB) Output(calldepth int, s string) error {
	println(s)
	return nil
}
