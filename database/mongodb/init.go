package mongodb

import (
	"errors"
	"github.com/ottemo/foundation/database"
	"github.com/ottemo/foundation/env"
	"labix.org/v2/mgo"
)

func init() {
	instance := new(MongoDB)

	env.RegisterOnConfigIniStart( instance.Startup )
	database.RegisterDBEngine( instance )
}


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

	it.session  = session
	it.database = session.DB(DBName)
	it.DBName = DBName
	it.collections = map[string]bool{}

	if collectionsList, err := it.database.CollectionNames(); err == nil {
		for _, collection := range collectionsList {
			it.collections[collection] = true
		}
	}

	database.OnDatabaseStart()

	return nil
}
