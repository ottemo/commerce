package rts

import (
	"github.com/ottemo/foundation/app/models/order"
	"time"
	"github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/utils"
)

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

	if delta.Hours() > 48 {
//		// group by days
		for _, order := range list {
			timestamp := utils.InterfaceToTime(order["created_at"])
//			year := time.Unix(int64(timestamp), 0).Year()
//			month := time.Unix(int64(timestamp), 0).Month()
//			day := time.Unix(int64(timestamp), 0).Day()
//			timeGroup := time.Date(year, month, day, 0, 0, 0, 0, time.Local)
//
//			mapIndex := fmt.Sprintf("%v", int32(timeGroup.Unix()))
			println(utils.InterfaceToTime(timestamp))
//			salesDetail[hash].Data[mapIndex] += 1
		}
//	} else {
//		// group by days
//		for date := fromDate; int32(date.Unix()) < int32(toDate.Unix()); date = date.AddDate(0, 0, 1) {
//			for _, timestamp := range visits.Data[date.Format("2006-01-02")] {
//				hour := time.Unix(int64(timestamp), 0).Hour()
//				year := time.Unix(int64(timestamp), 0).Year()
//				month := time.Unix(int64(timestamp), 0).Month()
//				day := time.Unix(int64(timestamp), 0).Day()
//				timeGroup := time.Date(year, month, day, hour, 0, 0, 0, time.Local)
//
//				mapIndex := fmt.Sprintf("%v", int32(timeGroup.Unix()))
//				result[mapIndex] += 1
//			}
//		}
	}

	salesDetail[hash].lastUpdate = time.Now().Unix()

	return nil
}
