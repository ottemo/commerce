package page

import (
	"github.com/ottemo/foundation/api"
	"github.com/ottemo/foundation/db"

	"github.com/ottemo/foundation/app/models"
	"github.com/ottemo/foundation/app/models/cms"
)

// module entry point before app start
func init() {
	cmsPageInstance := new(DefaultCMSPage)
	var _ cms.I_CMSPage = cmsPageInstance
	models.RegisterModel(cms.MODEL_NAME_CMS_PAGE, cmsPageInstance)

	cmsPageCollectionInstance := new(DefaultCMSPageCollection)
	var _ cms.I_CMSPageCollection = cmsPageCollectionInstance
	models.RegisterModel(cms.MODEL_NAME_CMS_PAGE_COLLECTION, cmsPageCollectionInstance)

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
