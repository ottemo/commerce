package mssql

import (
	"database/sql"
	"errors"
	"time"

	"github.com/ottemo/commerce/env"
	"github.com/ottemo/commerce/utils"
)

// ------------------------------------------------------------------------------------
// InterfaceDBConnector for MS SQL Server implementation based on package "github.com/ottemo/commerce/db/interfaces"
// ------------------------------------------------------------------------------------

// GetConnectionParams returns configured DB connection parameters
func (it *DBEngine) GetConnectionParams() interface{} {
	var connectionParams = connectionParamsType{
		uri:             "sqlserver://127.0.0.1:1433",
		dbName:          "commerce",
		poolConnections: 10,
		maxConnections:  0,
	}

	if iniConfig := env.GetIniConfig(); iniConfig != nil {
		if iniValue := iniConfig.GetValue("db.mssql.uri", connectionParams.uri); iniValue != "" {
			connectionParams.uri = iniValue
		}

		if iniValue := iniConfig.GetValue("db.mssql.db", connectionParams.dbName); iniValue != "" {
			connectionParams.dbName = iniValue
		}

		if iniValue := iniConfig.GetValue("db.mssql.maxConnections", ""); iniValue != "" {
			connectionParams.maxConnections = utils.InterfaceToInt(iniValue)
		}

		if iniValue := iniConfig.GetValue("db.mssql.poolConnections", ""); iniValue != "" {
			connectionParams.poolConnections = utils.InterfaceToInt(iniValue)
		}
	}

	return connectionParams
}

// Connect establishes DB connection
func (it *DBEngine) Connect(srcConnectionParams interface{}) error {
	connectionParams, ok := srcConnectionParams.(connectionParamsType)
	if !ok {
		return errors.New("invalid connection parameters type")
	}

	// performing DB connection
	newConnection, err := sql.Open("sqlserver", "sqlserver://"+connectionParams.uri)
	if err != nil {
		return err
	}
	it.connection = newConnection

	//if connectionParams.poolConnections > 0 {
	//	it.connection.SetMaxIdleConns(connectionParams.poolConnections)
	//}
	//if connectionParams.maxConnections > 0 {
	//	it.connection.SetMaxOpenConns(connectionParams.maxConnections)
	//}

	// making sure DB selected otherwise trying to obtain DB
	rows, err := it.connection.Query("SELECT DB_NAME()")
	if err != nil {
		return err
	}

	var database string
	if !rows.Next() || rows.Scan(&database) != nil || database == "" || database != connectionParams.dbName {
		if _, err := it.connection.Exec("USE " + connectionParams.dbName); err != nil {
			if _, err = it.connection.Exec("CREATE DATABASE " + connectionParams.dbName); err != nil {
				return err
			}
			if _, err = it.connection.Exec("USE " + connectionParams.dbName); err != nil {
				return err
			}
		}
		connectionParams.uri += "?database=" + connectionParams.dbName
		return it.Connect(connectionParams)
	}

	return nil
}

// AfterConnect makes initialization of DB engine
func (it *DBEngine) AfterConnect(srcConnectionParams interface{}) error {
	SQL := "IF (OBJECT_ID('" + ConstCollectionNameColumnInfo + "', 'U') IS NULL) "
	SQL += "CREATE TABLE [" + ConstCollectionNameColumnInfo + "] ( " +
		"[_id]        INTEGER NOT NULL IDENTITY(1,1)," +
		"[collection] VARCHAR(255)," +
		"[column]     VARCHAR(255)," +
		"[type]       VARCHAR(255)," +
		"[indexed]    BIT," +
		"PRIMARY KEY([_id]) )"

	_, err := it.connection.Exec(SQL)
	if err != nil {
		return sqlError(SQL, err)
	}

	return nil
}

// Reconnect tries to reconnect to DB
func (it *DBEngine) Reconnect(connectionParams interface{}) error {
	// Call Connect to "select database".
	// Simple it.Ping is not enough - it restores connection, but doesn't "select database".
	return it.Connect(connectionParams)
}

// IsConnected returns connection status
func (it *DBEngine) IsConnected() bool {
	return it.isConnected
}

// SetConnected sets connection status
func (it *DBEngine) SetConnected(connected bool) {
	it.isConnected = connected
}

// Ping checks connection alive
func (it *DBEngine) Ping() error {
	_, err := it.connection.Query("SELECT 1")
	return err
}

// GetValidationInterval returns delay between Ping
func (it *DBEngine) GetValidationInterval() time.Duration {
	return ConstConnectionValidateInterval
}

// GetEngineName returns DBEngine name (InterfaceDBConnector)
func (it *DBEngine) GetEngineName() string {
	return it.GetName()
}

// LogConnection outputs message to log
func (it *DBEngine) LogConnection(message string) {
	env.Log("mssql.log", "DEBUG", message)
}
