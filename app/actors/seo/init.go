package seo

import (
	"github.com/ottemo/foundation/api"
	"github.com/ottemo/foundation/db"
	"github.com/ottemo/foundation/env"
)

// module entry point before app start
func init() {
	api.RegisterOnRestServiceStart(setupAPI)
	db.RegisterOnDatabaseStart(setupDB)
}

// DB preparations for current model implementation
func setupDB() error {

	if collection, err := db.GetCollection(COLLECTION_NAME_URL_REWRITES); err == nil {
		collection.AddColumn("url", "varchar(255)", true)
		collection.AddColumn("type", "varchar(255)", true)
		collection.AddColumn("rewrite", "varchar(255)", false)
		collection.AddColumn("title", "varchar(255)", false)
		collection.AddColumn("meta_keywords", "varchar(255)", false)
		collection.AddColumn("meta_description", "varchar(255)", false)
	} else {
		return env.ErrorDispatch(err)
	}

	return nil
}
