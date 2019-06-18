package db

import (
	"testing"
	"fmt"
	"time"
	"github.com/ottemo/commerce/utils"
)

// TestSimple is a general test for database adapter
func TestSimple(t *testing.T) {

	const collectionName = "db_test"

	// obtaining current DB engine
	db := GetDBEngine()
	if db == nil {
		t.Fatal("Cant get database engine")
	}

	// TODO: there is no function in DBEngine to remove existing collection
	// making new collection if it does not exists
	if !db.HasCollection(collectionName) {
		if err := db.CreateCollection(collectionName); err != nil {
			t.Fatal(err)
		}
	}

	// getting collection instance
	collection, err := db.GetCollection(collectionName)
	if err != nil {
		t.Fatal(err)
	}

	// truncating it
	collection.Delete()

	// defining required collection columns
	columns := map[string]string {
		"bool": ConstTypeBoolean,
		"char": ConstTypeVarchar,
		"text": ConstTypeText,
		"int": ConstTypeInteger,
		"decimal": ConstTypeDecimal,
		"money": ConstTypeMoney,
		"float": ConstTypeFloat,
		"date": ConstTypeDatetime,
		"json": ConstTypeJSON,
	}

	// making necessary collection columns
	for name, kind := range columns {
		if !collection.HasColumn(name) {
			if err := collection.AddColumn(name, kind, true); err != nil {
				t.Fatal(err)
			}
		} else if collection.GetColumnType(name) != kind {
			if err:= collection.RemoveColumn(name); err != nil {
				t.Fatal(err)
			}
			if err := collection.AddColumn(name, kind, true); err != nil {
				t.Fatal(err)
			}
		}
	}

	// checking created columns
	collectionColumns := collection.ListColumns()
	for name, kind := range columns {
		if value, present := collectionColumns[name]; !present || value != kind {
			t.Fatal("test collection column " + name + " issue: '" + kind + "' != '" + value + "'")
		}
	}

	// filling collection with sample data
	values := make([]map[string]interface{}, 0, 10)
	ids := make([]string, 0, 10)
	for i := 0; i < 10; i++ {
		value := map[string]interface{} {
			"bool": i % 2 == 0,
			"char": fmt.Sprintf("char %v", i),
			"text": fmt.Sprintf("text %v", i),
			"int": i,
			"decimal": utils.Round(float64(i) / 1.1, 0.5, 4),
			"money": utils.RoundPrice(float64(i) / 2.2),
			// TODO: figure out if the precision 6 is ok as utils.InterfaceToString converts to that precision
			"float": utils.Round(float64(i) / 3.3, 0.5, 6),
			"date": time.Now().Add(time.Hour * time.Duration(i)),
			"json": map[string]interface{} {
				"a": i % 2,
				"b": fmt.Sprintf("test %v", i),
				"c": i,
				"d": float64(i) / 1.1,
			},
		}

		id, err := collection.Save(value)
		if err != nil {
			t.Fatal(err, fmt.Sprintf("can't store: %v", value))
		}

		values = append(values,value)
		ids = append(ids, id)
	}

	// checking filled data
	for i, id := range ids {
		data, err := collection.LoadByID(id)
		if err != nil {
			t.Fatal(fmt.Sprintf("item %v collection.LoadByID(%s): %v", i, id, err))
		}

		for key, expect := range values[i] {
			value, present := data[key]
			if key == "date" {
				valueTime := utils.InterfaceToTime(value).Unix()
				expectTime := utils.InterfaceToTime(expect).Unix()
				if valueTime != expectTime {
					t.Fatal(fmt.Sprintf("key %v does not match (%v != %v)", key, expectTime, valueTime))
				}
			} else if !present || (key != "json" && value != expect) {
				t.Fatal(fmt.Sprintf("key %v does not match: expect - %v, got - %v)", key, expect, value))
			}
		}
	}

	// checking "distinct" operation
	collection.ClearFilters()
	if data, err := collection.Distinct("bool"); err != nil || len(data) != 2 {
		t.Fatal("distinct operation failed")
	}

	// checking collection filters
	collection.SetupFilterGroup("default", true, "")
	collection.SetupFilterGroup("case1", false, "default")
	collection.SetupFilterGroup("case2", false, "default")
	collection.AddGroupFilter("case1", "bool", "=", true)
	collection.AddGroupFilter("case1", "int", "<", 5)
	collection.AddGroupFilter("case2", "int", "=",  9)

	// checking "count" operation
	if cnt, err := collection.Count(); err != nil || cnt != 4 {
		t.Fatal(fmt.Sprintf("invalid filter result count %v != 4", cnt))
	}

	// checking sorting and limiting functionality
	collection.AddSort("int", true)
	collection.SetLimit(1, 2)

	data, err := collection.Load()
	if err != nil || len(data) != 2 {
		t.Fatal(fmt.Sprintf("invalid collection load result count %v != 2", len(data)))
	}

	if value, present := data[0]["int"]; !present || value != 4 {
		t.Fatal(fmt.Sprintf("un-expected value loaded: %v, expected: 4", value))
	}
}
