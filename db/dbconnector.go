package db

import (
	"fmt"
	"time"

	"github.com/ottemo/commerce/utils"
)

// Connection establishes a connection to the db and keeps it alive
type Connection struct {
	connector InterfaceDBConnector
}

// NewConnection returns new Connection instance
func NewConnection(connector InterfaceDBConnector) *Connection {
	return &Connection{connector: connector}
}

// Async instantiates the connection process asynchronously
func (it *Connection) Async() error {
	go func() {
		if err := it.Connect(); err != nil {
			it.log("Internal error: unable to start " + it.connector.GetEngineName() + " DBEngine: " + err.Error())
		}
	}()
	return nil
}

// Connect implements connection process algorithm, however the process will quit if a connection cannot be
// created. To not crash the app when database connection failures happen, use the wrapper
// db.Connection.Async()
func (it *Connection) Connect() error {
	var connectionParams = it.connector.GetConnectionParams()

	it.log("Connecting to " + it.connector.GetEngineName() + " DB. Timeout: 10 seconds.")

	var ticker = time.NewTicker(it.connector.GetValidationInterval())
	var isConnectMessageShown = false
	for !it.connector.IsConnected() {
		err := it.connector.Connect(connectionParams)

		if err != nil {
			if !isConnectMessageShown {
				it.log("Can't connect to DBEngine: " + err.Error())
				isConnectMessageShown = true
			}

			it.log("Wait " + utils.InterfaceToString(10) + " seconds to reconnect.")
			_ = <-ticker.C
		} else {
			it.connector.SetConnected(true)
			it.log("DB connection established.")
		}
	}

	err := it.connector.AfterConnect(connectionParams)
	if err != nil {
		it.log("After DB connect error:")
	}

	go func() {
		for range ticker.C {
			if err := it.connector.Ping(); err != nil {
				it.connector.SetConnected(false)
				it.log("DB connection lost. Reconnect in 60 seconds.")
				if err := it.connector.Reconnect(connectionParams); err != nil {
					it.log("Unable to reconnect: " + err.Error())
				}
			} else if !it.connector.IsConnected() { // Show message only once
				it.connector.SetConnected(true)
				it.log("DB connection restored.")
			}
		}
	}()

	err = OnDatabaseStart()
	if err != nil {
		return err
	}

	err = OnDatabaseStart()
	return err
}

// log outputs messages to stdout and the connection endpoint
func (it *Connection) log(message string) {
	// output to stdout
	fmt.Println(time.Now().Format(time.RFC3339), message)
	it.connector.LogConnection(message)
}
