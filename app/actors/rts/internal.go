package rts

import (
	"encoding/json"
	"fmt"
	"github.com/ottemo/foundation/app"
	"github.com/ottemo/foundation/app/models/order"
	"github.com/ottemo/foundation/app/models/product"
	"github.com/ottemo/foundation/db"
	"github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/utils"
	"regexp"
	"strings"
	"time"
)

func GetReferrer(url string) (string, error) {
	excludeUrls := []string{app.GetFoundationUrl(""), app.GetDashboardUrl("")}

	r := regexp.MustCompile(`^(http|https):\/\/(.+)\/.*$`)
	groups := r.FindStringSubmatch(url)
	if len(groups) == 0 {
		return "", env.ErrorNew("Invalid URL in referrer")
	}
	result := groups[2]

	for index := 0; index < len(excludeUrls); index += 1 {
		if strings.Contains(excludeUrls[index], result) {
			return "", env.ErrorNew("Invalid URL in referrer")
		}
	}

	return result, nil
}

func IncreaseOnline(typeCounter int) {
	switch typeCounter {
	case REFERRER_TYPE_DIRECT:
		OnlineDirect += 1
		if OnlineDirect > OnlineDirectMax {
			OnlineDirectMax = OnlineDirect
		}
		break
	case REFERRER_TYPE_SEARCH:
		OnlineSearch += 1
		if OnlineSearch > OnlineSearchMax {
			OnlineSearchMax = OnlineSearch
		}
		break
	case REFERRER_TYPE_SITE:
		OnlineSite += 1
		if OnlineSite > OnlineSiteMax {
			OnlineSiteMax = OnlineSite
		}
		break
	}
}

func DecreaseOnline(typeCounter int) {
	switch typeCounter {
	case REFERRER_TYPE_DIRECT:
		if OnlineDirect != 0 {
			OnlineDirect -= 1
		}
		break
	case REFERRER_TYPE_SEARCH:
		if OnlineSearch != 0 {
			OnlineSearch -= 1
		}

		break
	case REFERRER_TYPE_SITE:
		if OnlineSite != 0 {
			OnlineSite -= 1
		}
		break
	}
}

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

func GetTodayVisitorsData() error {
	year := time.Now().Year()
	month := time.Now().Month()
	day := time.Now().Day()
	today := time.Date(year, month, day, 0, 0, 0, 0, time.UTC)

	if visitorsInfoToday.Day == today.In(time.Local) {
		return nil
	} else {
		year := time.Now().AddDate(0, 0, -1).Year()
		month := time.Now().AddDate(0, 0, -1).Month()
		day := time.Now().AddDate(0, 0, -1).Day()
		yesterday := time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
		if visitorsInfoToday.Day == yesterday {
			SaveVisitorData()
			visitorsInfoYesterday = visitorsInfoToday
		}
		visitorInfoCollection, err := db.GetCollection(COLLECTION_NAME_VISITORS)
		if err == nil {
			visitorInfoCollection.AddFilter("day", "=", today)
			dbRecord, _ := visitorInfoCollection.Load()

			if len(dbRecord) > 0 {
				visitorsInfoToday.Id = utils.InterfaceToString(dbRecord[0]["_id"])
				visitorsInfoToday.Day = utils.InterfaceToTime(dbRecord[0]["day"])
				visitorsInfoToday.Visitors = utils.InterfaceToInt(dbRecord[0]["visitors"])
				visitorsInfoToday.Cart = utils.InterfaceToInt(dbRecord[0]["cart"])
				visitorsInfoToday.Checkout = utils.InterfaceToInt(dbRecord[0]["checkout"])
				visitorsInfoToday.Sales = utils.InterfaceToInt(dbRecord[0]["sales"])
				visitorsInfoToday.Details = RtsDecodeDetails(utils.InterfaceToString(dbRecord[0]["details"]))

				return nil
			}
		}
	}

	visitorsInfoToday = new(dbVisitorRow)
	visitorsInfoToday.Id = ""
	visitorsInfoToday.Day = today
	visitorsInfoToday.Details = make(map[string]*VisitorDetail)

	return nil
}

func GetYesterdayVisitorsData() error {
	year := time.Now().AddDate(0, 0, -1).Year()
	month := time.Now().AddDate(0, 0, -1).Month()
	day := time.Now().AddDate(0, 0, -1).Day()
	yesterday := time.Date(year, month, day, 0, 0, 0, 0, time.Local)

	if visitorsInfoYesterday.Day == yesterday {
		return nil
	} else {

		visitorInfoCollection, err := db.GetCollection(COLLECTION_NAME_VISITORS)
		if err == nil {
			visitorInfoCollection.AddFilter("day", "=", yesterday)
			dbRecord, _ := visitorInfoCollection.Load()

			if len(dbRecord) > 0 {
				visitorsInfoYesterday.Id = utils.InterfaceToString(dbRecord[0]["_id"])
				visitorsInfoYesterday.Day = utils.InterfaceToTime(dbRecord[0]["day"])
				visitorsInfoYesterday.Visitors = utils.InterfaceToInt(dbRecord[0]["visitors"])
				visitorsInfoYesterday.Cart = utils.InterfaceToInt(dbRecord[0]["cart"])
				visitorsInfoYesterday.Checkout = utils.InterfaceToInt(dbRecord[0]["checkout"])
				visitorsInfoYesterday.Sales = utils.InterfaceToInt(dbRecord[0]["sales"])
				visitorsInfoYesterday.Details = RtsDecodeDetails(utils.InterfaceToString(dbRecord[0]["details"]))

				return nil
			}
		}
	}

	visitorsInfoYesterday = new(dbVisitorRow)
	visitorsInfoYesterday.Id = ""
	visitorsInfoYesterday.Day = yesterday
	visitorsInfoYesterday.Details = make(map[string]*VisitorDetail)

	return nil
}

func SaveVisitorData() error {
	visitorInfoCollection, err := db.GetCollection(COLLECTION_NAME_VISITORS)
	if err == nil {
		visitorInfoRow := make(map[string]interface{})
		if "" != visitorsInfoToday.Id {
			visitorInfoRow["_id"] = visitorsInfoToday.Id
		}

		visitorInfoRow["day"] = visitorsInfoToday.Day
		visitorInfoRow["visitors"] = visitorsInfoToday.Visitors
		visitorInfoRow["cart"] = visitorsInfoToday.Cart
		visitorInfoRow["checkout"] = visitorsInfoToday.Checkout
		visitorInfoRow["sales"] = visitorsInfoToday.Sales
		visitorInfoRow["details"] = RtsEncodeDetails(visitorsInfoToday.Details)

		_, err = visitorInfoCollection.Save(visitorInfoRow)
		if err != nil {
			return env.ErrorDispatch(err)
		}
	}

	return nil
}

func RtsEncodeDetails(details map[string]*VisitorDetail) string {
	jsonString, _ := json.Marshal(details)

	return string(jsonString)
}

func RtsDecodeDetails(detailsString string) map[string]*VisitorDetail {
	var details map[string]*VisitorDetail
	_ = json.Unmarshal([]byte(detailsString), &details)

	return details
}
