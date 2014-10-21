package sqlite

import (
	"crypto/rand"
	"encoding/hex"
	"strconv"
	"strings"
	"time"

	sqlite3 "github.com/mxk/go-sqlite/sqlite3"
	"github.com/ottemo/foundation/db"
	"github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/utils"
)

// makes SQL filter string based on ColumnName, Operator and Value parameters or returns nil
//   - internal usage function for AddFilter and AddStaticFilter routines
func (it *SQLiteCollection) makeSQLFilterString(ColumnName string, Operator string, Value interface{}) (string, error) {
	if !it.HasColumn(ColumnName) {
		return "", env.ErrorNew("can't find column '" + ColumnName + "'")
	}

	Operator = strings.ToUpper(Operator)
	allowedOperators := []string{"=", "!=", "<>", ">", ">=", "<", "<=", "LIKE", "IN"}

	if !utils.IsInListStr(Operator, allowedOperators) {
		return "", env.ErrorNew("unknown operator '" + Operator + "' for column '" + ColumnName + "', allowed: '" + strings.Join(allowedOperators, "', ") + "'")
	}

	switch Operator {
	case "LIKE":
		if typedValue, ok := Value.(string); ok && !strings.Contains(typedValue, "%") {
			Value = "'%" + typedValue + "%'"
		}

	case "IN":
		if typedValue, ok := Value.(SQLiteCollection); ok {
			Value = "(" + typedValue.getSelectSQL() + ")"
		} else {
			newValue := "("
			for _, arrayItem := range utils.InterfaceToArray(Value) {
				newValue += convertValueForSQL(arrayItem) + ", "
			}
			newValue = strings.TrimRight(newValue, ", ") + ")"
			Value = newValue
		}

	default:
		Value = convertValueForSQL(Value)
	}

	return "`" + ColumnName + "` " + Operator + " " + utils.InterfaceToString(Value), nil
}

// returns SQL select statement for current collection
func (it *SQLiteCollection) getSelectSQL() string {
	SQL := "SELECT " + it.getSQLResultColumns() + " FROM " + it.Name + it.getSQLFilters() + it.getSQLOrder() + it.Limit
	return SQL
}

// un-serialize object values
func (it *SQLiteCollection) modifyResultRow(row sqlite3.RowMap) sqlite3.RowMap {

	for columnName, columnValue := range row {
		columnType, present := dbEngine.attributeTypes[it.Name][columnName]
		if !present {
			columnType = ""
		}

		if columnName != "_id" && columnType != "" {
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
	sqlColumns := "`" + strings.Join(it.ResultColumns, "`, `") + "`"
	if sqlColumns == "``" {
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

func (it *SQLiteCollection) makeUUID(id string) string {

	if len(id) != 24 {
		timeStamp := strconv.FormatInt(time.Now().Unix(), 16)

		randomBytes := make([]byte, 8)
		rand.Reader.Read(randomBytes)

		randomHex := make([]byte, 16)
		hex.Encode(randomHex, randomBytes)

		id = timeStamp + string(randomHex)
	}

	return id
}
