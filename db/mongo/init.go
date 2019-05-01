package mongo

import (
	"github.com/ottemo/commerce/db"
	"github.com/ottemo/commerce/env"
)

// init makes package self-initialization routine
func init() {
	instance := new(DBEngine)

	var dbConnection = db.NewConnection(instance)
	env.RegisterOnConfigIniStart(dbConnection.Async)

	_ = db.RegisterDBEngine(instance)
}

// Output is a implementation of mgo.log_Logger interface
func (it *DBEngine) Output(calldepth int, s string) error {
	env.Log("mongo.log", "DEBUG", s)
	return nil
}
