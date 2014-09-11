package attributes

import (
	"github.com/ottemo/foundation/db"
)

// package self initializer
func init() {
	db.RegisterOnDatabaseStart(SetupDB)
}

// setups database tables for model usage
func SetupDB() error {

	if collection, err := db.GetCollection("custom_attributes"); err == nil {
		collection.AddColumn("model", "text", true)
		collection.AddColumn("collection", "text", true)
		collection.AddColumn("attribute", "text", true)
		collection.AddColumn("type", "text", false)
		collection.AddColumn("required", "bool", false)
		collection.AddColumn("label", "text", true)
		collection.AddColumn("group", "text", false)
		collection.AddColumn("editors", "text", false)
		collection.AddColumn("options", "text", false)
		collection.AddColumn("default", "text", false)
		collection.AddColumn("validators", "text", false)
		collection.AddColumn("layered", "bool", false)

	} else {
		return err
	}

	return nil
}
