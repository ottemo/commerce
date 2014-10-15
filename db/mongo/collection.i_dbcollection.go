package mongo

import (
	"github.com/ottemo/foundation/env"
	"labix.org/v2/mgo/bson"
	"sort"
)

// loads one record from DB by record _id
func (it *MongoDBCollection) LoadById(id string) (map[string]interface{}, error) {
	result := make(map[string]interface{})

	err := it.collection.FindId(id).One(&result)

	return result, env.ErrorDispatch(err)
}

// loads records from DB for current collection and filter if it set
func (it *MongoDBCollection) Load() ([]map[string]interface{}, error) {
	result := make([]map[string]interface{}, 0)

	err := it.prepareQuery().All(&result)

	return result, env.ErrorDispatch(err)
}

// applies [iterator] function to each record, stops on return false
func (it *MongoDBCollection) Iterate(iteratorFunc func(record map[string]interface{}) bool) error {
	record := make(map[string]interface{})

	iterator := it.prepareQuery().Iter()
	for iterator.Next(&record) {
		proceed := iteratorFunc(record)

		if !proceed {
			err := iterator.Close()
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// returns count of rows matching current select statement
func (it *MongoDBCollection) Count() (int, error) {
	return it.collection.Find(it.makeSelector()).Count()
}

// returns distinct values of specified attribute
func (it *MongoDBCollection) Distinct(columnName string) ([]interface{}, error) {
	result := make([]interface{}, 0)

	err := it.prepareQuery().Distinct(columnName, &result)

	return result, env.ErrorDispatch(err)
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

	return id, env.ErrorDispatch(err)
}

// removes records that matches current select statement from DB, returns amount of affected rows
func (it *MongoDBCollection) Delete() (int, error) {
	changeInfo, err := it.collection.RemoveAll(it.makeSelector())

	return changeInfo.Removed, env.ErrorDispatch(err)
}

// removes record from DB by is's id
func (it *MongoDBCollection) DeleteById(id string) error {
	return it.collection.RemoveId(id)
}

// setups filter group params for collection
func (it *MongoDBCollection) SetupFilterGroup(groupName string, orSequence bool, parentGroup string) error {
	if _, present := it.FilterGroups[parentGroup]; !present && parentGroup != "" {
		return env.ErrorNew("invalid parent group")
	}

	filterGroup := it.getFilterGroup(groupName)
	filterGroup.OrSequence = orSequence
	filterGroup.ParentGroup = parentGroup

	return nil
}

// removes filter group for collection
func (it *MongoDBCollection) RemoveFilterGroup(GroupName string) error {
	if _, present := it.FilterGroups[GroupName]; !present {
		return env.ErrorNew("invalid group name")
	}
	delete(it.FilterGroups, GroupName)
	return nil
}

// adds selection filter to specific filter group (all filter groups will be joined before db query)
func (it *MongoDBCollection) AddGroupFilter(GroupName string, ColumnName string, Operator string, Value interface{}) error {
	err := it.updateFilterGroup(GroupName, ColumnName, Operator, Value)
	if err != nil {
		return err
	}
	return nil
}

// adds selection filter that will not be cleared by ClearFilters() function
func (it *MongoDBCollection) AddStaticFilter(ColumnName string, Operator string, Value interface{}) error {
	err := it.updateFilterGroup(FILTER_GROUP_STATIC, ColumnName, Operator, Value)
	if err != nil {
		return err
	}
	return nil
}

// adds selection filter to current collection object
func (it *MongoDBCollection) AddFilter(ColumnName string, Operator string, Value interface{}) error {
	err := it.updateFilterGroup(FILTER_GROUP_DEFAULT, ColumnName, Operator, Value)
	if err != nil {
		return err
	}
	return nil
}

// removes all filters that were set for current collection
func (it *MongoDBCollection) ClearFilters() error {
	for filterGroup, _ := range it.FilterGroups {
		if filterGroup != FILTER_GROUP_STATIC {
			delete(it.FilterGroups, filterGroup)
		}
	}
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
		it.ResultAttributes = []string{}

		if !it.HasColumn(columnName) {
			return env.ErrorNew("there is no column " + columnName + " found")
		}

		it.ResultAttributes = append(it.ResultAttributes, columnName)
	}
	return nil
}

// returns attributes available for current collection
func (it *MongoDBCollection) ListColumns() map[string]string {

	result := map[string]string{}

	infoCollection := it.database.C(COLLECTION_NAME_COLUMN_INFO)
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

	// updating cached attribute types information
	if _, present := attributeTypes[it.Name]; !present {
		attributeTypes[it.Name] = make(map[string]string)
	}

	attributeTypesMutex.Lock()
	for attributeName, attributeType := range result {
		attributeTypes[it.Name][attributeName] = attributeType
	}
	attributeTypesMutex.Unlock()

	return result
}

// returns SQL like type of attribute in current collection, or if not present ""
func (it *MongoDBCollection) GetColumnType(ColumnName string) string {
	// _id - has static type
	if ColumnName == "_id" {
		return "string"
	}

	// looking in cache first
	attributeType, present := attributeTypes[it.Name][ColumnName]
	if !present {
		// updating cache, and looking again
		it.ListColumns()
		attributeType, present = attributeTypes[it.Name][ColumnName]
	}

	return attributeType
}

// check for attribute presence in current collection
func (it *MongoDBCollection) HasColumn(ColumnName string) bool {
	// _id - always present
	if ColumnName == "_id" {
		return true
	}

	// looking in cache first
	_, present := attributeTypes[it.Name][ColumnName]
	if !present {
		// updating cache, and looking again
		it.ListColumns()
		_, present = attributeTypes[it.Name][ColumnName]
	}

	return present
}

// adds new attribute to current collection
func (it *MongoDBCollection) AddColumn(ColumnName string, ColumnType string, indexed bool) error {

	infoCollection := it.database.C(COLLECTION_NAME_COLUMN_INFO)

	selector := map[string]interface{}{"collection": it.Name, "column": ColumnName}
	data := map[string]interface{}{"collection": it.Name, "column": ColumnName, "type": ColumnType, "indexed": indexed}

	_, err := infoCollection.Upsert(selector, data)

	return env.ErrorDispatch(err)
}

// removes attribute from current collection
//   - for MoongoDB it means update "collection_column_info" collection
//   and, update all objects of current collection to exclude attribute
func (it *MongoDBCollection) RemoveColumn(ColumnName string) error {

	infoCollection := it.database.C(COLLECTION_NAME_COLUMN_INFO)
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
