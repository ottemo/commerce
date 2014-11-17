// Package mongo is a "MongoDB" implementation of interfaces declared in
// "github.com/ottemo/foundation/db" package
package mongo

import (
	"sync"

	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
)

var (
	attributeTypes      = make(map[string]map[string]string) // cached values of collection attribute types
	attributeTypesMutex sync.RWMutex                         // syncronization for attributeTypes modification
)

const (
	MONGO_DEBUG = false // flag which indicates to perform log on each operation

	FILTER_GROUP_STATIC  = "static"  // name for static filter, ref. to AddStaticFilter(...)
	FILTER_GROUP_DEFAULT = "default" // name for default filter, ref. to by AddFilter(...)

	COLLECTION_NAME_COLUMN_INFO = "collection_column_info" // collection name to hold Ottemo types of attributes
)

// structure to hold information of named collection filter
type T_DBFilterGroup struct {
	Name         string
	FilterValues []bson.D
	ParentGroup  string
	OrSequence   bool
}

// I_DBCollection implementer class
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

// I_DBEngine implementer class
type MongoDB struct {
	database *mgo.Database
	session  *mgo.Session

	DBName      string
	collections map[string]bool
}
