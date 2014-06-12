package mongodb

import (
	"strings"
	"labix.org/v2/mgo"
	"github.com/ottemo/foundation/database"
)

func (it *MongoDB) GetName() string { return  "MongoDB" }


func (it *MongoDB) HasCollection(CollectionName string) bool {
	CollectionName = strings.ToLower(CollectionName)

	return true
}

func (it *MongoDB) CreateCollection(CollectionName string) error {
	CollectionName = strings.ToLower(CollectionName)

	err := it.database.C(CollectionName).Create( new(mgo.CollectionInfo) )
	//it.database.C(CollectionName).EnsureIndex(mgo.Index{Key: []string{"_id"}, Unique: true})

	//CMD := "db.createCollection(\"" + CollectionName + "\")"
	//println(CMD)
	//err := it.database.Run(CMD, nil)

	return err
}

func (it *MongoDB) GetCollection(CollectionName string) (database.I_DBCollection, error) {
	CollectionName = strings.ToLower(CollectionName)

	if _, present := it.collections[CollectionName]; !present {
		if err := it.CreateCollection(CollectionName); err != nil {
			return nil, err
		}
		it.collections[CollectionName] = true
	}

	mgoCollection := it.database.C(CollectionName)

	result := &MongoDBCollection {
		Selector: map[string]interface{} {},
		Name: CollectionName,
		database: it.database,
		collection: mgoCollection }

	return result, nil
}
