package config

import (
	"github.com/ottemo/foundation/api"
	"github.com/ottemo/foundation/db"
	"github.com/ottemo/foundation/env"
)

const (
	// database collection name to store config values
	CONFIG_COLLECTION_NAME = "config"
)

// init makes package self-initialization routine
func init() {
	instance := &DefaultConfig{
		configValues:     make(map[string]interface{}),
		configTypes:      make(map[string]string),
		configValidators: make(map[string]env.F_ConfigValueValidator)}

	db.RegisterOnDatabaseStart(setupDB)
	db.RegisterOnDatabaseStart(instance.Load)

	api.RegisterOnRestServiceStart(setupAPI)

	env.RegisterConfig(instance)
}

// setupDB prepares system database for package usage
func setupDB() error {
	collection, err := db.GetCollection(CONFIG_COLLECTION_NAME)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	collection.AddColumn("path", "varchar(255)", true)
	collection.AddColumn("value", "text", false)

	collection.AddColumn("type", "text", false)

	collection.AddColumn("editor", "text", false)
	collection.AddColumn("options", "text", false)

	collection.AddColumn("label", "varchar(255)", false)
	collection.AddColumn("description", "text", false)

	collection.AddColumn("image", "varchar(255)", false)

	return nil
}
