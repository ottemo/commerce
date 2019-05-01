package mssql

import (
	"crypto/rand"
	"encoding/hex"
	"strconv"
	"strings"
	"time"

	"github.com/ottemo/commerce/db"
	"github.com/ottemo/commerce/env"
	"github.com/ottemo/commerce/utils"
)

// creates SQL filter string based on ColumnName, Operator and Value parameters or returns nil
//   - this is an internal helper function for the AddFilter and the AddStaticFilter routines
func (it *DBCollection) makeSQLFilterString(ColumnName string, Operator string, Value interface{}) (string, error) {
	if !it.HasColumn(ColumnName) {
		return "", env.ErrorNew(ConstErrorModule, ConstErrorLevel, "8a113e66-b2e8-472d-a92c-476794b2d127", "can't find column '"+ColumnName+"'")
	}

	Operator = strings.ToUpper(Operator)
	allowedOperators := []string{"=", "!=", "<>", ">", ">=", "<", "<=", "LIKE", "IN"}

	if !utils.IsInListStr(Operator, allowedOperators) {
		return "", env.ErrorNew(ConstErrorModule, ConstErrorLevel, "11a51df3-83bb-4250-bff7-60e2e8bb6b49", "unknown operator '"+Operator+"' for column '"+ColumnName+"', allowed: '"+strings.Join(allowedOperators, "', ")+"'")
	}

	columnType := it.GetColumnType(ColumnName)

	// array column - this is a special case
	if strings.HasPrefix(columnType, "[]") {
		value := strings.Trim(convertValueForSQL(Value, columnType), "'")
		template := "(', ' || [" + ColumnName + "] || ',') LIKE '%, $value,%'"

		var resultItems []string
		for _, arrayItem := range strings.Split(value, ", ") {
			item := utils.InterfaceToString(arrayItem)
			resultItems = append(resultItems, strings.Replace(template, "$value", item, 1))
		}

		if len(resultItems) == 1 {
			return resultItems[0], nil
		}
		return strings.Join(resultItems, " OR "), nil
	}

	// regular columns - the default case
	switch Operator {
	case "LIKE":
		if typedValue, ok := Value.(string); ok {
			if !strings.Contains(typedValue, "%") {
				Value = "'%" + typedValue + "%'"
			} else {
				newValue := strings.Trim(Value.(string), "'")
				newValue = strings.Trim(newValue, "\"")
				Value = "'" + newValue + "'"
			}
		} else {
			Value = "''"
		}

	case "IN":
		if typedValue, ok := Value.(*DBCollection); ok {
			Value = "(" + typedValue.getSelectSQL() + ")"
		} else {
			newValue := "("
			for _, arrayItem := range utils.InterfaceToArray(Value) {
				newValue += convertValueForSQL(arrayItem, columnType) + ", "
			}
			newValue = strings.TrimRight(newValue, ", ") + ")"
			Value = newValue
		}

	default:
		Value = convertValueForSQL(Value, columnType)
	}

	if columnType == db.ConstTypeText {
		return "CONVERT(VARCHAR, [" + ColumnName + "] " + Operator + " " + utils.InterfaceToString(Value), nil
	}

	return "[" + ColumnName + "] " + Operator + " " + utils.InterfaceToString(Value), nil
}

// getSelectSQL returns the SQL select statement for the current collection
func (it *DBCollection) getSelectSQL() string {
	limit := " "
	if it.Limit > 0 {
		limit = " TOP " + utils.InterfaceToString(it.Limit) + " "
	}

	SQL := "SELECT" + limit + it.getSQLResultColumns() + " FROM " + it.Name + it.getSQLFilters() + it.getSQLOrder()
	return SQL
}

// modifyResultRow will un-serialize object values
func (it *DBCollection) modifyResultRow(row RowMap) RowMap {

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

// getSQLResultColumns will join the column results into a string
func (it *DBCollection) getSQLResultColumns() string {
	sqlColumns := "[" + strings.Join(it.ResultColumns, "], [") + "]"
	if sqlColumns == "[]" {
		sqlColumns = "*"
	}

	return sqlColumns
}

// getSQLOrder joins the order columns into one string with the preceding keyword
func (it *DBCollection) getSQLOrder() string {
	sqlOrder := strings.Join(it.Order, ", ")
	if sqlOrder != "" {
		sqlOrder = " ORDER BY " + sqlOrder
	}

	return sqlOrder
}

// getSQLFilters collects all filters into a single string
func (it *DBCollection) getSQLFilters() string {

	var collectSubfilters func(string) []string

	collectSubfilters = func(parentGroupName string) []string {
		var result []string

		for filterGroupName, filterGroup := range it.FilterGroups {
			if filterGroup.ParentGroup == parentGroupName {
				joinOperator := " AND "
				if filterGroup.OrSequence {
					joinOperator = " OR "
				}
				subFilters := collectSubfilters(filterGroupName)
				if len(subFilters) > 0 {
					subFilters = append([]string{""}, subFilters...)
				}
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

// getFilterGroup returns the filter group, it will create a new one if one does not already exist
func (it *DBCollection) getFilterGroup(groupName string) *StructDBFilterGroup {
	filterGroup, present := it.FilterGroups[groupName]
	if !present {
		filterGroup = &StructDBFilterGroup{Name: groupName, FilterValues: make([]string, 0)}
		it.FilterGroups[groupName] = filterGroup
	}
	return filterGroup
}

// updateFilterGroup adds a filter(combination of [column, operator, value]) in the named filter group
func (it *DBCollection) updateFilterGroup(groupName string, columnName string, operator string, value interface{}) error {

	/*if !it.HasColumn(columnName) {
		return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "67fa360b-ba93-407c-a787-1c013ebb8947", "column not found " + columnName)
	}*/

	newValue, err := it.makeSQLFilterString(columnName, operator, value)
	if err != nil {
		return err
	}

	filterGroup := it.getFilterGroup(groupName)
	filterGroup.FilterValues = append(filterGroup.FilterValues, newValue)

	return nil
}

// makeUUID generates a new UUID for _id column
func (it *DBCollection) makeUUID(id string) string {

	if len(id) != 24 {
		timeStamp := strconv.FormatInt(time.Now().Unix(), 16)

		randomBytes := make([]byte, 8)
		if _, err := rand.Reader.Read(randomBytes); err != nil {
			_ = env.ErrorDispatch(err)
		}

		randomHex := make([]byte, 16)
		hex.Encode(randomHex, randomBytes)

		id = timeStamp + string(randomHex)
	}

	return id
}
