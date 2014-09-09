package sqlite

import (
	"errors"
	"regexp"
	"strconv"
	"strings"

	sqlite3 "github.com/mxk/go-sqlite/sqlite3"
)

// close sqlite3 statement routine
func closeStatement(statement *sqlite3.Stmt) {
	if statement != nil {
		statement.Close()
	}
}

// formats SQL query error for output to log
func sqlError(SQL string, err error) error {
	return errors.New("SQL \"" + SQL + "\" error: " + err.Error())
}

// returns string that represents value for SQL query
func convertValueForSQL(value interface{}) string {
	result := ""

	switch typedValue := value.(type) {
	case string:
		result = strings.Replace(typedValue, "'", "''", -1)
		result = strings.Replace(result, "\\", "\\\\", -1)
		result = "'" + typedValue + "'"

	case []string:
		result := ""
		for _, item := range typedValue {
			if result != "" {
				result += ", "
			}
			result += convertValueForSQL(item)
		}
		return result
	}

	return result
}

// returns type used inside sqlite for given general name
func GetDBType(ColumnType string) (string, error) {
	ColumnType = strings.ToLower(ColumnType)
	switch {
	case ColumnType == "id" || ColumnType == "int" || ColumnType == "integer":
		return "INTEGER", nil
	case ColumnType == "real" || ColumnType == "float":
		return "REAL", nil
	case ColumnType == "string" || ColumnType == "text" || strings.Contains(ColumnType, "char"):
		return "TEXT", nil
	case ColumnType == "blob" || ColumnType == "struct" || ColumnType == "data":
		return "BLOB", nil
	case strings.Contains(ColumnType, "numeric") || strings.Contains(ColumnType, "decimal") || ColumnType == "money":
		return "NUMERIC", nil
	case ColumnType == "bool" || ColumnType == "boolean":
		return "NUMERIC", nil
	}

	return "?", errors.New("Unknown type '" + ColumnType + "'")
}

// makes SQL filter string based on ColumnName, Operator and Value parameters or returns nil
//   - internal usage function for AddFilter and AddStaticFilter routines
func (it *SQLiteCollection) makeSQLFilterString(ColumnName string, Operator string, Value interface{}) (string, error) {
	if it.HasColumn(ColumnName) {
		Operator = strings.ToUpper(Operator)
		if Operator == "" || Operator == "=" || Operator == "<>" || Operator == ">" || Operator == "<" || Operator == "LIKE" || Operator == "IN" {
			return ColumnName + " " + Operator + " " + convertValueForSQL(Value), nil
		} else {
			return "", errors.New("unknown operator '" + Operator + "' supposed  '', '=', '>', '<', '<>', 'LIKE', 'IN' " + ColumnName + "'")
		}
	} else {
		return "", errors.New("can't find column '" + ColumnName + "'")
	}
}

// collects all filters in a single string (for internal usage)
func (it *SQLiteCollection) getSQLFilters() string {
	sqlFilter := ""
	for _, value := range it.StaticFilters {
		if sqlFilter != "" {
			sqlFilter += " AND "
		}
		sqlFilter += value
	}

	for _, value := range it.Filters {
		if sqlFilter != "" {
			sqlFilter += " AND "
		}
		sqlFilter += value
	}

	if sqlFilter != "" {
		sqlFilter = " WHERE " + sqlFilter
	}

	return sqlFilter
}

// loads record from DB by it's id
func (it *SQLiteCollection) LoadById(id string) (map[string]interface{}, error) {
	row := make(sqlite3.RowMap)

	sqlColumns := strings.Join(it.ResultColumns, ", ")
	if sqlColumns == "" {
		sqlColumns = "*"
	}

	SQL := "SELECT " + sqlColumns + " FROM " + it.TableName + " WHERE _id = " + id

	if DEBUG_SQL {
		println(SQL)
	}

	stmt, err := it.Connection.Query(SQL)
	defer closeStatement(stmt)

	if err == nil {
		if err := stmt.Scan(row); err == nil {

			row["_id"] = strconv.FormatInt(row["_id"].(int64), 10)

			return map[string]interface{}(row), nil
		} else {
			return nil, err
		}
	} else {
		return nil, err
	}
}

// loads records from DB for current collection and filter if it set
func (it *SQLiteCollection) Load() ([]map[string]interface{}, error) {
	result := make([]map[string]interface{}, 0, 10)

	row := make(sqlite3.RowMap)

	sqlLoadFilter := it.getSQLFilters()

	sqlOrder := strings.Join(it.Order, ", ")
	if sqlOrder != "" {
		sqlOrder = " ORDER BY " + sqlOrder
	}

	sqlColumns := strings.Join(it.ResultColumns, ", ")
	if sqlColumns == "" {
		sqlColumns = "*"
	}

	SQL := "SELECT " + sqlColumns + " FROM " + it.TableName + sqlLoadFilter + sqlOrder + it.Limit

	if DEBUG_SQL {
		println(SQL)
	}

	stmt, err := it.Connection.Query(SQL)
	defer closeStatement(stmt)

	if err == nil {
		for ; err == nil; err = stmt.Next() {

			if err := stmt.Scan(row); err == nil {

				row["_id"] = strconv.FormatInt(row["_id"].(int64), 10)

				result = append(result, row)
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

	return "", nil
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

// adds selection filter that will not be cleared by ClearFilters() function
func (it *SQLiteCollection) AddStaticFilter(ColumnName string, Operator string, Value interface{}) error {

	sqlFilter, err := it.makeSQLFilterString(ColumnName, Operator, Value)
	if err != nil {
		return err
	}
	it.StaticFilters[ColumnName] = sqlFilter

	return nil
}

// adds selection filter to current collection(table) object
func (it *SQLiteCollection) AddFilter(ColumnName string, Operator string, Value interface{}) error {

	sqlFilter, err := it.makeSQLFilterString(ColumnName, Operator, Value)
	if err != nil {
		return err
	}
	it.Filters[ColumnName] = sqlFilter

	return nil
}

// removes all filters that were set for current collection
func (it *SQLiteCollection) ClearFilters() error {
	it.Filters = make(map[string]string)
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
