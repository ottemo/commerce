package rts

import (
	"github.com/ottemo/foundation/utils"
	"regexp"
	"time"
	"github.com/ottemo/foundation/db"
	"fmt"
)

func referrerHandler(event string, data map[string]interface{}) bool {

	if "api.referer" != event || "" == utils.InterfaceToString(data["referer"]) {
		return true
	}

	str := utils.InterfaceToString(data["referer"])

	r := regexp.MustCompile(`(^(http|https):\/\/.+\/).*$`)
	groups := r.FindStringSubmatch(str)
	if len(groups) == 0 {
		return true
	}

	referrer := groups[1]
	sessionId := utils.InterfaceToString(data["sessionId"])

	if _, ok := referrers[referrer]; !ok {
		referrers[referrer] = &ReferrerData{Data: make(map[string]map[string]bool), Count: 0}
	}

	currentDay := time.Now().Format("2006-01-02")
	if _, ok := referrers[referrer].Data[currentDay]; !ok {
		referrers[referrer].Data[currentDay] = make(map[string]bool)
	}

	if _, ok := referrers[referrer].Data[currentDay][sessionId]; !ok {
		referrers[referrer].Data[currentDay][sessionId] = true
		referrers[referrer].Count += 1
	}

	return true
}

func visitsHandler(event string, data map[string]interface{}) bool {

	if "api.visits" != event {
		return true
	}

	sessionId := utils.InterfaceToString(data["sessionId"])
	today := time.Now().Format("2006-01-02")
	yesterday := time.Now().AddDate(0, 0, -1).Format("2006-01-02")

	if visits.Today != today {
		visits.Yesterday = yesterday
		visits.Today = today
	}

	if _, ok := visits.Data[today]; !ok {
		// 2,1,2
		visits.Data["2014-10-22"] = make(map[string]int32)
		visits.Data["2014-10-22"]["b0"] = int32(1413939794) // 2014-10-22 01:03:14
		visits.Data["2014-10-22"]["b1"] = int32(1413940274) // 2014-10-22 01:11:14
		visits.Data["2014-10-22"]["b2"] = int32(1413990674) // 2014-10-22 15:11:14
		visits.Data["2014-10-22"]["b3"] = int32(1414005074) // 2014-10-22 19:11:14
		visits.Data["2014-10-22"]["b4"] = int32(1414005062) // 2014-10-22 19:11:02

		// 1,2,2,3,1,1
		visits.Data["2014-10-23"] = make(map[string]int32)
		visits.Data["2014-10-23"]["a2"] = int32(1414028132) // 2014-10-23 01:35:32
		visits.Data["2014-10-23"]["a3"] = int32(1414042542) // 2014-10-23 05:35:42
		visits.Data["2014-10-23"]["a5"] = int32(1414040722) // 2014-10-23 05:05:22
		visits.Data["2014-10-23"]["a0"] = int32(1414065912) // 2014-10-23 12:05:12
		visits.Data["2014-10-23"]["a4"] = int32(1414066552) // 2014-10-23 12:15:52
		visits.Data["2014-10-23"]["a6"] = int32(1414069962) // 2014-10-23 13:12:42
		visits.Data["2014-10-23"]["a7"] = int32(1414071672) // 2014-10-23 13:41:12
		visits.Data["2014-10-23"]["a8"] = int32(1414072342) // 2014-10-23 13:52:22
		visits.Data["2014-10-23"]["a1"] = int32(1414077202) // 2014-10-23 15:13:22
		visits.Data["2014-10-23"]["a9"] = int32(1414102532) // 2014-10-23 22:15:32

		// 1,2,2,3,1,1
		visits.Data[yesterday] = make(map[string]int32)
		visits.Data[yesterday]["a2"] = int32(1414373732) // 2014-10-27 01:35:32
		visits.Data[yesterday]["a3"] = int32(1414388142) // 2014-10-27 05:35:42
		visits.Data[yesterday]["a5"] = int32(1414386322) // 2014-10-27 05:05:22
		visits.Data[yesterday]["a0"] = int32(1414411512) // 2014-10-27 12:05:12
		visits.Data[yesterday]["a4"] = int32(1414412152) // 2014-10-27 12:15:52
		visits.Data[yesterday]["a6"] = int32(1414415562) // 2014-10-27 13:12:42
		visits.Data[yesterday]["a7"] = int32(1414417272) // 2014-10-27 13:41:12
		visits.Data[yesterday]["a8"] = int32(1414417942) // 2014-10-27 13:52:22
		visits.Data[yesterday]["a1"] = int32(1414422802) // 2014-10-27 15:13:22
		visits.Data[yesterday]["a9"] = int32(1414448132) // 2014-10-27 22:15:32

		visits.Data[today] = make(map[string]int32)
	}

	if _, ok := visits.Data[today][sessionId]; !ok {
		visits.Data[today][sessionId] = int32(time.Now().Unix())
		if _, ok := conversions["visitors"]; !ok {
			conversions["visitors"] = make(map[string]int)
		}
		conversions["visitors"]["count"] += 1
	}

	return true
}

func addToCartHandler(event string, data map[string]interface{}) bool {

	if "api.addToCart" != event {
		return true
	}

	sessionId := utils.InterfaceToString(data["sessionId"])
	if "" == sessionId {
		return true
	}

	if _, ok := conversions["addedToCart"]; !ok {
		conversions["addedToCart"] = make(map[string]int)
	}

	if _, ok := conversions["addedToCart"][sessionId]; !ok {
		conversions["addedToCart"][sessionId] = 1
	}


	return true
}

func reachedCheckoutHandler(event string, data map[string]interface{}) bool {

	if "api.reachedCheckout" != event {
		return true
	}

	sessionId := utils.InterfaceToString(data["sessionId"])
	if "" == sessionId {
		return true
	}

	if _, ok := conversions["reachedCheckout"]; !ok {
		conversions["reachedCheckout"] = make(map[string]int)
	}

	if _, ok := conversions["reachedCheckout"][sessionId]; !ok {
		conversions["reachedCheckout"][sessionId] = 1
	}

	return true
}

func purchasedHandler(event string, data map[string]interface{}) bool {

	if "api.purchased" != event {
		return true
	}

	sessionId := utils.InterfaceToString(data["sessionId"])
	if "" == sessionId {
		return true
	}

	if _, ok := conversions["purchased"]; !ok {
		conversions["purchased"] = make(map[string]int)
	}

	if _, ok := conversions["purchased"][sessionId]; !ok {
		conversions["purchased"][sessionId] = 1
	}

	return true
}

func salesHandler(event string, data map[string]interface{}) bool {

	if "api.sales" != event || len(data) == 0 {
		return true
	}
	salesData := make(map[string]int)

	salesHistoryCollection, err := db.GetCollection(COLLECTION_NAME_SALES_HISTORY)
	if err != nil {
		return true
	}


	for productId, count := range data {
		year := time.Now().Year()
		month := time.Now().Month()
		day := time.Now().Day()
		date := time.Date(year, month, day, 0, 0, 0, 0, time.Local)
		salesHistoryRow := make(map[string]interface{})
		salesData[productId] = utils.InterfaceToInt(count)

		salesHistoryCollection.ClearFilters()
		salesHistoryCollection.AddFilter("created_at", "=", date)
		salesHistoryCollection.AddFilter("product_id", "=", productId)
		dbSaleRow, _ := salesHistoryCollection.Load()
		fmt.Println(dbSaleRow)
		newCount  := utils.InterfaceToInt(count)
		if len(dbSaleRow) > 0 {
			salesHistoryRow["_id"] = utils.InterfaceToString(dbSaleRow[0]["_id"])
			oldCount := utils.InterfaceToInt(dbSaleRow[0]["count"])
			newCount += oldCount
		}

		// Add history row
		salesHistoryRow["product_id"] = productId
		salesHistoryRow["created_at"] = date
		salesHistoryRow["count"] = newCount
		_, err = salesHistoryCollection.Save(salesHistoryRow)
		if err != nil {
			return true
		}
	}

	SaveSalesData(salesData)

	return true
}
