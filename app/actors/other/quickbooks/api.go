package quickbooks

import (
	"encoding/csv"
	"github.com/ottemo/foundation/api"
	"github.com/ottemo/foundation/app/models/order"
	"github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/utils"
	"strings"
)

// setupAPI setups package related API endpoint routines
func setupAPI() error {
	var err error

	err = api.GetRestService().RegisterAPI("quickbook/export", api.ConstRESTOperationGet, APIExportOrders)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	return nil
}

// APIGetGiftCard return gift card info buy it's code
func APIExportOrders(context api.InterfaceApplicationContext) (interface{}, error) {
	var itemCSVRecords [][]string

	orderCollectionModel, err := order.GetOrderCollectionModel()
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	dbOrderCollection := orderCollectionModel.GetDBCollection()
	dbOrderCollection.AddSort("created_at", false)

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

	exportFilename := "ordersTestExport.csv"

	context.SetResponseContentType("text/csv")
	context.SetResponseSetting("Content-disposition", "attachment;filename="+exportFilename)

	for _, csvRecord := range itemCSVRecords {
		csvWriter.Write(csvRecord)
	}
	csvWriter.Flush()

	return nil, nil
}

//func convertRecordToRow (record map[string]interface {}, key string) []string {
//	result := make([]string, 0)
//
//	return ""
//}
