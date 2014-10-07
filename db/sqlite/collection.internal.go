package sqlite

import (
	"strings"

	sqlite3 "github.com/mxk/go-sqlite/sqlite3"
	"github.com/ottemo/foundation/db"
	"github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/utils"
)

// close sqlite3 statement routine
func closeStatement(statement *sqlite3.Stmt) {
	if statement != nil {
		statement.Close()
	}
}

// formats SQL query error for output to log
func sqlError(SQL string, err error) error {
	return env.ErrorNew("SQL \"" + SQL + "\" error: " + err.Error())
}

// returns string that represents value for SQL query
func convertValueForSQL(value interface{}) string {

	switch value.(type) {
	case bool:
		if value.(bool) {
			return "1"
		}
		return "0"

	case string:
		result := value.(string)
		result = strings.Replace(result, "'", "''", -1)
		result = strings.Replace(result, "\\", "\\\\", -1)
		result = "'" + result + "'"

		return result

	case map[string]interface{}, map[string]string:
		return convertValueForSQL(utils.EncodeToJsonString(value))

	case []string:
		result := ""
		for _, item := range value.([]string) {
			if result != "" {
				result += ", "
			}
			result += convertValueForSQL(item)
		}
		return result
	}

	return convertValueForSQL(utils.InterfaceToString(value))
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
	case strings.Contains(ColumnType, "date") || strings.Contains(ColumnType, "time"):
		return "NUMERIC", nil
	case ColumnType == "bool" || ColumnType == "boolean":
		return "NUMERIC", nil
	}

	return "?", env.ErrorNew("Unknown type '" + ColumnType + "'")
}

// makes SQL filter string based on ColumnName, Operator and Value parameters or returns nil
//   - internal usage function for AddFilter and AddStaticFilter routines
func (it *SQLiteCollection) makeSQLFilterString(ColumnName string, Operator string, Value interface{}) (string, error) {
	if !it.HasColumn(ColumnName) {
		return "", env.ErrorNew("can't find column '" + ColumnName + "'")
	}

	Operator = strings.ToUpper(Operator)
	allowedOperators := []string{"=", "!=", "<>", ">", "<", "LIKE", "IN"}

	if !utils.IsInListStr(Operator, allowedOperators) {
		return "", env.ErrorNew("unknown operator '" + Operator + "' for column '" + ColumnName + "', allowed: '" + strings.Join(allowedOperators, "', ") + "'")
	}

	switch Operator {
	case "LIKE":
		if typedValue, ok := Value.(string); ok && !strings.Contains(typedValue, "%") {
			Value = "%" + typedValue + "%"
		}

	case "IN":
		if typedValue, ok := Value.(SQLiteCollection); ok {
			Value = "(" + typedValue.getSelectSQL() + ")"
		} else {
			newValue := "("
			for _, arrayItem := range utils.InterfaceToArray(Value) {
				newValue += utils.InterfaceToString(arrayItem)
			}
			newValue += ")"
			Value = newValue
		}
	}

	return ColumnName + " " + Operator + " " + convertValueForSQL(Value), nil
}

func (it *SQLiteCollection) getSelectSQL() string {
	SQL := "SELECT " + it.getSQLResultColumns() + " FROM " + it.Name + it.getSQLFilters() + it.getSQLOrder() + it.Limit
	return SQL
}

// un-serialize object values
func (it *SQLiteCollection) modifyResultRow(row sqlite3.RowMap) sqlite3.RowMap {

	for columnName, columnValue := range row {
		columnType := it.GetColumnType(columnName)
		if columnType != "" {
			row[columnName] = db.ConvertTypeFromDbToGo(columnValue, columnType)
		}
	}

	if _, present := row["_id"]; present {
		row["_id"] = utils.InterfaceToString(row["_id"])
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
		return env.ErrorNew("not existing column " + columnName)
	}*/

	newValue, err := it.makeSQLFilterString(columnName, operator, value)
	if err != nil {
		return err
	}

	filterGroup := it.getFilterGroup(groupName)
	filterGroup.FilterValues = append(filterGroup.FilterValues, newValue)

	return nil
}
