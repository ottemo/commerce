package sqlite

import (
	"errors"
	"io"
	"regexp"
	"strconv"
	"strings"

	sqlite3 "github.com/mxk/go-sqlite/sqlite3"
)

// loads record from DB by it's id
func (it *SQLiteCollection) LoadById(id string) (map[string]interface{}, error) {
	var result map[string]interface{} = nil

	err := it.Iterate(func(row map[string]interface{}) bool {
		result = row
		return false
	})

	return result, err
}

// loads records from DB for current collection and filter if it set
func (it *SQLiteCollection) Load() ([]map[string]interface{}, error) {
	result := make([]map[string]interface{}, 0)

	err := it.Iterate(func(row map[string]interface{}) bool {
		result = append(result, row)
		return true
	})

	return result, err
}

// applies [iterator] function to each record, stops on return false
func (it *SQLiteCollection) Iterate(iteratorFunc func(record map[string]interface{}) bool) error {

	SQL := "SELECT " + it.getSQLResultColumns() + " FROM " + it.TableName + it.getSQLFilters() + it.getSQLOrder() + it.Limit
	if DEBUG_SQL {
		println(SQL)
	}

	stmt, err := it.Connection.Query(SQL)
	defer closeStatement(stmt)

	if err != nil {
		return sqlError(SQL, err)
	}

	row := make(sqlite3.RowMap)
	for ; err == nil; err = stmt.Next() {
		if err := stmt.Scan(row); err == nil {
			it.modifyResultRowId(row)

			if !iteratorFunc(row) {
				break
			}
		}
	}

	if err == io.EOF {
		err = nil
	}

	return err
}

// returns distinct values of specified attribute
func (it *SQLiteCollection) Distinct(columnName string) ([]interface{}, error) {

	if len(it.ResultColumns) != 1 {
		return nil, errors.New("should be 1 result column")
	}

	SQL := "SELECT " + it.getSQLResultColumns() + " FROM " + it.TableName + it.getSQLFilters() + it.getSQLOrder() + it.Limit
	if DEBUG_SQL {
		println(SQL)
	}

	stmt, err := it.Connection.Query(SQL)
	defer closeStatement(stmt)

	result := make([]interface{}, 0)
	if err == nil {
		row := make(sqlite3.RowMap)

		for ; err == nil; err = stmt.Next() {
			if err := stmt.Scan(row); err == nil {
				it.modifyResultRowId(row)

				for _, columnValue := range row {
					result = append(result, columnValue)
					break
				}
			} else {
				return result, err
			}
		}

	} else {
		err = sqlError(SQL, err)
	}

	return result, err
}

// returns count of rows matching current select statement
func (it *SQLiteCollection) Count() (int, error) {
	sqlLoadFilter := it.getSQLFilters()

	row := make(sqlite3.RowMap)

	SQL := "SELECT COUNT(*) AS cnt FROM " + it.TableName + sqlLoadFilter

	if DEBUG_SQL {
		println(SQL)
	}

	stmt, err := it.Connection.Query(SQL)
	defer closeStatement(stmt)

	if err == nil {
		if err := stmt.Scan(row); err == nil {
			cnt := int(row["cnt"].(int64)) //TODO: check this assertion works
			return cnt, err
		} else {
			return 0, err
		}
	} else {
		return 0, err
	}
}

// stores record in DB for current collection
func (it *SQLiteCollection) Save(Item map[string]interface{}) (string, error) {

	// _id in SQLite supposed to be auto-incremented int but for MongoDB it forced to be string
	// collection interface also forced us to use string but we still want it ti be int in DB
	// to make that we need to convert it before save from  string to int or nil
	// and after save get auto-incremented id as convert to string
	if Item["_id"] != nil {
		if intValue, err := strconv.ParseInt(Item["_id"].(string), 10, 64); err == nil {
			Item["_id"] = intValue
		} else {
			Item["_id"] = nil
		}
	}

	// SQL generation
	columns := make([]string, 0, len(Item))
	args := make([]string, 0, len(Item))
	values := make([]interface{}, 0, len(Item))

	for k, v := range Item {
		columns = append(columns, "\""+k+"\"")
		args = append(args, "$_"+k)
		values = append(values, v)
	}

	SQL := "INSERT OR REPLACE INTO  " + it.TableName +
		" (" + strings.Join(columns, ",") + ") VALUES " +
		" (" + strings.Join(args, ",") + ")"

	if DEBUG_SQL {
		println(SQL)
	}

	if err := it.Connection.Exec(SQL, values...); err == nil {

		// auto-incremented _id back to string
		newIdInt := it.Connection.LastInsertId()
		newIdString := strconv.FormatInt(newIdInt, 10)
		Item["_id"] = newIdString

		return newIdString, nil
	} else {
		return "", sqlError(SQL, err)
	}
}

// removes records that matches current select statement from DB
//   - returns amount of affected rows
func (it *SQLiteCollection) Delete() (int, error) {
	sqlDeleteFilter := it.getSQLFilters()

	SQL := "DELETE FROM " + it.TableName + sqlDeleteFilter

	if DEBUG_SQL {
		println(SQL)
	}

	err := it.Connection.Exec(SQL)
	affected := it.Connection.RowsAffected()

	return affected, err
}

// removes record from DB by is's id
func (it *SQLiteCollection) DeleteById(id string) error {
	SQL := "DELETE FROM " + it.TableName + " WHERE _id = " + id

	if DEBUG_SQL {
		println(SQL)
	}

	return it.Connection.Exec(SQL)
}

// setups filter group params for collection
func (it *SQLiteCollection) SetupFilterGroup(groupName string, orSequence bool, parentGroup string) error {
	if _, present := it.FilterGroups[parentGroup]; !present && parentGroup != "" {
		return errors.New("invalid parent group")
	}

	filterGroup := it.getFilterGroup(groupName)
	filterGroup.OrSequence = orSequence
	filterGroup.ParentGroup = parentGroup

	return nil
}

// removes filter group for collection
func (it *SQLiteCollection) RemoveFilterGroup(groupName string) error {
	if _, present := it.FilterGroups[groupName]; !present {
		return errors.New("invalid group name")
	}
	delete(it.FilterGroups, groupName)
	return nil
}

// adds selection filter to specific filter group (all filter groups will be joined before db query)
func (it *SQLiteCollection) AddGroupFilter(groupName string, columnName string, operator string, value interface{}) error {
	err := it.updateFilterGroup(groupName, columnName, operator, value)
	if err != nil {
		return err
	}
	return nil
}

// adds selection filter that will not be cleared by ClearFilters() function
func (it *SQLiteCollection) AddStaticFilter(columnName string, operator string, value interface{}) error {

	err := it.updateFilterGroup(FILTER_GROUP_STATIC, columnName, operator, value)
	if err != nil {
		return err
	}
	return nil
}

// adds selection filter to current collection(table) object
func (it *SQLiteCollection) AddFilter(ColumnName string, Operator string, Value interface{}) error {

	err := it.updateFilterGroup(FILTER_GROUP_DEFAULT, ColumnName, Operator, Value)
	if err != nil {
		return err
	}
	return nil
}

// removes all filters that were set for current collection, except static
func (it *SQLiteCollection) ClearFilters() error {
	for filterGroup, _ := range it.FilterGroups {
		if filterGroup != FILTER_GROUP_STATIC {
			delete(it.FilterGroups, filterGroup)
		}
	}
	return nil
}

// adds sorting for current collection
func (it *SQLiteCollection) AddSort(ColumnName string, Desc bool) error {
	if it.HasColumn(ColumnName) {
		if Desc {
			it.Order = append(it.Order, ColumnName+" DESC")
		} else {
			it.Order = append(it.Order, ColumnName)
		}
	} else {
		return errors.New("can't find column '" + ColumnName + "'")
	}

	return nil
}

// removes any sorting that was set for current collection
func (it *SQLiteCollection) ClearSort() error {
	it.Order = make([]string, 0)
	return nil
}

// limits column selection for Load() and LoadById()function
func (it *SQLiteCollection) SetResultColumns(columns ...string) error {
	for _, columnName := range columns {
		if !it.HasColumn(columnName) {
			return errors.New("there is no column " + columnName + " found")
		}

		it.ResultColumns = append(it.ResultColumns, columnName)
	}
	return nil
}

// results pagination
func (it *SQLiteCollection) SetLimit(Offset int, Limit int) error {
	if Limit == 0 {
		it.Limit = ""
	} else {
		it.Limit = " LIMIT " + strconv.Itoa(Limit) + " OFFSET " + strconv.Itoa(Offset)
	}
	return nil
}

// updates information about attributes(columns) for current collection(table) (loads them from DB)
func (it *SQLiteCollection) RefreshColumns() {
	SQL := "PRAGMA table_info(" + it.TableName + ")"

	if DEBUG_SQL {
		println(SQL)
	}

	row := make(sqlite3.RowMap)

	stmt, err := it.Connection.Query(SQL)
	defer closeStatement(stmt)

	for ; err == nil; err = stmt.Next() {
		stmt.Scan(row)

		key := row["name"].(string)
		value := row["type"].(string)
		it.Columns[key] = value
	}
}

// returns attributes(columns) available for current collection(table)
func (it *SQLiteCollection) ListColumns() map[string]string {
	it.RefreshColumns()
	return it.Columns
}

// returns SQL like type of attribute in current collection, or if not present ""
func (it *SQLiteCollection) GetColumnType(columnName string) string {
	if columnName == "_id" {
		return "string"
	}

	if columnType, present := it.Columns[columnName]; present {
		return columnType
	}

	return ""
}

// check for attribute(column) presence in current collection
func (it *SQLiteCollection) HasColumn(ColumnName string) bool {
	if _, present := it.Columns[ColumnName]; present {
		return true
	} else {
		it.RefreshColumns()
		_, present := it.Columns[ColumnName]
		return present
	}
}

// adds new attribute(column) to current collection(table)
func (it *SQLiteCollection) AddColumn(ColumnName string, ColumnType string, indexed bool) error {

	// TODO: there probably need column name check to be only lowercase, exclude some chars, etc.

	if it.HasColumn(ColumnName) {
		return errors.New("column '" + ColumnName + "' already exists for '" + it.TableName + "' collection")
	}

	if ColumnType, err := GetDBType(ColumnType); err == nil {

		SQL := "ALTER TABLE " + it.TableName + " ADD COLUMN \"" + ColumnName + "\" " + ColumnType

		if DEBUG_SQL {
			println(SQL)
		}

		if err := it.Connection.Exec(SQL); err == nil {
			return nil
		} else {
			return sqlError(SQL, err)
		}

	} else {
		return err
	}

}

// removes attribute(column) to current collection(table)
//   - sqlite do not have alter DROP COLUMN statements so it is hard task...
func (it *SQLiteCollection) RemoveColumn(ColumnName string) error {

	// checking column in table
	if !it.HasColumn(ColumnName) {
		return errors.New("column '" + ColumnName + "' not exists in '" + it.TableName + "' collection")
	}

	// getting table create SQL to take columns from
	//----------------------------------------------
	var tableCreateSQL string = ""

	SQL := "SELECT sql FROM sqlite_master WHERE tbl_name='" + it.TableName + "' AND type='table'"

	stmt, err := it.Connection.Query(SQL)
	if err != nil {
		return sqlError(SQL, err)
	}

	err = stmt.Scan(&tableCreateSQL)
	if err != nil {
		return err
	}

	closeStatement(stmt)

	// parsing create SQL, making same but w/o deleting column
	//--------------------------------------------------------
	tableColumnsWTypes := ""
	tableColumnsWoTypes := ""

	re := regexp.MustCompile("CREATE TABLE .*\\((.*)\\)")
	if regexMatch := re.FindStringSubmatch(tableCreateSQL); len(regexMatch) >= 2 {
		tableColumnsList := strings.Split(regexMatch[1], ", ")

		for _, tableColumn := range tableColumnsList {
			tableColumn = strings.Trim(tableColumn, "\n\t ")

			if !strings.HasPrefix(tableColumn, ColumnName) && !strings.HasPrefix(tableColumn, "\""+ColumnName+"\"") {
				if tableColumnsWTypes != "" {
					tableColumnsWTypes += ", "
					tableColumnsWoTypes += ", "
				}
				tableColumnsWTypes += tableColumn
				tableColumnsWoTypes += tableColumn[0:strings.Index(tableColumn, " ")]
			}

		}
	} else {
		return errors.New("can't find table create columns in '" + tableCreateSQL + "', found [" + strings.Join(regexMatch, ", ") + "]")
	}

	// making new table without removing column, and filling with values from old table
	//---------------------------------------------------------------------------------
	SQL = "CREATE TABLE " + it.TableName + "_removecolumn (" + tableColumnsWTypes + ") "
	if err := it.Connection.Exec(SQL); err != nil {
		return sqlError(SQL, err)
	}

	SQL = "INSERT INTO " + it.TableName + "_removecolumn (" + tableColumnsWoTypes + ") SELECT " + tableColumnsWoTypes + " FROM " + it.TableName
	if err := it.Connection.Exec(SQL); err != nil {
		return sqlError(SQL, err)
	}

	// switching newly created table, deleting old table
	//---------------------------------------------------
	SQL = "ALTER TABLE " + it.TableName + " RENAME TO " + it.TableName + "_fordelete"
	if err := it.Connection.Exec(SQL); err != nil {
		return sqlError(SQL, err)
	}

	SQL = "ALTER TABLE " + it.TableName + "_removecolumn RENAME TO " + it.TableName
	if err := it.Connection.Exec(SQL); err != nil {
		return sqlError(SQL, err)
	}

	SQL = "DROP TABLE " + it.TableName + "_fordelete"
	if err := it.Connection.Exec(SQL); err != nil {
		return sqlError(SQL, err)
	}

	return nil
}
