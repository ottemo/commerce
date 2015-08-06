package quickbooks

import (
	"encoding/csv"
	"strings"
	"time"

	"github.com/ottemo/foundation/api"
	"github.com/ottemo/foundation/app/models/order"
	"github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/utils"
)

// setupAPI configures the package related API endpoints
func setupAPI() error {
	var err error

	err = api.GetRestService().RegisterAPI("quickbooks/export", api.ConstRESTOperationGet, APIExportOrders)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	return nil
}

// APIExportOrders returns a list of orders in Quickbooks IIF format
// - returns all orders in IIF format with no parameters
// - returns orders specified in url parameters
func APIExportOrders(context api.InterfaceApplicationContext) (interface{}, error) {
	var itemCSVRecords [][]string
	var requestedOrdersIDs []interface{}

	if context.GetRequestArgument("orders") != "" {
		requestedOrdersIDs = utils.InterfaceToArray(context.GetRequestArgument("orders"))
	}

	//	withCustomers := utils.InterfaceToBool(context.GetRequestArgument("customers"))

	orderCollectionModel, err := order.GetOrderCollectionModel()
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	dbOrderCollection := orderCollectionModel.GetDBCollection()
	dbOrderCollection.AddSort("created_at", false)

	if requestedOrdersIDs != nil && len(requestedOrdersIDs) > 0 {
		dbOrderCollection.AddFilter("_id", "in", requestedOrdersIDs)
	}

	dbRecords, err := dbOrderCollection.Load()
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	for _, columnsHeaders := range orderFields {
		itemCSVRecords = append(itemCSVRecords, columnsHeaders)
	}

	for _, record := range dbRecords {
		orderRecord := utils.InterfaceToMap(record)

		for _, inputValues := range dataSeted {
			var rowData []string

			for _, value := range inputValues {
				switch typedValue := value.(type) {
				case string:
					if strings.Index(typedValue, "$") == 0 {
						rowData = append(rowData, utils.InterfaceToString(orderRecord[strings.Replace(typedValue, "$", "", 1)]))
						break
					}
					rowData = append(rowData, typedValue)
					break

				case func(record map[string]interface{}) string:
					rowData = append(rowData, typedValue(orderRecord))
					break

				default:
					rowData = append(rowData, "")

				}
			}

			itemCSVRecords = append(itemCSVRecords, rowData)
		}
	}

	// preparing csv writer
	csvWriter := csv.NewWriter(context.GetResponseWriter())

	exportFilename := "orders_export_" + time.Now().Format(time.RFC3339) + ".iif"

	context.SetResponseContentType("text/csv")
	context.SetResponseSetting("Content-disposition", "attachment;filename="+exportFilename)

	for _, csvRecord := range itemCSVRecords {
		csvWriter.Write(csvRecord)
	}
	csvWriter.Flush()

	return nil, nil
}
