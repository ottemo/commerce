package page

import (
	"github.com/ottemo/foundation/api"
	"github.com/ottemo/foundation/db"

	"github.com/ottemo/foundation/app/models"
	"github.com/ottemo/foundation/app/models/cms"
)

// module entry point before app start
func init() {
	pageInstance := new(DefaultCMSPage)

	//checking interface implementation
	(func(cms.I_CMSPage) {})(pageInstance)

	models.RegisterModel(cms.CMS_PAGE_MODEL_NAME, pageInstance)

	db.RegisterOnDatabaseStart(setupDB)
	api.RegisterOnRestServiceStart(setupAPI)
}

// DB preparations for current model implementation
func setupDB() error {
	collection, err := db.GetCollection(CMS_PAGE_COLLECTION_NAME)
	if err != nil {
		return err
	}

	collection.AddColumn("identifier", "varchar(255)", true)
	collection.AddColumn("url", "varchar(255)", true)
	collection.AddColumn("title", "varchar(255)", false)
	collection.AddColumn("content", "text", false)
	collection.AddColumn("meta_keywords", "varchar(255)", false)
	collection.AddColumn("meta_description", "text", false)
	collection.AddColumn("created_at", "datetime", false)
	collection.AddColumn("updated_at", "datetime", false)

	return nil
}
