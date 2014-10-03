package mongo

import (
	"strings"

	"github.com/ottemo/foundation/db"
	"github.com/ottemo/foundation/env"

	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
)

// returns current DB engine name
func (it *MongoDB) GetName() string {
	return "MongoDB"
}

// checks if collection already exists
func (it *MongoDB) HasCollection(CollectionName string) bool {
	CollectionName = strings.ToLower(CollectionName)

	return true
}

// creates cllection by it's name
func (it *MongoDB) CreateCollection(CollectionName string) error {
	CollectionName = strings.ToLower(CollectionName)

	err := it.database.C(CollectionName).Create(new(mgo.CollectionInfo))
	//it.database.C(CollectionName).EnsureIndex(mgo.Index{Key: []string{"_id"}, Unique: true})

	//CMD := "db.createCollection(\"" + CollectionName + "\")"
	//println(CMD)
	//err := it.database.Run(CMD, nil)

	return env.ErrorDispatch(err)
}

// returns collection by name or creates new one
func (it *MongoDB) GetCollection(CollectionName string) (db.I_DBCollection, error) {
	CollectionName = strings.ToLower(CollectionName)

	if _, present := it.collections[CollectionName]; !present {
		if err := it.CreateCollection(CollectionName); err != nil {
			return nil, err
		}
		it.collections[CollectionName] = true
	}

	mgoCollection := it.database.C(CollectionName)

	result := &MongoDBCollection{
		Name:             CollectionName,
		FilterGroups:     make(map[string]*T_DBFilterGroup),
		ResultAttributes: make([]string, 0),
		database:         it.database,
		collection:       mgoCollection,
	}

	return result, nil
}

// executes raw query for DB engine.
//   This function makes eval commang on mongo db (http://docs.mongodb.org/manual/reference/command/eval/#dbcmd.eval)
//   so if you are using "db.collection.find()" - it returns cursor object, do not forget to add ".toArray()", i.e.
//   "db.collection.find().toArray()"
func (it *MongoDB) RawQuery(query string) (map[string]interface{}, error) {
	result := make(map[string]interface{})
	err := it.database.Run(bson.D{{"eval", query}}, result)
	return result, err
}
