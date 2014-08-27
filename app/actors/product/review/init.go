package review

import (
	"errors"

	"github.com/ottemo/foundation/api"
	"github.com/ottemo/foundation/db"
)

// module entry point before app start
func init() {
	db.RegisterOnDatabaseStart(setupDB)
	api.RegisterOnRestServiceStart(setupAPI)
}

// DB preparations for current model implementation
func setupDB() error {

	dbEngine := db.GetDBEngine()
	if dbEngine == nil {
		return errors.New("Can't get database engine")
	}

	if collection, err := dbEngine.GetCollection("review"); err == nil {
		collection.AddColumn("product_id", "id", true)
		collection.AddColumn("visitor_id", "id", true)
		collection.AddColumn("username", "varchar(100)", true)
		collection.AddColumn("rating", "int", false)
		collection.AddColumn("review", "text", false)
		collection.AddColumn("created_at", "datetime", false)
	} else {
		return err
	}

	if collection, err := dbEngine.GetCollection("rating"); err == nil {
		collection.AddColumn("product_id", "id", true)
		collection.AddColumn("1star", "int", false)
		collection.AddColumn("2star", "int", false)
		collection.AddColumn("3star", "int", false)
		collection.AddColumn("4star", "int", false)
		collection.AddColumn("5star", "int", false)
	} else {
		return err
	}

	return nil
}
