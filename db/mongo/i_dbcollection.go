package mongo

import (
	"errors"
	"sort"
	"strings"

	"labix.org/v2/mgo/bson"
)

const (
	COLUMN_INFO_COLLECTION = "collection_column_info"
)

// returns string that represents value for SQL query
func sqlError(SQL string, err error) error {
	return errors.New("SQL \"" + SQL + "\" error: " + err.Error())
}

//func getDBType(ColumnType string) (string, error) {
//	ColumnType = strings.ToLower(ColumnType)
//	switch ColumnType {
//	case ColumnType == "int" || ColumnType == "integer":
//		return "INTEGER", nil
//	case ColumnType == "real" || ColumnType == "float":
//		return "REAL", nil
//	case ColumnType == "string" || ColumnType == "text" || strings.Contains(ColumnType, "char"):
//		return "TEXT", nil
//	case ColumnType == "blob" || ColumnType == "struct" || ColumnType == "data":
//		return "BLOB", nil
//	case strings.Contains(ColumnType, "numeric") || strings.Contains(ColumnType, "decimal") || ColumnType == "money":
//		return "NUMERIC", nil
//	}
//
//	return "?", errors.New("Unknown type '" + ColumnType + "'")
//}

// converts well knows filter operator used for SQL to mongoDB one if possible
func getMongoOperator(Operator string, Value interface{}) (string, interface{}, error) {
	Operator = strings.ToLower(Operator)

	switch Operator {
	case "=":
		return "", Value, nil
	case ">":
		return "$gt", Value, nil
	case ">=":
		return "$gte", Value, nil
	case "<":
		return "$lt", Value, nil
	case "<=":
		return "$lte", Value, nil
	case "like":
		Value, ok := Value.(string)
		if ok {
			Value = strings.Replace("%", ".*", Value, -1)
			return "$regex", Value, nil
		}
	}

	return "?", "?", errors.New("Unknown operator '" + Operator + "'")
}

// internal usage function for AddFilter and AddStaticFilter routines
func (it *MongoDBCollection) makeSelector(ColumnName string, Operator string, Value interface{}) (interface{}, error) {
	newOperator, newValue, err := getMongoOperator(Operator, Value)
	if err != nil {
		return nil, err
	}

	if newOperator != "" {
		return map[string]interface{}{newOperator: newValue}, nil
	} else {
		return newValue, nil
	}
}

// function to join static filters with dynamic filters in one selector
func (it *MongoDBCollection) joinSelectors() interface{} {
	result := make(map[string]interface{})

	for column, value := range it.StaticSelector {
		result[column] = value
	}

	for column, value := range it.Selector {
		if prevValue, present := result[column]; present {
			result[column] = map[string]interface{}{"$and:": []interface{}{prevValue, value}}
		} else {
			result[column] = value
		}
	}

	return result
}

// loads record from DB by it's id
func (it *MongoDBCollection) LoadById(id string) (map[string]interface{}, error) {
	result := make(map[string]interface{})

	err := it.collection.FindId(id).One(&result)

	return result, err
}

// loads records from DB for current collection and filter if it set
func (it *MongoDBCollection) Load() ([]map[string]interface{}, error) {
	result := make([]map[string]interface{}, 0)

	query := it.collection.Find(it.joinSelectors())

	if len(it.Sort) > 0 {
		query.Sort(it.Sort...)
	}

	if it.Offset > 0 {
		query = query.Skip(it.Offset)
	}
	if it.Limit > 0 {
		query = query.Limit(it.Limit)
	}

	err := query.All(&result)

	return result, err
}

// returns count of rows matching current select statement
func (it *MongoDBCollection) Count() (int, error) {
	return it.collection.Find(it.Selector).Count()
}

// stores record in DB for current collection
func (it *MongoDBCollection) Save(Item map[string]interface{}) (string, error) {

	// id validation/updating
	//-----------------------
	id := bson.NewObjectId().Hex()

	if _id, present := Item["_id"]; present {
		if _id, ok := _id.(string); ok && _id != "" {
			if bson.IsObjectIdHex(_id) {
				id = _id
			}
		}
	}
	Item["_id"] = id

	// sorting by attribute name
	//--------------------------
	bsonDocument := make(bson.D, 0, len(Item))
	keysList := make([]string, 0, len(Item))

	for key, _ := range Item {
		keysList = append(keysList, key)
	}
	sort.Strings(keysList)

	for _, key := range keysList {
		bsonDocument = append(bsonDocument, bson.DocElem{Name: key, Value: Item[key]})
	}

	// saving document to DB
	//----------------------
	changeInfo, err := it.collection.UpsertId(id, bsonDocument)

	if changeInfo != nil && changeInfo.UpsertedId != nil {
		//id = changeInfo.UpsertedId
	}

	return id, err
}

// removes records that matches current select statement from DB
//   - returns amount of affected rows
func (it *MongoDBCollection) Delete() (int, error) {
	changeInfo, err := it.collection.RemoveAll(it.Selector)

	return changeInfo.Removed, err
}

// removes record from DB by is's id
func (it *MongoDBCollection) DeleteById(id string) error {

	return it.collection.RemoveId(id)
}

// adds selection filter that will not be cleared by ClearFilters() function
func (it *MongoDBCollection) AddStaticFilter(ColumnName string, Operator string, Value interface{}) error {
	selector, err := it.makeSelector(ColumnName, Operator, Value)
	if err != nil {
		return err
	}
	it.StaticSelector[ColumnName] = selector

	return nil
}

// adds selection filter to current collection object
func (it *MongoDBCollection) AddFilter(ColumnName string, Operator string, Value interface{}) error {
	selector, err := it.makeSelector(ColumnName, Operator, Value)
	if err != nil {
		return err
	}
	it.Selector[ColumnName] = selector

	return nil
}

// removes all filters that were set for current collection
func (it *MongoDBCollection) ClearFilters() error {
	it.Selector = make(map[string]interface{})
	return nil
}

// adds sorting for current collection
func (it *MongoDBCollection) AddSort(ColumnName string, Desc bool) error {
	if Desc {
		it.Sort = append(it.Sort, "-"+ColumnName)
	} else {
		it.Sort = append(it.Sort, ColumnName)
	}
	return nil
}

// removes any sorting that was set for current collection
func (it *MongoDBCollection) ClearSort() error {
	it.Sort = make([]string, 0)
	return nil
}

// results pagination
func (it *MongoDBCollection) SetLimit(Offset int, Limit int) error {
	it.Limit = Limit
	it.Offset = Offset

	return nil
}

// limits column selection for Load() and LoadById()function
func (it *MongoDBCollection) SetResultColumns(columns ...string) error {
	for _, columnName := range columns {
		if !it.HasColumn(columnName) {
			return errors.New("there is no column " + columnName + " found")
		}

		it.ResultAttributes = append(it.ResultAttributes, columnName)
	}
	return nil
}

// Collection columns stuff
//--------------------------

// returns attributes available for current collection
func (it *MongoDBCollection) ListColumns() map[string]string {

	result := map[string]string{}

	infoCollection := it.database.C(COLUMN_INFO_COLLECTION)
	selector := map[string]string{"collection": it.Name}
	iter := infoCollection.Find(selector).Iter()

	row := map[string]string{}
	for iter.Next(&row) {
		colName, okColumn := row["column"]
		colType, okType := row["type"]

		if okColumn && okType {
			result[colName] = colType
		}
	}

	return result
}

// check for attribute presence in current collection
func (it *MongoDBCollection) HasColumn(ColumnName string) bool {

	infoCollection := it.database.C(COLUMN_INFO_COLLECTION)
	selector := map[string]interface{}{"collection": it.Name, "column": ColumnName}
	count, _ := infoCollection.Find(selector).Count()

	return count > 0
}

// adds new attribute to current collection
func (it *MongoDBCollection) AddColumn(ColumnName string, ColumnType string, indexed bool) error {

	infoCollection := it.database.C(COLUMN_INFO_COLLECTION)

	selector := map[string]interface{}{"collection": it.Name, "column": ColumnName}
	data := map[string]interface{}{"collection": it.Name, "column": ColumnName, "type": ColumnType, "indexed": indexed}

	_, err := infoCollection.Upsert(selector, data)

	return err
}

// removes attribute from current collection
//   - for MoongoDB it means update "collection_column_info" collection
//   and, update all objects of current collection to exclude attribute
func (it *MongoDBCollection) RemoveColumn(ColumnName string) error {

	infoCollection := it.database.C(COLUMN_INFO_COLLECTION)
	removeSelector := map[string]string{"collection": it.Name, "column": ColumnName}

	err := infoCollection.Remove(removeSelector)
	if err != nil {
		return err
	}

	updateSelector := map[string]interface{}{ColumnName: map[string]interface{}{"$exists": true}}
	data := map[string]interface{}{"$unset": map[string]interface{}{ColumnName: ""}}

	_, err = it.collection.UpdateAll(updateSelector, data)

	if err != nil {
		return err
	}

	return nil
}
