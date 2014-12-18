package mongo

import (
	"errors"
	"strings"

	"gopkg.in/mgo.v2/bson"
)

// Constants for the mongodb package
const (
	COLUMN_INFO_COLLECTION = "collection_column_info"
)

func sqlError(SQL string, err error) error {
	return errors.New("SQL \"" + SQL + "\" error: " + err.Error())
}

func getMongoOperator(Operator string, Value string) (string, string, error) {
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
		Value = strings.Replace("%", ".*", Value, -1)
		return "$regex", Value, nil
	}

	return "?", "?", errors.New("Unknown operator '" + Operator + "'")
}

// LoadById will return a map of interfaces when provided an ID.
func (it *MongoDBCollection) LoadById(id string) (map[string]interface{}, error) {
	result := make(map[string]interface{})

	err := it.collection.FindId(id).One(&result)

	return result, err
}

// Load will return a map of interfaces matching the the Collection selector.
func (it *MongoDBCollection) Load() ([]map[string]interface{}, error) {
	result := make([]map[string]interface{}, 0)

	query := it.collection.Find(it.Selector)

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

// Save will persist the map of interfaces to the database.
func (it *MongoDBCollection) Save(Item map[string]interface{}) (string, error) {

	id := bson.NewObjectId().Hex()

	if _id, present := Item["_id"]; present {
		if _id, ok := _id.(string); ok && _id != "" {
			if bson.IsObjectIdHex(_id) {
				id = _id
			}
		}
	}

	Item["_id"] = id

	changeInfo, err := it.collection.UpsertId(id, Item)

	if changeInfo != nil && changeInfo.UpsertedId != nil {
		//id = changeInfo.UpsertedId
	}

	return id, err
}

// Delete will remove the given Selector from the database.
func (it *MongoDBCollection) Delete() (int, error) {
	changeInfo, err := it.collection.RemoveAll(it.Selector)

	return changeInfo.Removed, err
}

// DeleteById will remove the data identified by ID.
func (it *MongoDBCollection) DeleteById(id string) error {

	return it.collection.RemoveId(id)
}

// AddFilter will insert a new Filter.
func (it *MongoDBCollection) AddFilter(ColumnName string, Operator string, Value string) error {

	newOperator, newValue, err := getMongoOperator(Operator, Value)
	if err != nil {
		return err
	}

	var filterValue interface{} = newValue
	if newOperator != "" {
		filterValue = map[string]interface{}{newOperator: newValue}
	} else {
		filterValue = newValue
	}

	it.Selector[ColumnName] = filterValue

	return nil
}

// ClearFilters will remove all filters for the Collection.
func (it *MongoDBCollection) ClearFilters() error {
	it.Selector = make(map[string]interface{})
	return nil
}

// AddSort will append a new Sort.
func (it *MongoDBCollection) AddSort(ColumnName string, Desc bool) error {
	if Desc {
		it.Sort = append(it.Sort, "-"+ColumnName)
	} else {
		it.Sort = append(it.Sort, ColumnName)
	}
	return nil
}

// ClearSort will clear all Sorts.
func (it *MongoDBCollection) ClearSort() error {
	it.Sort = make([]string, 0)
	return nil
}

// SetLimit will set the Offset and Limit for the Collection.
func (it *MongoDBCollection) SetLimit(Offset int, Limit int) error {
	it.Limit = Limit
	it.Offset = Offset

	return nil
}

// Collection columns stuff
//--------------------------

// ListColumns will return the Columns of the Collection as map of strings.
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

// HasColumn will return true/false if the current Collection has a Column.
func (it *MongoDBCollection) HasColumn(ColumnName string) bool {

	infoCollection := it.database.C(COLUMN_INFO_COLLECTION)
	selector := map[string]interface{}{"collection": it.Name, "column": ColumnName}
	count, _ := infoCollection.Find(selector).Count()

	return count > 0
}

// AddColumn will add a Column to the current Collection.
func (it *MongoDBCollection) AddColumn(ColumnName string, ColumnType string, indexed bool) error {

	infoCollection := it.database.C(COLUMN_INFO_COLLECTION)

	selector := map[string]interface{}{"collection": it.Name, "column": ColumnName}
	data := map[string]interface{}{"collection": it.Name, "column": ColumnName, "type": ColumnType, "indexed": indexed}

	_, err := infoCollection.Upsert(selector, data)

	return err
}

// RemoveColumn will remove the givne Column from the current Collection.
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
