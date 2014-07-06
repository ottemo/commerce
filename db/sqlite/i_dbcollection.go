package sqlite

import (
	"strings"
	"strconv"
	"errors"
	"regexp"

	sqlite3 "code.google.com/p/go-sqlite/go1/sqlite3"
)

func sqlError(SQL string, err error) error {
	return errors.New("SQL \"" + SQL + "\" error: " + err.Error())
}


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
	}

	return "?", errors.New("Unknown type '" + ColumnType + "'")
}


func (it *SQLiteCollection) LoadById(id string) (map[string]interface{}, error) {
	row := make(sqlite3.RowMap)

	SQL := "SELECT * FROM " + it.TableName + " WHERE _id = " + id
	if s, err := it.Connection.Query(SQL); err == nil {
		if err := s.Scan(row); err == nil {
			return map[string]interface{}(row), nil
		} else {
			return nil, err
		}
	} else {
		return nil, err
	}
}

func (it *SQLiteCollection) Load() ([]map[string]interface{}, error) {
	result := make([] map[string]interface{}, 0, 10)
	var err error = nil

	row := make(sqlite3.RowMap)

	sqlLoadFilter := strings.Join(it.Filters, " AND ")
	if sqlLoadFilter != "" {
		sqlLoadFilter = " WHERE " + sqlLoadFilter
	}

	sqlOrder := strings.Join(it.Order, ", ")
	if sqlOrder != "" {
		sqlOrder = " ORDER BY " + sqlOrder
	}

	SQL := "SELECT * FROM " + it.TableName + sqlLoadFilter + sqlOrder + it.Limit
	if s, err := it.Connection.Query(SQL); err == nil {
		for ; err == nil; err = s.Next() {

			if err := s.Scan(row); err == nil {
				result = append(result, row )
			} else {
				return result, err
			}
		}
	} else {
		err = sqlError(SQL, err)
	}

	return result, err
}

func (it *SQLiteCollection) Save(Item map[string]interface{}) (string, error) {

	// _id in SQLite supposed to be auto-incremented int but for MongoDB it forced to be string
	// collection interface also forced us to use string but we still want it ti be int in DB
	// to make that we need to convert it before save from  string to int or nil
	// and after save get auto-incremented id as convert to string
	if Item["_id"] != nil {
		if intValue, err := strconv.ParseInt( Item["_id"].(string), 10, 64); err == nil {
			Item["_id"] = intValue
		}else{
			Item["_id"] = nil
		}
	}

	// SQL generation
	columns := make([]string, 0, len(Item))
	   args := make([]string, 0, len(Item))
	 values := make([]interface{}, 0, len(Item))
	
	for k,v := range Item {
		columns = append(columns, "\"" + k + "\"")
		   args = append(args, "$_" + k )
		 values = append(values, v)
	}

	SQL := "INSERT OR REPLACE INTO  " + it.TableName +
			" (" + strings.Join(columns, ",") + ") VALUES " +
			" (" + strings.Join(args, ",") + ")"
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

func (it *SQLiteCollection) Delete() (int, error) {
	sqlDeleteFilter := strings.Join(it.Filters, " AND ")
	if sqlDeleteFilter != "" {
		sqlDeleteFilter = " WHERE " + sqlDeleteFilter
	}

	SQL := "DELETE FROM " + it.TableName + sqlDeleteFilter
	err := it.Connection.Exec(SQL)
	affected := it.Connection.RowsAffected()

	return affected, err
}

func (it *SQLiteCollection) DeleteById(id string) error {
	SQL := "DELETE FROM " + it.TableName + " WHERE _id = " + id
	return it.Connection.Exec(SQL)
}


func (it *SQLiteCollection) AddFilter(ColumnName string, Operator string, Value string) error {
	if it.HasColumn(ColumnName) {
		Operator = strings.ToUpper(Operator)
		if Operator == "" || Operator == "=" || Operator == "<>" || Operator == ">" || Operator == "<" || Operator == "LIKE" {
			it.Filters = append(it.Filters, ColumnName + " " + Operator + " " + Value)
		} else {
			return errors.New("unknown operator '" + Operator + "' supposed  '', '=', '>', '<', '<>', 'LIKE' " + ColumnName + "'")
		}
	} else {
		return errors.New("can't find column '" + ColumnName + "'")
	}

	return nil
}

func (it *SQLiteCollection) ClearFilters() error {
	it.Filters = make([]string, 0)
	return nil
}

func (it *SQLiteCollection) AddSort(ColumnName string, Desc bool) error {
	if it.HasColumn(ColumnName) {
		if Desc {
			it.Order = append(it.Order, ColumnName + " DESC")
		} else {
			it.Order = append(it.Order, ColumnName)
		}
	} else {
		return errors.New("can't find column '" + ColumnName + "'")
	}

	return nil
}

func (it *SQLiteCollection) ClearSort() error {
	it.Order = make([]string, 0)
	return nil
}



func (it *SQLiteCollection) SetLimit(Offset int, Limit int) error {
	if Limit == 0 {
		it.Limit = ""
	} else {
		it.Limit = " LIMIT " + strconv.Itoa(Limit) + " OFFSET " + strconv.Itoa(Offset)
	}
	return nil
}


// Collection columns stuff
//--------------------------

func (it *SQLiteCollection) RefreshColumns() {
	SQL := "PRAGMA table_info(" + it.TableName + ")"

	row := make(sqlite3.RowMap)
	for stmt, err := it.Connection.Query(SQL); err == nil; err = stmt.Next() {
		stmt.Scan(row)

		  key := row["name"].(string)
		value := row["type"].(string)
		it.Columns[key] = value
	}
}

func (it *SQLiteCollection) ListColumns() map[string]string {
	it.RefreshColumns()
	return it.Columns
}

func (it *SQLiteCollection) HasColumn(ColumnName string) bool {
	if _, present := it.Columns[ColumnName]; present {
		return true
	} else {
		it.RefreshColumns()
		_, present := it.Columns[ColumnName]
		return present
	}
}

func (it *SQLiteCollection) AddColumn(ColumnName string, ColumnType string, indexed bool) error {

	// TODO: there probably need column name check to be only lowercase, exclude some chars, etc.

	if it.HasColumn(ColumnName) {
		return errors.New("column '" + ColumnName + "' already exists for '" + it.TableName + "' collection")
	}

	if ColumnType, err := GetDBType(ColumnType); err == nil {

		SQL := "ALTER TABLE " + it.TableName + " ADD COLUMN \"" + ColumnName + "\" " + ColumnType
		if err := it.Connection.Exec(SQL); err == nil {
			return nil
		} else {
			return sqlError(SQL, err)
		}

	}else {
		return err
	}

}

func (it *SQLiteCollection) RemoveColumn(ColumnName string) error {

	if !it.HasColumn(ColumnName) {
		return errors.New("column '" + ColumnName + "' not exists in '" + it.TableName + "' collection")
	}

	it.Connection.Begin()
	defer it.Connection.Commit()

	var SQL string

	SQL = "SELECT sql FROM sqlite_master WHERE tbl_name='" + it.TableName + "' AND type='table'"
	if stmt, err := it.Connection.Query(SQL); err == nil {

		var tableCreateSQL string = ""

		if err := stmt.Scan(&tableCreateSQL); err == nil {

			tableColumnsWTypes := ""
			tableColumnsWoTypes := ""

			re := regexp.MustCompile("CREATE TABLE .*\\((.*)\\)")
			if regexMatch := re.FindStringSubmatch(tableCreateSQL); len(regexMatch) >=2 {
				tableColumnsList := strings.Split(regexMatch[1], ", ")


				for _, tableColumn := range tableColumnsList {
					tableColumn = strings.Trim(tableColumn, "\n\t ")
					if !strings.HasPrefix(tableColumn, ColumnName) {
						if tableColumnsWTypes != "" {
							tableColumnsWTypes += ", "
							tableColumnsWoTypes += ", "
						}
						tableColumnsWTypes += "\"" + tableColumn + "\""
						tableColumnsWoTypes += "\"" + tableColumn[0 : strings.Index(tableColumn, " ")] + "\""
					}

				}
			} else {
				return errors.New("can't find table create columns in '" + tableCreateSQL + "', found [" + strings.Join(regexMatch, ", ") + "]")
			}

			SQL = "CREATE TABLE " + it.TableName + "_removecolumn (" + tableColumnsWTypes + ") "
			if err := it.Connection.Exec(SQL); err != nil {
				return sqlError(SQL, err)
			}

			SQL = "INSERT INTO " + it.TableName + "_removecolumn (" + tableColumnsWoTypes + ") SELECT " + tableColumnsWoTypes + " FROM " + it.TableName
			if err := it.Connection.Exec(SQL); err != nil {
				return sqlError(SQL, err)
			}

			// SQL = "DROP TABLE " + it.TableName
			SQL = "ALTER TABLE " + it.TableName + " RENAME TO " + it.TableName + "_fordelete"
			if err := it.Connection.Exec(SQL); err != nil {
				return sqlError(SQL, err)
			}

			SQL = "ALTER TABLE " + it.TableName + "_removecolumn RENAME TO " + it.TableName
			if err := it.Connection.Exec(SQL); err != nil {
				return sqlError(SQL, err)
			}

			it.Connection.Commit()

			// TODO: Fix this issue with table lock on DROP table

			// SQL = "DROP TABLE " + it.TableName + "_fordelete"
			// if err := it.Connection.Exec(SQL); err != nil {
			// 	return sqlError(SQL, err)
			// }

		}else{
			return err
		}

	} else {
		return sqlError(SQL, err)
	}

	return nil
}
