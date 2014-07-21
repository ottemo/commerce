package category

import (
	"errors"
	"github.com/ottemo/foundation/app/models"
	"github.com/ottemo/foundation/db"

	"github.com/ottemo/foundation/api"

	"github.com/ottemo/foundation/app/models/category"
)

func init() {
	instance := new(DefaultCategory)

	ifce := interface{} (instance)
	if _, ok := ifce.(models.I_Model); !ok {
		panic("DefaultCategory - I_Model interface not implemented")
	}
	if _, ok := ifce.(models.I_Object); !ok {
		panic("DefaultCategory - I_Object interface not implemented")
	}
	if _, ok := ifce.(models.I_Storable); !ok {
		panic("DefaultCategory - I_Storable interface not implemented")
	}
	if _, ok := ifce.(models.I_Listable); !ok {
		panic("DefaultCategory - I_Listable interface not implemented")
	}
	if _, ok := ifce.(category.I_Category); !ok {
		panic("DefaultCategory - I_Category interface not implemented")
	}

	models.RegisterModel("Category", instance)
	db.RegisterOnDatabaseStart(instance.setupModel)

	api.RegisterOnRestServiceStart(instance.setupAPI)
}

func (it *DefaultCategory) setupModel() error {

	if dbEngine := db.GetDBEngine(); dbEngine != nil {
		collection, err := dbEngine.GetCollection(CATEGORY_COLLECTION_NAME)
		if err != nil {
			return err
		}

		collection.AddColumn("parent_id", "id", true)
		collection.AddColumn("path", "text", true)
		collection.AddColumn("name", "text", true)

		collection, err = dbEngine.GetCollection(CATEGORY_PRODUCT_JUNCTION_COLLECTION_NAME)
		if err != nil {
			return err
		}

		collection.AddColumn("category_id", "id", true)
		collection.AddColumn("product_id", "id", true)

	} else {
		return errors.New("Can't get database engine")
	}

	return nil
}
