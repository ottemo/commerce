package rts

import (
	"fmt"
	"github.com/ottemo/foundation/app/models/order"
	"github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/utils"
	"time"
	"github.com/ottemo/foundation/db"
	"github.com/ottemo/foundation/app/models/product"
)

func GetProducts() ([]map[string]interface{}, error) {
	productCollectionModel, err := product.GetProductCollectionModel()
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}
	dbProductCollection := productCollectionModel.GetDBCollection()
	dbProductCollection.SetResultColumns("_id")
	return dbProductCollection.Load()
}

func GetDateFrom() (time.Time, error) {
	result := time.Now()

	salesHistoryCollection, err := db.GetCollection(COLLECTION_NAME_SALES_HISTORY)
	if err == nil {
		salesHistoryCollection.SetResultColumns("created_at")
		salesHistoryCollection.AddSort("created_at", true)
		salesHistoryCollection.SetLimit(0, 1)
		dbRecord, _ := salesHistoryCollection.Load()

		if len(dbRecord) > 0 {
			datetime := utils.InterfaceToTime(dbRecord[0]["created_at"])
			year := datetime.Year()
			month := datetime.Month()
			day := datetime.Day()
			result := time.Date(year, month, day, 0, 0, 0, 0, time.Local)
			if time.Now().Format("2006-01-02") == datetime.Format("2006-01-02") {
				return result, env.ErrorNew("Sales history has last data")
			}

			return result, nil
		}
	}

	orderCollectionModel, err := order.GetOrderCollectionModel()
	dateFrom := time.Now()
	if err != nil {
		year := dateFrom.Year()
		month := dateFrom.Month()
		day := dateFrom.Day()
		result := time.Date(year, month, day, 0, 0, 0, 0, time.Local)
		return result, nil
	}
	dbOrderCollection := orderCollectionModel.GetDBCollection()
	dbOrderCollection.SetResultColumns("created_at")
	dbOrderCollection.AddSort("created_at", false)
	dbRecord, err := dbOrderCollection.Load()
	if len(dbRecord) > 0 {
		dateFrom = utils.InterfaceToTime(dbRecord[0]["created_at"])
		year := dateFrom.Year()
		month := dateFrom.Month()
		day := dateFrom.Day()
		result := time.Date(year, month, day, 0, 0, 0, 0, time.Local)
		return result, nil
	}

	return result, nil
}

func GetOrderItems(date time.Time, productId string) ([]map[string]interface{}, error) {

	orderItemCollectionModel, err := order.GetOrderItemCollectionModel()
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	dbOrderItemCollection := orderItemCollectionModel.GetDBCollection()
	dbOrderItemCollection.SetResultColumns("qty")

	orderCollectionModel, err := order.GetOrderCollectionModel()
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}
	dbOrderCollection := orderCollectionModel.GetDBCollection()
	dbOrderCollection.SetResultColumns("_id")

	year := date.Year()
	month := date.Month()
	day := date.Day()
	todayFrom := time.Date(year, month, day, 0, 0, 0, 0, time.Local)
	todayTo := time.Date(year, month, day, 23, 59, 59, 0, time.Local)
	dbOrderCollection.AddFilter("created_at", ">=", todayFrom)
	dbOrderCollection.AddFilter("created_at", "<=", todayTo)

	dbOrderItemCollection.AddFilter("order_id", "IN", dbOrderCollection)
	dbOrderItemCollection.AddFilter("product_id", "=", productId)

	return dbOrderItemCollection.Load()
}

func DeleteExistingRowHistory(date time.Time, productId string) error {
	salesHistoryCollection, err := db.GetCollection(COLLECTION_NAME_SALES_HISTORY)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	salesHistoryCollection.ClearFilters()
	salesHistoryCollection.AddFilter("product_id", "=", productId)
	salesHistoryCollection.AddFilter("created_at", "=", date)
	dbSalesHist, _ := salesHistoryCollection.Load()
	if len(dbSalesHist) > 0 {
		err = salesHistoryCollection.DeleteById(utils.InterfaceToString(dbSalesHist[0]["_id"]))
		if err != nil {
			return env.ErrorDispatch(err)
		}
	}

	return nil
}

func initSalesHistory() error {
	sales := make(map[string]int)
	dateTo := time.Now()

	dbProductRecord, _ := GetProducts()
	dateFrom, err := GetDateFrom()
	if err != nil {
		return nil
	}

	salesHistoryCollection, err := db.GetCollection(COLLECTION_NAME_SALES_HISTORY)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	for date := dateFrom; int32(date.Unix()) < int32(dateTo.Unix()); date = date.AddDate(0, 0, 1) {
		for _, productItem := range dbProductRecord {
			productId := utils.InterfaceToString(productItem["_id"])
			DeleteExistingRowHistory(date, productId)
			count := 0

			items, _ := GetOrderItems(date, productId)
			for _, item := range items {
				count += utils.InterfaceToInt(item["qty"])
			}

			// Add history row
			salesHistoryRow := make(map[string]interface{})
			salesHistoryRow["product_id"] = productId
			salesHistoryRow["created_at"] = date
			salesHistoryRow["count"] = count
			_, err = salesHistoryCollection.Save(salesHistoryRow)
			if err != nil {
				return env.ErrorDispatch(err)
			}
			sales[productId] += count
		}
		SaveSalesData(sales)
		sales = make(map[string]int)
	}

	return nil
}


func SaveSalesData(data map[string]int) error {

	if len(data) == 0 {
		return nil
	}

	salesCollection, err := db.GetCollection(COLLECTION_NAME_SALES)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	for productId, count := range data {
		// Add history row
		salesRow := make(map[string]interface{})

		salesCollection.ClearFilters()
		salesCollection.AddFilter("range", "=", GetSalesRange())
		salesCollection.AddFilter("product_id", "=", productId)
		dbSaleRow, _ := salesCollection.Load()
		if len(dbSaleRow) > 0 {
			salesRow["_id"] = utils.InterfaceToString(dbSaleRow[0]["_id"])
			oldCount := utils.InterfaceToInt(dbSaleRow[0]["count"])
			count += oldCount
		}

		salesRow["product_id"] = productId
		salesRow["count"] = count
		salesRow["range"] = GetSalesRange()
		_, err = salesCollection.Save(salesRow)
		if err != nil {
			return env.ErrorDispatch(err)
		}
	}

	return nil
}

func GetSalesRange() string {
	_range := "2014-01-01:"

	return _range
}

// DB preparations for current model implementation
func GetTotalSales(fromDate time.Time, toDate time.Time) error {

	orderCollectionModelT, err := order.GetOrderCollectionModel()
	if err != nil {
		return env.ErrorDispatch(err)
	}

	dbCollection := orderCollectionModelT.GetDBCollection()

	year := fromDate.Year()
	month := fromDate.Month()
	day := fromDate.Day()
	todayFrom := time.Date(year, month, day, 0, 0, 0, 0, time.Local)
	todayTo := time.Date(year, month, day, 23, 59, 59, 0, time.Local)

	dbCollection.AddFilter("created_at", ">=", todayFrom)
	dbCollection.AddFilter("created_at", "<=", todayTo)

	// filters handle for today
	today, err := dbCollection.Count()
	if err != nil {
		return env.ErrorDispatch(err)
	}

	dbCollection.ClearFilters()
	year = toDate.Year()
	month = toDate.Month()
	day = toDate.Day()
	yesterdayFrom := time.Date(year, month, day, 0, 0, 0, 0, time.Local)
	yesterdayTo := time.Date(year, month, day, 23, 59, 59, 0, time.Local)

	dbCollection.AddFilter("created_at", ">=", yesterdayFrom)
	dbCollection.AddFilter("created_at", "<=", yesterdayTo)

	// filters handle for yesterday
	yesterday, err := dbCollection.Count()
	if err != nil {
		return env.ErrorDispatch(err)
	}

	ratio := float64(1)
	if 0 != yesterday {
		ratio = float64(today)/float64(yesterday)-float64(1)
	}

	sales.ratio = ratio
	sales.today = today
	sales.lastUpdate = time.Now().Unix()
	sales.yesterday = yesterday

	return nil
}

func GetSalesDetail(fromDate time.Time, toDate time.Time, hash string) error {

	orderCollectionModelT, err := order.GetOrderCollectionModel()
	if err != nil {
		return env.ErrorDispatch(err)
	}

	dbCollection := orderCollectionModelT.GetDBCollection()
	dbCollection.SetResultColumns("_id", "created_at")
	dbCollection.AddSort("created_at", false)

	year := fromDate.Year()
	month := fromDate.Month()
	day := fromDate.Day()
	dateFrom := time.Date(year, month, day, 0, 0, 0, 0, time.Local)
	dbCollection.AddFilter("created_at", ">=", dateFrom)

	year = toDate.Year()
	month = toDate.Month()
	day = toDate.Day()
	dateTo := time.Date(year, month, day, 0, 0, 0, 0, time.Local)
	dbCollection.AddFilter("created_at", "<=", dateTo)

	// filters handle for yesterday
	list, err := dbCollection.Load()
	if err != nil {
		return env.ErrorDispatch(err)
	}

	delta := toDate.Sub(fromDate)
	salesDetail[hash] = &SalesDetailData{Data: make(map[string]int)}
	if delta.Hours() > 48 { // group by days
		// fills the data a zero
		for date := fromDate; int32(date.Unix()) < int32(toDate.Unix()); date = date.AddDate(0, 0, 1) {
			timestamp := int64(date.Unix())
			mapIndex := GetDayForTimestamp(timestamp, false)

			salesDetail[hash].Data[mapIndex] = 0
		}
		// counts items
		for _, order := range list {
			timestamp := int64(utils.InterfaceToTime(order["created_at"]).Unix())
			mapIndex := GetDayForTimestamp(timestamp, false)

			salesDetail[hash].Data[mapIndex] += 1
		}
	} else { // group by hours

		for date := fromDate; int32(date.Unix()) < int32(toDate.Unix()); date = date.AddDate(0, 0, 1) {
			timestamp := int64(date.Unix())
			year := time.Unix(int64(timestamp), 0).Year()
			month := time.Unix(int64(timestamp), 0).Month()
			day := time.Unix(int64(timestamp), 0).Day()
			for hour := 0; hour < 24; hour += 1 {
				timeGroup := time.Date(year, month, day, hour, 0, 0, 0, time.Local)
				if timeGroup.Unix() > time.Now().Unix() {
					break
				}
				mapIndex := fmt.Sprintf("%v", int32(timeGroup.Unix()))
				salesDetail[hash].Data[mapIndex] = 0
			}
		}
		for _, order := range list {
			timestamp := int64(utils.InterfaceToTime(order["created_at"]).Unix())
			mapIndex := GetDayForTimestamp(timestamp, true)

			salesDetail[hash].Data[mapIndex] += 1
		}
	}

	salesDetail[hash].lastUpdate = time.Now().Unix()

	return nil
}

// Rounds time to day or hour
func GetDayForTimestamp(timestamp int64, byHour bool) string {
	hour := 0
	if byHour {
		hour = time.Unix(timestamp, 0).Hour()
	}

	year := time.Unix(timestamp, 0).Year()
	month := time.Unix(timestamp, 0).Month()
	day := time.Unix(timestamp, 0).Day()
	timeGroup := time.Date(year, month, day, hour, 0, 0, 0, time.Local)
	mapIndex := fmt.Sprintf("%v", int32(timeGroup.Unix()))

	return mapIndex
}
