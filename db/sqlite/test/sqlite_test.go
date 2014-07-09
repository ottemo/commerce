package test

import (
	"github.com/ottemo/foundation/database/sqlite"
	"testing"
)

func TestCollection(t *testing.T) {
	sqliteEngine := new(sqlite.SQLite)

	sqliteEngine.Startup()

	collection, _ := sqliteEngine.GetCollection("test")
	if err := collection.CreateCollection(); err == nil {
		if err := collection.AddColumn("name", "varchar", false); err != nil {
			t.Error(err)
		}
		if err := collection.AddColumn("value", "int", false); err != nil {
			t.Error(err)
		}

		if err := collection.AddColumn("deleteme", "text", false); err != nil {
			t.Error(err)
		}
		if err := collection.RemoveColumn("deleteme"); err != nil {
			t.Error(err)
		}
	}

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
