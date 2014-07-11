package test

import (
	"testing"

	"github.com/ottemo/foundation/db/sqlite"
)

func TestCollection(t *testing.T) {
	sqliteEngine := new(sqlite.SQLite)

	sqliteEngine.Startup()

	if err := sqliteEngine.CreateCollection("test"); err == nil {
		collection, _ := sqliteEngine.GetCollection("test")
		if err := collection.AddColumn("name", "varchar", false); err != nil {
			t.Error(err)
		}
		if err := collection.AddColumn("value", "int", false); err != nil {
			t.Error(err)
		}

		if err := collection.AddColumn("deleteme", "text", false); err != nil {
			t.Error(err)
		}
		// TODO: implement and debug RemoveColumn
		// if err := collection.RemoveColumn("deleteme"); err != nil {
		// 	t.Error(err)
		// }

		x := map[string]interface{}{"name": "value_10", "value": 10}
		if _, err := collection.Save(x); err != nil {
			t.Fatal(err)
		}

		t.Logf("new ID: %i", x["_id"])

		if all, err := collection.Load(); err == nil {
			for _, x := range all {
				t.Log(x)
			}
		} else {
			t.Fatal(err)
		}
	}
}
