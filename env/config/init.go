package config

import (
	"errors"

	"github.com/ottemo/foundation/api"
	"github.com/ottemo/foundation/db"
	"github.com/ottemo/foundation/env"
)

const (
	CONFIG_COLLECTION_NAME = "config"
)

func init() {
	instance := &DefaultConfig{
		configValues:     make(map[string]interface{}),
		configValidators: make(map[string]env.F_ConfigValueValidator)}

	db.RegisterOnDatabaseStart(setupDB)
	db.RegisterOnDatabaseStart(instance.Load)

	api.RegisterOnRestServiceStart(setupAPI)

	env.RegisterConfig(instance)
}

// DB preparations for current model implementation
func setupDB() error {
	if dbEngine := db.GetDBEngine(); dbEngine != nil {
		collection, err := dbEngine.GetCollection(CONFIG_COLLECTION_NAME)
		if err != nil {
			return err
		}

		collection.AddColumn("path", "varchar(255)", true)
		collection.AddColumn("value", "text", false)

		collection.AddColumn("type", "text", false)

		collection.AddColumn("editor", "text", false)
		collection.AddColumn("options", "text", false)

		collection.AddColumn("label", "varchar(255)", false)
		collection.AddColumn("description", "text", false)

		collection.AddColumn("image", "varchar(255)", false)
	} else {
		return errors.New("Can't get database engine")
	}

	return nil
}
