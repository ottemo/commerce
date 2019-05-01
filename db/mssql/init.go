package mssql

import (
	"github.com/ottemo/commerce/db"
	"github.com/ottemo/commerce/env"
)

// init is the package bootstrap routine
func init() {
	dbEngine = new(DBEngine)
	dbEngine.attributeTypes = make(map[string]map[string]string)

	var _ db.InterfaceDBEngine = dbEngine

	var dbConnector = db.NewDBConnector(dbEngine)
	env.RegisterOnConfigIniStart(dbConnector.Connect)

	if err := db.RegisterDBEngine(dbEngine); err != nil {
		_ = env.ErrorDispatch(err)
	}
}
