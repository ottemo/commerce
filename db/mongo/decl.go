package mongo

import (
	"gopkg.in/mgo.v2"
)

// MongoDBCollection  holds the meta data related to all collections.
type MongoDBCollection struct {
	database   *mgo.Database
	collection *mgo.Collection
	Name       string

	Selector map[string]interface{}
	Sort     []string

	Limit  int
	Offset int
}

// MongoDB holds the connection and session information.
type MongoDB struct {
	database *mgo.Database
	session  *mgo.Session

	DBName      string
	collections map[string]bool
}
