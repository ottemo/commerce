package attributes

import (
	"errors"
	"github.com/ottemo/foundation/db"
)

func init(){
	db.RegisterOnDatabaseStart( SetupModel )
}

func SetupModel() error {

	if dbEngine := db.GetDBEngine(); dbEngine != nil {
		if collection, err := dbEngine.GetCollection("custom_attributes"); err == nil {
			collection.AddColumn("model", "text", true)
			collection.AddColumn("collection", "text", true)
			collection.AddColumn("attribute", "text", true)
			collection.AddColumn("type", "text", false)
			collection.AddColumn("label", "text", true)
			collection.AddColumn("group", "text", false)
			collection.AddColumn("editors", "text", false)
			collection.AddColumn("options", "text", false)
			collection.AddColumn("default", "text", false)

		} else {
			return err
		}
	} else {
		return errors.New("Can't get database engine")
	}

	return nil
}
