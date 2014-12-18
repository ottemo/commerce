package mongo

import (
	"strings"

	"github.com/ottemo/foundation/db"
	"gopkg.in/mgo.v2"
)

// GetName returns the name of the db provider.
func (it *MongoDB) GetName() string { return "MongoDB" }

// HasCollection returns true/false if the given Collection exists.
func (it *MongoDB) HasCollection(CollectionName string) bool {
	CollectionName = strings.ToLower(CollectionName)

	return true
}

// CreateCollection will create a MongoDB Collection with the provided name.
func (it *MongoDB) CreateCollection(CollectionName string) error {
	CollectionName = strings.ToLower(CollectionName)

	err := it.database.C(CollectionName).Create(new(mgo.CollectionInfo))
	//it.database.C(CollectionName).EnsureIndex(mgo.Index{Key: []string{"_id"}, Unique: true})

	//CMD := "db.createCollection(\"" + CollectionName + "\")"
	//println(CMD)
	//err := it.database.Run(CMD, nil)

	return err
}

// GetCollection will return an Interface to the provided Collection name.
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
		Selector:   map[string]interface{}{},
		Name:       CollectionName,
		database:   it.database,
		collection: mgoCollection}

	return result, nil
}
