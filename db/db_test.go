package db

import (
	"testing"
	"fmt"
	"time"
)

func TestSimple(t *testing.T) {
	db := GetDBEngine()
	if db == nil {
		t.Fatal("Cant get database engine")
	}

	if !db.HasCollection("db_test") {
		db.CreateCollection("db_test")
	}

	collection, err := db.GetCollection("db_test")
	if err != nil {
		t.Fatal(err)
	}

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

	for name, kind := range columns {
		if !collection.HasColumn(name) || collection.GetColumnType(name) != kind {
			if err := collection.AddColumn(name, kind, true); err != nil {
				t.Fatal(err)
			}
		}
	}

	collectionColumns := collection.ListColumns()
	for name, kind := range columns {
		if value, present := collectionColumns[name]; !present || value != kind {
			t.Fatal("test collection column " + name + " issue: '" + kind + "' != '" + value + "'")
		}
	}

	values := make([]map[string]interface{}, 0, 10)
	ids := make([]string, 0, 10)
	for i := 0; i < 10; i++ {
		value := map[string]interface{} {
			"bool": i % 2 == 0,
			"char": fmt.Sprintf("char %v", i),
			"text": fmt.Sprintf("text %v", i),
			"int": i,
			"decimal": float64(i) / 1.1,
			"money": float64(i) / 2.2,
			"float": float64(i) / 3.3,
			"date": time.Now().Add(time.Hour * i * 5),
			"json": map[string]interface{} {
				"a": i % 2,
				"b": fmt.Sprintf("test %v", i),
				"c": i,
				"d": float64(i) / 1.1,
			},
		}

		id, err := collection.Save(value)
		if err != nil {
			t.Fatal(fmt.Sprintf("can't store: %v", value))
		}

		values = append(values,value)
		ids = append(ids, id)
	}

	for i, id := range ids {
		data, err := collection.LoadByID(id)
		if err != nil {
			t.Fatal(err)
		}

		for key, expect := range values[i] {
			value, present := data[key]
			if !present || (key != "json" && value != expect) {
				t.Fatal(fmt.Sprintf("key %v does not match (%v=%v)", key, expect, value))
			}
		}
	}

	if data, err := collection.Distinct("bool"); err != nil || len(data) != 2 {
		t.Fatal("distinct operation failed")
	}

	collection.AddFilter("bool", "=", true)
	collection.AddFilter("int", "<", 5)
	collection.SetupFilterGroup("or", true, "")
	collection.AddGroupFilter("or", "int", "=",  10)

	if cnt, err := collection.Count(); err != nil || cnt != 4 {
		t.Fatal(fmt.Sprintf("invalid filter result count %v != 4", cnt))
	}

	collection.AddSort("i", true)
	collection.SetLimit(1, 2)

	data, err := collection.Load()
	if err != nil || len(data) != 2 {
		t.Fatal(fmt.Sprintf("invalid collection load result count %v != 2", len(data)))
	}

	if value, present := data[0]["int"]; !present || value != 10 {
		t.Fatal(fmt.Sprintf("un-expected value loaded: %v", value))
	}

}
