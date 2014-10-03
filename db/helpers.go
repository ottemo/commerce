package db

import (
	"github.com/ottemo/foundation/env"
)

// returns database collection or error otherwise
func GetCollection(CollectionName string) (I_DBCollection, error) {
	dbEngine := GetDBEngine()
	if dbEngine == nil {
		return nil, env.ErrorNew("Can't get DBEngine")
	}

	return dbEngine.GetCollection(CollectionName)
}
