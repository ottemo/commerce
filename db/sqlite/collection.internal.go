package sqlite

import (
	"errors"
	"strings"

	sqlite3 "github.com/mxk/go-sqlite/sqlite3"
	"strconv"
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

// converts _id field from int for string
func (it *SQLiteCollection) modifyResultRowId(row sqlite3.RowMap) sqlite3.RowMap {
	if _, present := row["_id"]; present {
		row["_id"] = strconv.FormatInt(row["_id"].(int64), 10)
	}

	return row
}

// joins result columns in string
func (it *SQLiteCollection) getSQLResultColumns() string {
	sqlColumns := strings.Join(it.ResultColumns, ", ")
	if sqlColumns == "" {
		sqlColumns = "*"
	}

	return sqlColumns
}

// joins order olumns in one string with preceding keyword
func (it *SQLiteCollection) getSQLOrder() string {
	sqlOrder := strings.Join(it.Order, ", ")
	if sqlOrder != "" {
		sqlOrder = " ORDER BY " + sqlOrder
	}

	return sqlOrder
}

// collects all filters in a single string (for internal usage)
func (it *SQLiteCollection) getSQLFilters() string {

	var collectSubfilters func(string) []string = nil

	collectSubfilters = func(parentGroupName string) []string {
		result := make([]string, 0)

		for filterGroupName, filterGroup := range it.FilterGroups {
			if filterGroup.ParentGroup == parentGroupName {
				joinOperator := " AND "
				if filterGroup.OrSequence {
					joinOperator = " OR "
				}
				subFilters := collectSubfilters(filterGroupName)
				filterValue := "(" + strings.Join(filterGroup.FilterValues, joinOperator) + strings.Join(subFilters, joinOperator) + ")"
				result = append(result, filterValue)
			}
		}

		return result
	}

	sqlFilters := strings.Join(collectSubfilters(""), " AND ")
	if sqlFilters != "" {
		sqlFilters = " WHERE " + sqlFilters
	}

	return sqlFilters
}

// returns filter group, creates new one if not exists
func (it *SQLiteCollection) getFilterGroup(groupName string) *T_DBFilterGroup {
	filterGroup, present := it.FilterGroups[groupName]
	if !present {
		filterGroup = &T_DBFilterGroup{Name: groupName, FilterValues: make([]string, 0)}
		it.FilterGroups[groupName] = filterGroup
	}
	return filterGroup
}

// adds filter(combination of [column, operator, value]) in named filter group
func (it *SQLiteCollection) updateFilterGroup(groupName string, columnName string, operator string, value interface{}) error {

	/*if !it.HasColumn(columnName) {
		return errors.New("not existing column " + columnName)
	}*/

	newValue, err := it.makeSQLFilterString(columnName, operator, value)
	if err != nil {
		return err
	}

	filterGroup := it.getFilterGroup(groupName)
	filterGroup.FilterValues = append(filterGroup.FilterValues, newValue)

	return nil
}
