package mssql

import (
	"fmt"
	"strings"
	"time"

	"database/sql"

	"github.com/ottemo/commerce/db"
	"github.com/ottemo/commerce/env"
	"github.com/ottemo/commerce/utils"
)

// exec routines
func connectionExecWLastInsertID(SQL string, args ...interface{}) (int64, error) {

	var lastInsertID int64

	rows, err := dbEngine.connection.Query(SQL + "; SELECT SCOPE_IDENTITY()")
	if err != nil {
		return lastInsertID, err
	}

	if rows.Next() {
		err = rows.Scan(&lastInsertID)
	}

	return lastInsertID, err
}

// exec routines
func connectionExecWAffected(SQL string, args ...interface{}) (int64, error) {

	if ConstDebugSQL {
		env.Log(ConstDebugFile, env.ConstLogPrefixInfo, SQL)
	}

	result, err := dbEngine.connection.Exec(SQL, args...)
	if err != nil {
		return 0, err
	}

	return result.RowsAffected()
}

// exec routines
func connectionExec(SQL string, args ...interface{}) error {

	if ConstDebugSQL {
		env.Log(ConstDebugFile, env.ConstLogPrefixInfo, SQL)
	}

	_, err := dbEngine.connection.Exec(SQL, args...)

	return err
}

// query routines
func connectionQuery(SQL string) (*sql.Rows, error) {
	if ConstDebugSQL {
		env.Log(ConstDebugFile, env.ConstLogPrefixInfo, SQL)
	}

	return dbEngine.connection.Query(SQL)
}

// closeCursor closes cursor statement routine
func closeCursor(cursor *sql.Rows) {
	if cursor != nil {
		if err := cursor.Close(); err != nil {
			_ = env.ErrorDispatch(err)
		}
	}
}

// formats SQL query error for output to log
func sqlError(SQL string, err error) error {
	return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "8c3f2f99-4d08-412b-9dd4-fb6834c44c2b", "SQL \""+SQL+"\" error: "+err.Error())
}

// returns string that represents value for SQL query
func convertValueForSQL(value interface{}, dbType string) string {

	var sqlValue string

	switch value.(type) {
	case *DBCollection:
		sqlValue = value.(*DBCollection).getSelectSQL()

	case bool:
		if value.(bool) {
			sqlValue = "1"
		} else {
			sqlValue = "0"
		}

	case string:
		result := value.(string)
		result = strings.Replace(result, "'", "''", -1)
		result = strings.Replace(result, "\\", "\\\\", -1)
		result = "'" + result + "'"

		sqlValue = result

	case int, int32, int64:
		sqlValue = utils.InterfaceToString(value)

	case map[string]interface{}, map[string]string:
		sqlValue = convertValueForSQL(utils.EncodeToJSONString(value), dbType)

	case []string, []int, []int64, []int32, []float64, []bool:
		sqlValue = convertValueForSQL(utils.InterfaceToArray(value), dbType)

	case time.Time:
		sqlValue = convertValueForSQL(value.(time.Time).Unix(), dbType)

	case []interface{}:
		result := ""
		for _, item := range value.([]interface{}) {
			if result != "" {
				result += ", "
			}
			result += strings.Replace(utils.InterfaceToString(item), ",", "#2C;", -1)
		}
		sqlValue = convertValueForSQL(result, dbType)
	default:
		sqlValue = convertValueForSQL(utils.InterfaceToString(value), dbType)
	}

	if dbType == db.ConstTypeID && !ConstUseUUIDids {
		sqlValue = utils.InterfaceToString(utils.InterfaceToInt(value))
	} else if db.TypeIsString(dbType) && sqlValue != "" && sqlValue[0:1] != "'" {
		sqlValue = "'" + sqlValue + "'"
	} else if dbType == db.ConstTypeDatetime {
		// TODO: convert datetime
		//if value, err := strconv.ParseInt("1405544146", 10, 64); err == nil {
		//	sqlValue = "'" + strings.Replace(time.Unix(value, 0).Format("2006-01-02 15:04:05"), "T", " ", -1) + "'"
		//}
	}

	return sqlValue
}

func getRowAsStringMap(rows *sql.Rows) (RowMap, error) {
	row := make(RowMap)

	columns, err := rows.Columns()
	if err != nil {
		return row, env.ErrorDispatch(err)
	}

	values := make([]sql.NullString, len(columns))

	scanArgs := make([]interface{}, len(values))
	for i := range values {
		scanArgs[i] = &values[i]
	}

	err = rows.Scan(scanArgs...)
	if err != nil {
		if rows.Next() {
			err = rows.Scan(scanArgs...)
		}

		if err != nil {
			return row, env.ErrorDispatch(err)
		}
	}

	for idx, column := range columns {
		if values[idx].Valid {
			row[column] = values[idx].String
		} else {
			row[column] = nil
		}
	}

	return row, nil
}

// GetDBType returns type used inside mssql for given general name
func GetDBType(ColumnType string) (string, error) {
	ColumnType = strings.ToLower(ColumnType)
	switch {
	case strings.HasPrefix(ColumnType, "[]"):
		return "TEXT", nil
	case ColumnType == db.ConstTypeID:
		if ConstUseUUIDids {
			return "VARCHAR(24)", nil
		}
		return "BIGINT", nil
	case ColumnType == "int" || ColumnType == "integer":
		return "INT", nil
	case ColumnType == "real" || ColumnType == "float":
		return "DECIMAL(10,10)", nil
	case strings.Contains(ColumnType, "char") || ColumnType == "string":
		dataType := utils.DataTypeParse(ColumnType)
		if dataType.Precision > 0 {
			return fmt.Sprintf("VARCHAR(%d)", dataType.Precision), nil
		}
		return "VARCHAR(255)", nil
	case ColumnType == "text" || ColumnType == "json":
		return "TEXT", nil
	case ColumnType == "blob" || ColumnType == "struct" || ColumnType == "data":
		return "BLOB", nil
	case strings.Contains(ColumnType, "numeric") || strings.Contains(ColumnType, "decimal"):
		return "DECIMAL(10,5)", nil
	case ColumnType == "money":
		return "DECIMAL(10,2)", nil
	case strings.Contains(ColumnType, "date") || strings.Contains(ColumnType, "time"):
		return "INT", nil
		// TODO: convert to DATETIME
		// return "DATETIME", nil
	case ColumnType == "bool" || ColumnType == "boolean":
		return "BIT", nil
	}

	return "?", env.ErrorNew(ConstErrorModule, ConstErrorLevel, "80757774-5967-4e47-8429-7ca2cbcea72c", "Unknown type '"+ColumnType+"'")
}
