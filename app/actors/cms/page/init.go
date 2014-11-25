package page

import (
	"github.com/ottemo/foundation/api"
	"github.com/ottemo/foundation/app/models"
	"github.com/ottemo/foundation/app/models/cms"
	"github.com/ottemo/foundation/db"
	"github.com/ottemo/foundation/env"
)

// init makes package self-initialization routine
func init() {
	cmsPageInstance := new(DefaultCMSPage)
	var _ cms.InterfaceCMSPage = cmsPageInstance
	models.RegisterModel(cms.ConstModelNameCMSPage, cmsPageInstance)

	cmsPageCollectionInstance := new(DefaultCMSPageCollection)
	var _ cms.InterfaceCMSPageCollection = cmsPageCollectionInstance
	models.RegisterModel(cms.ConstModelNameCMSPageCollection, cmsPageCollectionInstance)

	db.RegisterOnDatabaseStart(setupDB)
	api.RegisterOnRestServiceStart(setupAPI)
}

// setupDB prepares system database for package usage
func setupDB() error {
	collection, err := db.GetCollection(ConstCmsPageCollectionName)
	if err != nil {
		return env.ErrorDispatch(err)
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
