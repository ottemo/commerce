package mongo

import (
	"sync"

	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
)

var (
	attributeTypes      = make(map[string]map[string]string)
	attributeTypesMutex sync.RWMutex
)

const (
	MONGO_DEBUG = false

	FILTER_GROUP_STATIC  = "static"
	FILTER_GROUP_DEFAULT = "default"

	COLLECTION_NAME_COLUMN_INFO = "collection_column_info"
)

type T_DBFilterGroup struct {
	Name         string
	FilterValues []bson.D
	ParentGroup  string
	OrSequence   bool
}

type MongoDBCollection struct {
	database   *mgo.Database
	collection *mgo.Collection

	subcollections []*MongoDBCollection
	subresults     []*bson.Raw

	Name string

	FilterGroups map[string]*T_DBFilterGroup

	Sort []string

	ResultAttributes []string

	Limit  int
	Offset int
}

type MongoDB struct {
	database *mgo.Database
	session  *mgo.Session

	DBName      string
	collections map[string]bool
}
