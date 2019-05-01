package test

import (
	"testing"

	"githbu.com/ottemo/commerce/db"
	"github.com/ottemo/commerce/utils"
)

// function used to test most product model operations
func TestDatabaseOperations(t *testing.T) {
	err := StartAppInTestingMode()
	if err != nil {
		t.Error(err)
	}

	dbEngine := db.GetDBEngine()
	if dbEngine == nil {
		t.Error("Can't obtain database engine")
	}
	if !dbEngine.HasCollection("test") {
		dbEngine.CreateCollection("test")
	}

	dbCollection, err := dbEngine.GetCollection("test")
	if err != nil {
		t.Error(err)
	}

	data := map[string]interface{}{
		"type_" + db.ConstTypeBoolean:  true,
		"type_" + db.ConstTypeVarchar:  "varchar",
		"type_" + db.ConstTypeText:     "text",
		"type_" + db.ConstTypeInteger:  1,
		"type_" + db.ConstTypeDecimal:  1.5,
		"type_" + db.ConstTypeMoney:    10.33,
		"type_" + db.ConstTypeFloat:    0.667341,
		"type_" + db.ConstTypeDatetime: utils.Time(),
		"type_" + db.ConstTypeJSON:     map[string]interface{}{"a": 10, "b": 25},
	}

	for column, _ := range data {
		dbType := column[5:]
		if !dbCollection.HasColumn(column) {
			dbCollection.AddColumn(column, dbType, false)
		}
	}

	id, err := dbCollection.Save(data)
	if err != nil {
		t.Error(err)
	}

	saved, err := dbCollection.LoadByID(id)
	if err != nil {
		t.Error(err)
	}

	if !utils.MatchMapAValuesToMapB(data, saved) {
		t.Error("Saved values does not match to originals")
	}

	err = dbCollection.DeleteByID(id)
	if err != nil {
		t.Error(err)
	}

	_, err = dbCollection.Delete()
	if err != nil {
		t.Error(err)
	}
}
