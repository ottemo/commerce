package db

import (
	"errors"
)

// returns database collection or error otherwise
func GetCollection(CollectionName string) (I_DBCollection, error) {
	dbEngine := GetDBEngine()
	if dbEngine == nil {
		return nil, errors.New("Can't get DBEngine")
	}

	return dbEngine.GetCollection(CollectionName)
}
