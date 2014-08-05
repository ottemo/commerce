package mongo

import (
	"github.com/ottemo/foundation/db"
	"labix.org/v2/mgo"
	"strings"
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

	return err
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
		Selector:         make(map[string]interface{}),
		StaticSelector:   make(map[string]interface{}),
		ResultAttributes: make([]string, 0),
		database:         it.database,
		collection:       mgoCollection,
	}

	return result, nil
}
