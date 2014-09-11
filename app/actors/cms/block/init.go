package block

import (
	"github.com/ottemo/foundation/api"
	"github.com/ottemo/foundation/db"

	"github.com/ottemo/foundation/app/models"
	"github.com/ottemo/foundation/app/models/cms"
)

// module entry point before app start
func init() {
	cmsBlockInstance := new(DefaultCMSBlock)
	var _ cms.I_CMSBlock = cmsBlockInstance
	models.RegisterModel(cms.MODEL_NAME_CMS_BLOCK, cmsBlockInstance)

	cmsBlockCollectionInstance := new(DefaultCMSBlockCollection)
	var _ cms.I_CMSBlockCollection = cmsBlockCollectionInstance
	models.RegisterModel(cms.MODEL_NAME_CMS_BLOCK_COLLECTION, cmsBlockCollectionInstance)

	db.RegisterOnDatabaseStart(setupDB)
	api.RegisterOnRestServiceStart(setupAPI)
}

// DB preparations for current model implementation
func setupDB() error {
	collection, err := db.GetCollection(CMS_BLOCK_COLLECTION_NAME)
	if err != nil {
		return err
	}

	collection.AddColumn("identifier", "varchar(255)", true)
	collection.AddColumn("content", "text", false)
	collection.AddColumn("created_at", "datetime", false)
	collection.AddColumn("updated_at", "datetime", false)

	return nil
}
